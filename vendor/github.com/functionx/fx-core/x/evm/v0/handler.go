package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/evm/v0/keeper"

	"github.com/functionx/fx-core/x/evm/v0/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(server types.MsgServer) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgEthereumTx:
			// execute state transition
			res, err := server.EthereumTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			err := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, err
		}
	}
}

func NewEvmProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.InitEvmParamsProposal:
			if ctx.BlockHeight() < fxtype.EvmV0SupportBlock() {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "evm module not enable")
			}
			return k.HandleInitEvmParamsProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized evm proposal content type: %T", c)
		}
	}
}
