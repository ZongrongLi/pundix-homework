package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/gravity/types"
)

const (
	flagEthKeyType  = "eth-key-type"
	flagEthKeystore = "eth-keystore"
	flagEthPassword = "eth-password"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gravity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand([]*cobra.Command{
		// set delegate address
		CmdSetOrchestratorAddress(),
		// send to eth
		CmdSendToEth(),
		CmdCancelSendToEth(),
		CmdRequestBatch(),

		// validator consensus confirm
		CmdValidatorSetConfirm(),
		CmdRequestBatchConfirm(),
	}...)

	return cmd
}

func CmdSetOrchestratorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-orchestrator-address [validator-address] [orchestrator-address] [eth-address]",
		Short: "Allows validators to delegate their voting responsibilities to a given key.",
		Example: "fxcored tx gravity set-orchestrator-address fxvaloper1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8wpk9j2 " +
			"fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2 0xb86d4DC8e2C57190c1cfb834fE69A7a65E2756C2",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			orchestratorAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			ethAddress := args[2]
			if !gethcommon.IsHexAddress(ethAddress) {
				return fmt.Errorf("invalid eth address:%v", ethAddress)
			}
			msg := types.MsgSetOrchestratorAddress{
				Validator:    validatorAddress.String(),
				Orchestrator: orchestratorAddress.String(),
				EthAddress:   ethAddress,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSendToEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-to-eth [eth-dest] [amount] [bridge-fee]",
		Short:   "Adds a new entry to the transaction pool to withdraw an amount from the Ethereum bridge contract",
		Example: "fxcored tx gravity send-to-eth 0xb86d4DC8e2C57190c1cfb834fE69A7a65E2756C2 1000FX 20FX",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := cliCtx.GetFromAddress()

			if !gethcommon.IsHexAddress(args[0]) {
				return sdkerrors.Wrap(fmt.Errorf("invalid bsd-dest address:%v", args[0]), "eth-dest")
			}
			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}
			bridgeFee, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "bridge fee")
			}

			if len(amount) > 1 || len(bridgeFee) > 1 {
				return fmt.Errorf("coin amounts too long, expecting just 1 coin amount for both amount and bridgeFee")
			}

			// Make the message
			msg := types.MsgSendToEth{
				Sender:    sender.String(),
				EthDest:   args[0],
				Amount:    amount[0],
				BridgeFee: bridgeFee[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelSendToEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-eth [txID]",
		Short: "Cancel transaction send to eth",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			senderAddr := cliCtx.GetFromAddress()
			txId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgCancelSendToEth(senderAddr, txId)
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build-batch [token_contract] [minimum_fee] [eth_fee_receive]",
		Short:   "Build a new batch on the fxcore side for pooled withdrawal transactions",
		Example: "fxcored tx gravity build-batch 0xb86d4DC8e2C57190c1cfb834fE69A7a65E2756C2 1 0xb86d4DC8e2C57190c1cfb834fE69A7a65E2756C2",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := cliCtx.GetFromAddress()

			minimumFee, ok := sdk.NewIntFromString(args[1])
			if !ok || minimumFee.IsNegative() {
				return fmt.Errorf("miniumu fee is valid!!!fee:%v\n", args[1])
			}
			ethFeeReceive := args[2]
			if !gethcommon.IsHexAddress(ethFeeReceive) {
				return fmt.Errorf("invalid ethFeeReceive address:%v", args[2])
			}
			baseFee := sdk.ZeroInt()
			baseFeeStr, err := cmd.Flags().GetString("base-fee")
			if err == nil {
				baseFeeStr = strings.TrimSpace(baseFeeStr)
				if len(baseFeeStr) > 0 {
					baseFee, ok = sdk.NewIntFromString(baseFeeStr)
					if !ok {
						return fmt.Errorf("invalid baseFee:%v", baseFeeStr)
					}
				}
			}
			msg := types.NewMsgRequestBatch(fromAddress, args[0], minimumFee, ethFeeReceive, baseFee)
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("base-fee", "", "requestBatch baseFee, is empty is sdk.ZeroInt")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatchConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-confirm [contractAddress] [nonce] [...hexEthPrivate]",
		Short: "Send valset confirm msg",
		Example: fmt.Sprintf("1. fxcored tx gravity valset-confirm 0x30dA8589BFa1E509A319489E014d384b87815D89 1 hexEthPrivateKey --eth-key-type=hex\n" +
			"2. fxcored tx gravity valset-confirm 0x30dA8589BFa1E509A319489E014d384b87815D89 1 --eth-key-type=keystore --eth-keystore=./key --eth-password=./password"),
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()
			contractAddress := args[0]
			if !gethcommon.IsHexAddress(contractAddress) {
				return fmt.Errorf("invalid contract address:%v", contractAddress)
			}

			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			ethKeyType, err := cmd.Flags().GetString(flagEthKeyType)
			if err != nil {
				return err
			}

			var ethPrivateKey *ecdsa.PrivateKey
			switch ethKeyType {
			case "hex":
				if len(args) < 3 {
					return fmt.Errorf("eth-key-type=hex must input hexEthPrivateKey")
				}
				ethPrivateKey, err = recoveryPrivateKeyHexPrivateKey(args[2])
				if err != nil {
					return err
				}
			case "keystore":
				keystoreFile, err := cmd.Flags().GetString(flagEthKeystore)
				if err != nil {
					return err
				}
				passwordFile, err := cmd.Flags().GetString(flagEthPassword)
				if err != nil {
					return err
				}
				ethPrivateKey, err = recoveryPrivateKeyByKeystore(keystoreFile, passwordFile)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown eth-key-type flag:%v, support:(keystore|hex)", ethKeyType)
			}
			ethAddress := ethCrypto.PubkeyToAddress(ethPrivateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			batchRequestByNonceResp, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				Nonce:         nonce,
				TokenContract: contractAddress,
			})
			if err != nil {
				return err
			}
			if batchRequestByNonceResp.Batch == nil {
				return fmt.Errorf("not found batch request by nonce!!!contractAddress:[%v], nonce:[%v]", contractAddress, nonce)
			}
			// Determine whether it has been confirmed
			batchConfirmResp, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				Nonce:         nonce,
				TokenContract: contractAddress,
				Address:       fromAddress.String(),
			})
			if err != nil {
				return err
			}
			if batchConfirmResp.GetConfirm() != nil {
				confirm := batchConfirmResp.GetConfirm()
				return clientCtx.PrintString(fmt.Sprintf("already confirm requestBatch!!!\n\tnonce:[%v]\n\ttokenContract:[%v]\n\torchestrator:[%v]\n\tethAddress:[%v]\n\tsignature:[%v]\n",
					confirm.Nonce, confirm.TokenContract, confirm.Orchestrator, confirm.EthSigner, confirm.Signature))
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			checkpoint, err := batchRequestByNonceResp.GetBatch().GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := types.NewEthereumSignature(checkpoint, ethPrivateKey)
			if err != nil {
				return err
			}
			msg := types.NewMsgConfirmBatch(nonce, contractAddress, ethAddress.String(), hex.EncodeToString(signature), fromAddress)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		}}

	cmd.Flags().String(flagEthKeyType, "keystore", "eth private key type(keystore|hex), default:keystore")
	cmd.Flags().String(flagEthKeystore, "", "eth keystore file")
	cmd.Flags().String(flagEthPassword, "", "eth keystore password file")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdValidatorSetConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valset-confirm [nonce] [...hexEthPrivate]",
		Short: "Send valset confirm msg",
		Example: fmt.Sprintf("1. fxcored tx gravity valset-confirm 1 hexEthPrivateKey --eth-key-type=hex\n" +
			"2. fxcored tx gravity valset-confirm 1 --eth-key-type=keystore --eth-keystore=./key --eth-password=./password"),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			ethKeyType, err := cmd.Flags().GetString(flagEthKeyType)
			if err != nil {
				return err
			}

			var ethPrivateKey *ecdsa.PrivateKey
			switch ethKeyType {
			case "hex":
				if len(args) < 2 {
					return fmt.Errorf("eth-key-type=hex must input hexEthPrivateKey")
				}
				ethPrivateKey, err = recoveryPrivateKeyHexPrivateKey(args[1])
				if err != nil {
					return errors.WithStack(err)
				}
			case "keystore":
				keystoreFile, err := cmd.Flags().GetString(flagEthKeystore)
				if err != nil {
					return err
				}
				passwordFile, err := cmd.Flags().GetString(flagEthPassword)
				if err != nil {
					return err
				}
				ethPrivateKey, err = recoveryPrivateKeyByKeystore(keystoreFile, passwordFile)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown eth-key-type flag:%v, support:(keystore|hex)", ethKeyType)
			}
			ethAddress := ethCrypto.PubkeyToAddress(ethPrivateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			valsetRequestResp, err := queryClient.ValsetRequest(cmd.Context(), &types.QueryValsetRequestRequest{Nonce: nonce})
			if err != nil {
				return err
			}
			// Determine whether it has been confirmed
			valsetConfirmResp, err := queryClient.ValsetConfirm(cmd.Context(), &types.QueryValsetConfirmRequest{
				Nonce:   nonce,
				Address: fromAddress.String(),
			})
			if err != nil {
				return err
			}
			if valsetConfirmResp.GetConfirm() != nil {
				confirm := valsetConfirmResp.GetConfirm()
				return fmt.Errorf("already confirm valset!!!\n\tnonce:[%v]\n\torchestrator:[%v]\n\tethAddress:[%v]\n\tsignature:[%v]\n", confirm.Nonce, confirm.Orchestrator, confirm.EthAddress, confirm.Signature)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			checkpoint := valsetRequestResp.GetValset().GetCheckpoint(paramsResp.Params.GetGravityId())
			signature, err := types.NewEthereumSignature(checkpoint, ethPrivateKey)
			if err != nil {
				return err
			}
			msg := types.NewMsgValsetConfirm(nonce, ethAddress.String(), fromAddress, hex.EncodeToString(signature))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		}}

	cmd.Flags().String(flagEthKeyType, "keystore", "eth private key type(keystore|hex), default:keystore")
	cmd.Flags().String(flagEthKeystore, "", "eth keystore file")
	cmd.Flags().String(flagEthPassword, "", "eth keystore password file")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func recoveryPrivateKeyByKeystore(keystoreFile, passwordFile string) (*ecdsa.PrivateKey, error) {
	keystoreData, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "keystoreFile:[%s]", keystoreFile)
	}
	passwordData, err := ioutil.ReadFile(passwordFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "passwordFile:[%s]", keystoreFile)
	}
	decryptKey, err := keystore.DecryptKey(keystoreData, string(passwordData))
	if err != nil {
		return nil, errors.WithMessagef(err, "decryptKey err")
	}
	return decryptKey.PrivateKey, nil
}

func recoveryPrivateKeyHexPrivateKey(hexPrivateKey string) (*ecdsa.PrivateKey, error) {
	return ethCrypto.HexToECDSA(hexPrivateKey)
}
