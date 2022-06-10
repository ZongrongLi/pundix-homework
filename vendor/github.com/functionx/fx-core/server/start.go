package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"

	abciserver "github.com/tendermint/tendermint/abci/server"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/rpc/client/local"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	"github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethdebug "github.com/functionx/fx-core/rpc/namespaces/ethereum/debug"
	serverconfig "github.com/functionx/fx-core/server/config"
)

const (
	flagWithTendermint     = "with-tendermint"
	flagAddress            = "address"
	flagTransport          = "transport"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"
	FlagInvCheckPeriod     = "inv-check-period"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
	FlagIndexEvents       = "index-events"
	FlagMinRetainBlocks   = "min-retain-blocks"
)

// GRPC-related flags.
const (
	flagGRPCEnable  = "grpc.enable"
	flagGRPCAddress = "grpc.address"

	flagRPCEnable         = "api.enable"
	flagEnabledUnsafeCors = "api.enabled-unsafe-cors"
)

// State sync-related flags.
const (
	FlagStateSyncSnapshotInterval   = "state-sync.snapshot-interval"
	FlagStateSyncSnapshotKeepRecent = "state-sync.snapshot-keep-recent"
)

// JSON-RPC flags
const (
	JSONRPCEnable          = "json-rpc.enable"
	JSONRPCAPI             = "json-rpc.api"
	JSONRPCAddress         = "json-rpc.address"
	JSONWsAddress          = "json-rpc.ws-address"
	JSONRPCGasCap          = "json-rpc.gas-cap"
	JSONRPCEVMTimeout      = "json-rpc.evm-timeout"
	JSONRPCTxFeeCap        = "json-rpc.txfee-cap"
	JSONRPCFilterCap       = "json-rpc.filter-cap"
	JSONRPCLogsCap         = "json-rpc.logs-cap"
	JSONRPCBlockRangeCap   = "json-rpc.block-range-cap"
	JSONRPCHTTPTimeout     = "json-rpc.http-timeout"
	JSONRPCHTTPIdleTimeout = "json-rpc.http-idle-timeout"
)

// EVM flags
const (
	EVMTracer         = "evm.tracer"
	EVMMaxTxGasWanted = "evm.max-tx-gas-wanted"
)

// TLS flags
const (
	TLSCertPath = "tls.certificate-path"
	TLSKeyPath  = "tls.key-path"
)

const (
	ServerStartTime = 5 * time.Second
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// Tendermint.
func StartCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			// Bind flags to the Context's Viper so the app construction can set
			// options accordingly.
			err := serverCtx.Viper.BindPFlags(cmd.Flags())
			if err != nil {
				return err
			}

			_, err = server.GetPruningOptionsFromFlags(serverCtx.Viper)
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			withTM, _ := cmd.Flags().GetBool(flagWithTendermint)
			if !withTM {
				serverCtx.Logger.Info("starting ABCI without Tendermint")
				return startStandAlone(serverCtx, appCreator)
			}

			serverCtx.Logger.Info("Unlocking keyring")

			// fire unlock precess for keyring
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if keyringBackend == keyring.BackendFile {
				_, err = clientCtx.Keyring.List()
				if err != nil {
					return err
				}
			}

			serverCtx.Logger.Info("starting ABCI with Tendermint")

			// amino is needed here for backwards compatibility of REST routes
			err = startInProcess(serverCtx, clientCtx, appCreator)
			errCode, ok := err.(server.ErrorCode)
			if !ok {
				return err
			}

			serverCtx.Logger.Debug(fmt.Sprintf("received quit signal: %d", errCode.Code))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTransport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)")
	cmd.Flags().IntSlice(FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().Uint64(FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(FlagPruning, storetypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningKeepEvery, 0, "Offset heights to keep on disk after 'keep-every' (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint(FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().Uint64(FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune Tendermint blocks")

	cmd.Flags().Bool(flagGRPCEnable, true, "Define if the gRPC server should be enabled")
	cmd.Flags().String(flagGRPCAddress, serverconfig.DefaultGRPCAddress, "the gRPC server address to listen on")

	cmd.Flags().Bool(flagRPCEnable, false, "Defines if Cosmos-sdk REST server should be enabled")
	cmd.Flags().Bool(flagEnabledUnsafeCors, false, "Defines if CORS should be enabled (unsafe - use it at your own risk)")

	cmd.Flags().Bool(JSONRPCEnable, true, "Define if the gRPC server should be enabled")
	cmd.Flags().StringSlice(JSONRPCAPI, serverconfig.GetDefaultAPINamespaces(), "Defines a list of JSON-RPC namespaces that should be enabled")
	cmd.Flags().String(JSONRPCAddress, serverconfig.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().String(JSONWsAddress, serverconfig.DefaultJSONRPCWsAddress, "the JSON-RPC WS server address to listen on")
	cmd.Flags().Uint64(JSONRPCGasCap, serverconfig.DefaultGasCap, "Sets a cap on gas that can be used in eth_call/estimateGas unit is FX (0=infinite)")
	cmd.Flags().Float64(JSONRPCTxFeeCap, serverconfig.DefaultTxFeeCap, "Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 photon)")
	cmd.Flags().Int32(JSONRPCFilterCap, serverconfig.DefaultFilterCap, "Sets the global cap for total number of filters that can be created")
	cmd.Flags().Duration(JSONRPCEVMTimeout, serverconfig.DefaultEVMTimeout, "Sets a timeout used for eth_call (0=infinite)")
	cmd.Flags().Duration(JSONRPCHTTPTimeout, serverconfig.DefaultHTTPTimeout, "Sets a read/write timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Duration(JSONRPCHTTPIdleTimeout, serverconfig.DefaultHTTPIdleTimeout, "Sets a idle timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Int32(JSONRPCLogsCap, serverconfig.DefaultLogsCap, "Sets the max number of results can be returned from single `eth_getLogs` query")
	cmd.Flags().Int32(JSONRPCBlockRangeCap, serverconfig.DefaultBlockRangeCap, "Sets the max block range allowed for `eth_getLogs` query")

	cmd.Flags().String(EVMTracer, serverconfig.DefaultEVMTracer, "the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)")
	cmd.Flags().Uint64(EVMMaxTxGasWanted, serverconfig.DefaultMaxTxGasWanted, "the gas wanted for each eth tx returned in ante handler in check tx mode")

	cmd.Flags().String(TLSCertPath, "", "the cert.pem file path for the server TLS configuration")
	cmd.Flags().String(TLSKeyPath, "", "the key.pem file path for the server TLS configuration")

	cmd.Flags().Uint64(FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(ctx *server.Context, appCreator types.AppCreator) error {
	addr := ctx.Viper.GetString(flagAddress)
	transport := ctx.Viper.GetString(flagTransport)
	home := ctx.Viper.GetString(flags.FlagHome)

	db, err := openDB(home)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			ctx.Logger.With("error", err).Error("error closing db")
		}
	}()

	traceWriterFile := ctx.Viper.GetString(flagTraceStore)
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	svr, err := abciserver.NewServer(addr, transport, app)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	svr.SetLogger(ctx.Logger.With("server", "abci"))

	err = svr.Start()
	if err != nil {
		tmos.Exit(err.Error())
	}

	defer func() {
		if err = svr.Stop(); err != nil {
			tmos.Exit(err.Error())
		}
	}()

	// Wait for SIGINT or SIGTERM signal
	return server.WaitForQuitSignals()
}

// legacyAminoCdc is used for the legacy REST API
func startInProcess(ctx *server.Context, clientCtx client.Context, appCreator types.AppCreator) (err error) {
	cfg := ctx.Config
	home := cfg.RootDir
	logger := ctx.Logger
	var cpuProfileCleanup func() error

	if cpuProfile := ctx.Viper.GetString(flagCPUProfile); cpuProfile != "" {
		fp, err := ethdebug.ExpandHome(cpuProfile)
		if err != nil {
			ctx.Logger.Debug("failed to get filepath for the CPU profile file", "error", err.Error())
			return err
		}
		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		cpuProfileCleanup = func() error {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			if err := f.Close(); err != nil {
				logger.Error("failed to close CPU profiler file", "error", err.Error())
				return err
			}
			return nil
		}
	}

	traceWriterFile := ctx.Viper.GetString(flagTraceStore)
	db, err := openDB(home)
	if err != nil {
		logger.Error("failed to open DB", "error", err.Error())
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			ctx.Logger.With("error", err).Error("error closing db")
		}
	}()

	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		logger.Error("failed to open trace writer", "error", err.Error())
		return err
	}

	config := serverconfig.GetConfig(ctx.Viper)

	if err := config.ValidateBasic(); err != nil {
		if strings.Contains(err.Error(), "set min gas price in app.toml or flag or env variable") {
			ctx.Logger.Error(
				"WARNING: The minimum-gas-prices config in app.toml is set to the empty string. " +
					"This defaults to 0 in the current version, but will error in the next version " +
					"(SDK v0.44). Please explicitly put the desired minimum-gas-prices in your app.toml.",
			)
		} else {
			return err
		}
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		logger.Error("failed load or gen node key", "error", err.Error())
		return err
	}

	genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)
	tmNode, err := node.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genDocProvider,
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("server", "node"),
	)
	if err != nil {
		logger.Error("failed init node", "error", err.Error())
		return err
	}
	logger.Debug("initialization: tmNode created")
	if err := tmNode.Start(); err != nil {
		logger.Error("failed start tendermint server", "error", err.Error())
		return err
	}
	logger.Debug("initialization: tmNode started")

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC  or JSONRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local tendermint RPC client.
	if config.API.Enable || config.GRPC.Enable || config.JSONRPC.Enable {
		clientCtx = clientCtx.WithClient(local.New(tmNode))

		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)
	}

	var apiSrv *api.Server
	if config.API.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}

		clientCtx := clientCtx.
			WithHomeDir(home).
			WithChainID(genDoc.ChainID)

		apiSrv = api.New(clientCtx, ctx.Logger.With("server", "api"))
		app.RegisterAPIRoutes(apiSrv, config.API)
		errCh := make(chan error)

		go func() {
			if err := apiSrv.Start(config.Config); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(ServerStartTime): // assume server started successfully
		}
	}

	var (
		grpcSrv *grpc.Server
	)
	if config.GRPC.Enable {
		grpcSrv, err = servergrpc.StartGRPCServer(clientCtx, app, config.GRPC.Address)
		if err != nil {
			return err
		}
	}

	var (
		httpSrv     *http.Server
		httpSrvDone chan struct{}
	)
	if config.JSONRPC.Enable {
		go func() {

			web3Logger := ctx.Logger.With("Web3JsonRpc")
			for {
				block, err := clientCtx.Client.Block(context.Background(), nil)
				if err != nil {
					web3Logger.Error(fmt.Sprintf("Query block failed!err:%s", err.Error()))
					time.Sleep(30 * time.Second)
					continue
				}
				if block == nil || block.Block == nil || block.Block.Height < fxtypes.EvmV1SupportBlock() {
					web3Logger.Debug("Evm Module not enable sleep 30s")
					time.Sleep(30 * time.Second)
					continue
				}
				genDoc, err := genDocProvider()
				if err != nil {
					web3Logger.Error("load genesis err!!!err", err)
					panic(fmt.Sprintf("load genesis err!!err:%v", err))
				}
				startJsonRpcClientCtx := clientCtx.WithChainID(genDoc.ChainID)
				tmRPCAddr := cfg.RPC.ListenAddress
				web3Logger.Info("Evm module enable!start json rpc server", "chainId", genDoc.ChainID, "rpcAddr", tmRPCAddr)
				httpSrv, httpSrvDone, err = StartJSONRPC(ctx, startJsonRpcClientCtx, tmRPCAddr, "/websocket", config)
				if err != nil {
					panic(fmt.Sprintf("start json rpc server !!tmRpcAddr:%v, err:%v", tmRPCAddr, err))
				}
				break
			}
		}()
	}

	defer func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}

		if cpuProfileCleanup != nil {
			_ = cpuProfileCleanup()
		}

		if apiSrv != nil {
			_ = apiSrv.Close()
		}

		if grpcSrv != nil {
			grpcSrv.Stop()
		}

		if httpSrv != nil {
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()

			if err := httpSrv.Shutdown(shutdownCtx); err != nil {
				logger.Error("HTTP server shutdown produced a warning", "error", err.Error())
			} else {
				logger.Info("HTTP server shut down, waiting 5 sec")
				select {
				case <-time.Tick(5 * time.Second):
				case <-httpSrvDone:
				}
			}
		}

		logger.Info("Bye!")
	}()

	// Wait for SIGINT or SIGTERM signal
	return server.WaitForQuitSignals()
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB("application", dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}

	filePath := filepath.Clean(traceWriterFile)
	return os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o600,
	)
}
