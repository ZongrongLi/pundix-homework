package clients

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/crypto/hd"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	rpcURI             = "https://fx-json.functionx.io:26657"
	singaporeValidator = "fxvaloper1a73plz6w7fc8ydlwxddanc7a239kk45jnl9xwj"
	userAccount1       = "fx15sy7ph7j6vma607y80cxdc7qg7pgvjdhnql3q6" // pick from explorer randomly
)

func newClientContext() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		// WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode("sync").
		WithHomeDir(app.DefaultNodeHome).
		WithViper("FX").
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithChainID("fxcore")

	clientCtx = clientCtx.WithNodeURI(rpcURI)

	client, err := rpchttp.New(rpcURI, "/websocket")
	if err != nil {
		panic(err)
	}

	return clientCtx.WithClient(client)
}
