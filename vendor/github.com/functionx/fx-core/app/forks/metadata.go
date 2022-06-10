package forks

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	fxtypes "github.com/functionx/fx-core/types"
)

func UpdateMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper) {
	bankKeeper.DeleteDenomMetaData(ctx, strings.ToLower(fxtypes.DefaultDenom))
	bankKeeper.SetDenomMetaData(ctx, fxtypes.GetFxBankMetaData(fxtypes.DefaultDenom))
}
