package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/evm/keeper"
	"github.com/functionx/fx-core/x/evm/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k *keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "evm module not enable")
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgEthereumTx:
			// execute state transition
			res, err := k.EthereumTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			err := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, err
		}
	}
}
