package keeper

import (
	"context"

	fxtypes "github.com/functionx/fx-core/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/feemarket/types"
)

var _ types.QueryServer = Keeper{}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
		return nil, status.Error(codes.InvalidArgument, types.ErrNotInitializedOrUnknownBlock.Error())
	}

	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}

// BaseFee implements the Query/BaseFee gRPC method
func (k Keeper) BaseFee(c context.Context, _ *types.QueryBaseFeeRequest) (*types.QueryBaseFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
		return nil, status.Error(codes.InvalidArgument, types.ErrNotInitializedOrUnknownBlock.Error())
	}

	res := &types.QueryBaseFeeResponse{}
	baseFee := k.GetBaseFee(ctx)

	if baseFee != nil {
		aux := sdk.NewIntFromBigInt(baseFee)
		res.BaseFee = &aux
	}

	return res, nil
}

// BlockGas implements the Query/BlockGas gRPC method
func (k Keeper) BlockGas(c context.Context, _ *types.QueryBlockGasRequest) (*types.QueryBlockGasResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
		return nil, status.Error(codes.InvalidArgument, types.ErrNotInitializedOrUnknownBlock.Error())
	}

	gas := k.GetBlockGasUsed(ctx)

	return &types.QueryBlockGasResponse{
		Gas: int64(gas),
	}, nil
}
