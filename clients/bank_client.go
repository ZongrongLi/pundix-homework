package clients

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

type BankQueryClient struct {
	Context client.Context
	Client  types.QueryClient
}

func (b *BankQueryClient) New() {
	b.Context = newClientContext()
	b.Client = types.NewQueryClient(b.Context)
}

func (b *BankQueryClient) Balance(address string) (*types.QueryBalanceResponse, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	params := types.NewQueryBalanceRequest(addr, "FX")
	res, err := b.Client.Balance(context.Background(), params)
	if err != nil {
		return nil, err
	}
	b.Context.PrintProto(res.Balance)
	return res, nil
}

func (b *BankQueryClient) TotalSupply() (*types.QuerySupplyOfResponse, error) {
	// pageReq = &query.PageRequest{
	// 	Limit: 100,
	// }

	res, err := b.Client.SupplyOf(context.Background(), &types.QuerySupplyOfRequest{Denom: "FX"})
	if err != nil {
		return nil, err
	}

	b.Context.PrintProto(&res.Amount)
	return res, nil
}
