package keeper

import (
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/x/erc20/types"
)

func (k Keeper) RelayTransferCrossChainProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) (err error) {
	fip20ABI := fxtypes.GetERC20(ctx.BlockHeight()).ABI
	for _, log := range receipt.Logs {
		tc, isOk, err := fxtypes.ParseTransferCrossChainEvent(fip20ABI, log)
		if err != nil {
			return err
		}
		if !isOk {
			continue
		}
		pair, found := k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		k.Logger(ctx).Info("transfer cross", "tx-hash", receipt.TxHash.Hex(), "from", from.Hex(), "to", to.Hex(), "token", pair.Erc20Address, "denom", pair.Denom)

		balances := k.bankKeeper.GetAllBalances(ctx, tc.From.Bytes())
		if !balances.IsAllGTE(tc.TotalAmount(pair.Denom)) {
			return fmt.Errorf("insufficient balance, have %s expected %s", balances.String(), tc.TotalAmount(pair.Denom).String())
		}

		targetType, target := tc.GetTarget()
		switch targetType {
		case fxtypes.FIP20TargetChain:
			err = k.TransferChainHandler(ctx, tc.GetFrom(), tc.Recipient, tc.GetAmount(pair.Denom), tc.GetFee(pair.Denom), target, receipt)
		case fxtypes.FIP20TargetIBC:
			err = k.TransferIBCHandler(ctx, tc.GetFrom(), tc.Recipient, tc.GetAmount(pair.Denom), tc.GetFee(pair.Denom), target, receipt)
		default:
			err = fmt.Errorf("traget unknown %d", targetType)
		}
		if err != nil {
			k.Logger(ctx).Error("failed to transfer cross chain", "tx-hash", receipt.TxHash.Hex(), "error", err.Error())
			return err
		}
		k.Logger(ctx).Info("transfer cross chain success", "tx-hash", receipt.TxHash.Hex())

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_transfer_cross_chain"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", pair.Erc20Address),
				telemetry.NewLabel("denom", pair.Denom),
				telemetry.NewLabel("type", targetType.String()),
				telemetry.NewLabel("target", target),
				telemetry.NewLabel("amount", tc.GetAmount(pair.Denom).String()),
			},
		)
	}
	return nil
}

func (k Keeper) TransferChainHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string, _ *ethtypes.Receipt) error {
	k.Logger(ctx).Info("transfer chain handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)
	router := k.GetRouter()
	if router == nil || !router.HasRoute(target) {
		return fmt.Errorf("target %s not support", target)
	}
	route, _ := router.GetRoute(target)
	return route.TransferAfter(ctx, from.String(), to, amount, fee)
}

func (k Keeper) TransferIBCHandler(ctx sdk.Context, from sdk.AccAddress, to string, amount, fee sdk.Coin, target string, receipt *ethtypes.Receipt) error {
	k.Logger(ctx).Info("transfer ibc handler", "from", from, "to", to, "amount", amount.String(), "fee", fee.String(), "target", target)
	targetIBC, ok := fxtypes.ParseTargetIBC(target)
	if !ok {
		return fmt.Errorf("invalid target ibc %s", target)
	}
	if _, err := sdk.GetFromBech32(to, targetIBC.Prefix); err != nil {
		return fmt.Errorf("invalid to address %s", to)
	}
	_, _, err := k.ibcChannelKeeper.GetChannelClientState(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if err != nil {
		return err
	}
	params := k.GetParams(ctx)
	ibcTimeoutHeight := ibcclienttypes.ZeroHeight()
	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(params.IbcTimeout)
	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if !found {
		return fmt.Errorf("ibc channel next sequence send not found, port %s, channel %s", targetIBC.SourcePort, targetIBC.SourceChannel)
	}
	ctx.Logger().Info("ibc transfer", "port", targetIBC.SourcePort, "channel", targetIBC.SourceChannel, "sequence", nextSequenceSend, "timeout-height", ibcTimeoutHeight)
	if err := k.ibcTransferKeeper.SendTransfer(
		ctx, targetIBC.SourcePort, targetIBC.SourceChannel, amount, from.Bytes(),
		to, ibcTimeoutHeight, ibcTimeoutTimestamp, "", fee); err != nil {
		return err
	}
	k.setIBCTransferHash(ctx, targetIBC.SourcePort, targetIBC.SourceChannel, nextSequenceSend, receipt.TxHash)
	return nil
}

func (k Keeper) setIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64, hash common.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(port, channel, sequence), hash.Bytes())
}

func (k Keeper) GetIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) (common.Hash, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIBCTransferKey(port, channel, sequence)
	if !store.Has(key) {
		return common.Hash{}, false
	}
	value := store.Get(key)
	return common.BytesToHash(value), true
}

func (k Keeper) HashIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetIBCTransferKey(port, channel, sequence))
}
