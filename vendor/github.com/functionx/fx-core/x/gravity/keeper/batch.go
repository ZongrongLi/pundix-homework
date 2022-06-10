package keeper

import (
	"fmt"
	"strings"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/gravity/types"
)

const OutgoingTxBatchSize = 100

// BuildOutgoingTXBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit withtout creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, contractAddress string, maxElements int, minimumFee sdk.Int, feeReceive string, baseFee sdk.Int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}
	lastBatch := k.GetLastOutgoingBatchByTokenType(ctx, contractAddress)

	// lastBatch may be nil if there are no existing batches, we only need
	// to perform this check if a previous batch exists
	if lastBatch != nil {
		// this traverses the current tx pool for this token type and determines what
		// fees a hypothetical batch would have if created
		currentFees := k.GetBatchFeesByTokenType(ctx, contractAddress, baseFee)
		if currentFees == nil {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "error getting fees from tx pool")
		}

		if lastBatch.GetFees().GT(currentFees.TotalFees) {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "new batch would not be more profitable")
		}
	}
	selectedTx, err := k.pickUnbatchedTX(ctx, contractAddress, maxElements, baseFee)
	if err != nil {
		return nil, err
	}
	if len(selectedTx) == 0 {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "no batch tx")
	}
	if types.OutgoingTransferTxs(selectedTx).TotalFee().LT(minimumFee) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "total fee less than minimum fee")
	}
	batchTimeout := k.GetBatchTimeoutHeight(ctx)
	if batchTimeout <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "batch timeout height")
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := &types.OutgoingTxBatch{
		BatchNonce:    nextID,
		BatchTimeout:  batchTimeout,
		Transactions:  selectedTx,
		TokenContract: contractAddress,
		FeeReceive:    feeReceive,
	}
	if err = k.StoreBatch(ctx, batch); err != nil {
		return nil, err
	}

	eventBatchNonceTxIds := strings.Builder{}
	eventBatchNonceTxIds.WriteString(fmt.Sprintf("%d", selectedTx[0].Id))
	for _, tx := range selectedTx[1:] {
		_, _ = eventBatchNonceTxIds.WriteString(fmt.Sprintf(",%d", tx.Id))
	}
	if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
		k.GetBridgeChainID(ctx) // gas used
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(nextID)),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxIds, eventBatchNonceTxIds.String()),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchTimeout, fmt.Sprint(batch.BatchTimeout)),
	))
	return batch, nil
}

// GetBatchTimeoutHeight This gets the batch timeout height in Ethereum blocks.
func (k Keeper) GetBatchTimeoutHeight(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	currentFxHeight := ctx.BlockHeight()
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedEthBlockHeight(ctx)
	if heights.FxBlockHeight == 0 || heights.EthBlockHeight == 0 {
		return 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentFxHeight) - heights.FxBlockHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projectedCurrentEthereumHeight := (projectedMillis / params.AverageEthBlockTime) + heights.EthBlockHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	blocksToAdd := params.TargetBatchTimeout / params.AverageEthBlockTime
	return projectedCurrentEthereumHeight + blocksToAdd
}

// OutgoingTxBatchExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, tokenContract string, nonce uint64) error {
	b := k.GetOutgoingTXBatch(ctx, tokenContract, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}

	// cleanup outgoing TX pool, while these transactions where hidden from GetPoolTransactions
	// they still exist in the pool and need to be cleaned up.
	for _, tx := range b.Transactions {
		k.removePoolEntry(ctx, tx.Id)
	}

	// Iterate through remaining batches
	k.IterateOutgoingTXBatches(ctx, func(key []byte, iterBatch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if iterBatch.BatchNonce < b.BatchNonce {
			if err := k.CancelOutgoingTXBatch(ctx, tokenContract, iterBatch.BatchNonce); err != nil {
				panic(fmt.Sprintf("Failed cancel out batch %s %d while trying to execute %s %d with %s", tokenContract, iterBatch.BatchNonce, tokenContract, nonce, err))
			}
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *b)
	return nil
}

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	store := ctx.KVStore(k.storeKey)
	// set the current block height when storing the batch
	batch.Block = uint64(ctx.BlockHeight())
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshal(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)

	if store.Has(blockKey) {
		return sdkerrors.Wrap(types.ErrDuplicate, fmt.Sprintf("block:[%v] has batch request", batch.Block))
	}
	store.Set(blockKey, k.cdc.MustMarshal(batch))
	return nil
}

// StoreBatchUnsafe stores a transaction batch w/o setting the height
func (k Keeper) StoreBatchUnsafe(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshal(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)
	store.Set(blockKey, k.cdc.MustMarshal(batch))
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTX(ctx sdk.Context, contractAddress string, maxElements int, baseFee sdk.Int) ([]*types.OutgoingTransferTx, error) {
	isSupportBaseFee := fxtypes.IsRequestBatchBaseFee(ctx.BlockHeight())
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateOutgoingPoolByFee(ctx, contractAddress, func(txID uint64, tx *types.OutgoingTransferTx) bool {
		if tx != nil && tx.Erc20Fee != nil {
			if isSupportBaseFee && tx.Erc20Fee.Amount.LT(baseFee) {
				return true
			}
			selectedTx = append(selectedTx, tx)
			err = k.removeFromUnbatchedTXIndex(ctx, *tx.Erc20Fee, txID)
			return err != nil || len(selectedTx) == maxElements
		} else {
			// we found a nil, exit
			return true
		}
	})
	return selectedTx, err
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, tokenContract string, nonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshal(bz, &b)
	for _, tx := range b.Transactions {
		tx.Erc20Token.Contract = tokenContract
		tx.Erc20Fee.Contract = tokenContract
	}
	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, tokenContract string, batchNonce uint64) error {
	batch := k.GetOutgoingTXBatch(ctx, tokenContract, batchNonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		tx.Erc20Fee.Contract = tokenContract
		k.prependToUnbatchedTXIndex(ctx, tokenContract, *tx.Erc20Fee, tx.Id)
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *batch)

	if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
		k.GetBridgeChainID(ctx) // gas used
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(batchNonce)),
	))
	return nil
}

// IterateOutgoingTXBatches iterates through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(key []byte, batch *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshal(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), &batch) {
			break
		}
	}
}

// GetOutgoingTxBatches returns the outgoing tx batches
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	batches := k.GetOutgoingTxBatches(ctx)
	var lastBatch *types.OutgoingTxBatch = nil
	lastNonce := uint64(0)
	for _, batch := range batches {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
	}
	return lastBatch
}

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBatchBlock, types.UInt64Bytes(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedBatchBlock)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetUnSlashedBatches returns all the unslashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (out []*types.OutgoingTxBatch) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)
	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock, maxHeight, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		if batch.Block > lastSlashedBatchBlock {
			out = append(out, batch)
		}
		return false
	})
	return
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(ctx sdk.Context, lastSlashedBatchBlock uint64, maxHeight uint64, cb func([]byte, *types.OutgoingTxBatch) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchBlockKey)
	iter := prefixStore.Iterator(types.UInt64Bytes(lastSlashedBatchBlock), types.UInt64Bytes(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var Batch types.OutgoingTxBatch
		k.cdc.MustUnmarshal(iter.Value(), &Batch)
		// cb returns true to stop early
		if cb(iter.Key(), &Batch) {
			break
		}
	}
}
