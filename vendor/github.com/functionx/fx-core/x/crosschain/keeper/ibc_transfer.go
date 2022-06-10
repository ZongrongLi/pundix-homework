package keeper

import (
	"encoding/hex"
	"fmt"

	gravitytypes "github.com/functionx/fx-core/x/gravity/types"

	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var targetEvmPrefix = hex.EncodeToString([]byte("module/evm"))

func (a Keeper) handlerRelayTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	if ctx.BlockHeight() >= fxtypes.EvmV1SupportBlock() && claim.TargetIbc == targetEvmPrefix {
		a.handlerEvmTransfer(ctx, claim, receiver, coin)
		return
	}
	a.handleIbcTransfer(ctx, claim, receiver, coin)
}

func (k Keeper) handleIbcTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiveAddr sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	targetIBC, ok := fxtypes.ParseHexTargetIBC(claim.TargetIbc)
	if !ok {
		logger.Error("convert target ibc data error!!!", "targetIbc", claim.GetTargetIbc())
		return
	}
	ibcReceiveAddress, err := gravitytypes.CovertIbcPacketReceiveAddressByPrefix(ctx.BlockHeight(), targetIBC.Prefix, receiveAddr)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "fxReceive", claim.Receiver,
			"ibcPrefix", targetIBC.Prefix, "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel, "error", err)
		return
	}

	_, clientState, err := k.ibcChannelKeeper.GetChannelClientState(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel)
		return
	}

	params := k.GetParams(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	destTimeoutHeight := clientStateHeight.GetRevisionHeight() + params.IbcTransferTimeoutHeight
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: destTimeoutHeight,
	}

	nextSequenceSend, found := k.ibcChannelKeeper.GetNextSequenceSend(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if !found {
		logger.Error("ibc channel next sequence send not found!!!", "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel)
		return
	}
	logger.Info("CrossChain start ibc transfer", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "destCurrentHeight", clientStateHeight.GetRevisionHeight(), "destTimeoutHeight", destTimeoutHeight, "nextSequenceSend", nextSequenceSend)

	if err = k.ibcTransferKeeper.SendTransfer(ctx,
		targetIBC.SourcePort, targetIBC.SourceChannel,
		coin, receiveAddr, ibcReceiveAddress,
		ibcTimeoutHeight, 0,
		"", sdk.NewCoin(coin.Denom, sdk.ZeroInt())); err != nil {
		logger.Error("CrossChain ibc transfer fail", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "err", err)
		return
	}
	k.SetIbcSequenceHeight(ctx, targetIBC.SourcePort, targetIBC.SourceChannel, nextSequenceSend, uint64(ctx.BlockHeight()))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(nextSequenceSend)),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, targetIBC.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, targetIBC.SourceChannel),
	))
}

func (k Keeper) handlerEvmTransfer(ctx sdk.Context, claim *types.MsgSendToFxClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	logger := k.Logger(ctx)
	receiverEthType := common.BytesToAddress(receiver.Bytes())
	logger.Info("convert denom to fip20", "sender", claim.Sender, "receiver", claim.Receiver,
		"receiver-eth-type", receiverEthType.String(), "amount", coin.String(), "target", claim.TargetIbc)
	err := k.erc20Keeper.RelayConvertCoin(ctx, receiver, receiverEthType, coin)
	if err != nil {
		logger.Error("evm transfer, convert denom to fip20 failed", "error", err.Error())
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	))
}
