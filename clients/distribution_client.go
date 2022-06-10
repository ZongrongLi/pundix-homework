package clients

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type DistributionQueryClient struct {
	Context client.Context
	Client  types.QueryClient
}

func (d *DistributionQueryClient) New() {
	d.Context = newClientContext()
	d.Client = types.NewQueryClient(d.Context)
}

func (d *DistributionQueryClient) QueryParams() (*types.QueryParamsResponse, error) {
	res, err := d.Client.Params(context.Background(), &types.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	if err = d.Context.PrintProto(&res.Params); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DistributionQueryClient) ValidatorOutstandingRewards(validitorAddr string) (*types.QueryValidatorOutstandingRewardsResponse, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(validitorAddr)
	if err != nil {
		return nil, err
	}

	res, err := d.Client.ValidatorOutstandingRewards(
		context.Background(),
		&types.QueryValidatorOutstandingRewardsRequest{ValidatorAddress: validatorAddr.String()},
	)
	if err != nil {
		return nil, err
	}

	if err = d.Context.PrintProto(&res.Rewards); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DistributionQueryClient) ValidatorCommission(validitorAddr string) (*types.QueryValidatorCommissionResponse, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(validitorAddr)
	if err != nil {
		return nil, err
	}

	res, err := d.Client.ValidatorCommission(
		context.Background(),
		&types.QueryValidatorCommissionRequest{ValidatorAddress: validatorAddr.String()},
	)

	if err = d.Context.PrintProto(&res.Commission); err != nil {
		return nil, err
	}
	return res, nil
}
func (d *DistributionQueryClient) ValidatorSlashes(validator string, startHeight, endHeight, limit uint64) (*types.QueryValidatorSlashesResponse, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return nil, err
	}

	pageReq := &query.PageRequest{
		Limit: uint64(limit),
	}

	res, err := d.Client.ValidatorSlashes(
		context.Background(),
		&types.QueryValidatorSlashesRequest{
			ValidatorAddress: validatorAddr.String(),
			StartingHeight:   startHeight,
			EndingHeight:     endHeight,
			Pagination:       pageReq,
		},
	)
	if err != nil {
		return nil, err
	}

	if err = d.Context.PrintProto(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DistributionQueryClient) CommunityPool() (*types.QueryCommunityPoolResponse, error) {
	res, err := d.Client.CommunityPool(context.Background(), &types.QueryCommunityPoolRequest{})
	if err != nil {
		return nil, err
	}

	if err = d.Context.PrintProto(res); err != nil {
		return nil, err
	}
	return res, nil
}
