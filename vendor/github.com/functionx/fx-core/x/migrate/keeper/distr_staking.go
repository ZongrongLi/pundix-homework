package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/x/migrate/types"
)

type DistrStakingMigrate struct {
	distrKey      sdk.StoreKey
	stakingKey    sdk.StoreKey
	stakingKeeper types.StakingKeeper
}

func NewDistrStakingMigrate(distrKey, stakingKey sdk.StoreKey, stakingKeeper types.StakingKeeper) MigrateI {
	return &DistrStakingMigrate{
		distrKey:      distrKey,
		stakingKey:    stakingKey,
		stakingKeeper: stakingKeeper,
	}
}

func (m *DistrStakingMigrate) Validate(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error {
	//check validator
	if _, found := m.stakingKeeper.GetValidator(ctx, sdk.ValAddress(from)); found {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, %s is the validator address", from.String())
	}
	if _, found := m.stakingKeeper.GetValidator(ctx, sdk.ValAddress(to)); found {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, %s is the validator address", to.String())
	}
	//check delegation
	if delegations := m.stakingKeeper.GetDelegatorDelegations(ctx, to, 1); len(delegations) > 0 {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, address %s has delegation record", to.String())
	}
	//check undelegatetion
	undelegations := m.stakingKeeper.GetUnbondingDelegations(ctx, to, 1)
	if len(undelegations) > 0 {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, address %s has undelegate record", to.String())
	}
	//check redelegation
	redelegations := m.stakingKeeper.GetRedelegations(ctx, to, 1)
	if len(redelegations) > 0 {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, address %s has redelegation record", to.String())
	}
	return nil
}

func (m *DistrStakingMigrate) Execute(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error {
	stakingStore := ctx.KVStore(m.stakingKey)
	distrStore := ctx.KVStore(m.distrKey)

	//migrate distribution withdraw address
	//if bz := distrStore.Get(distrtypes.GetDelegatorWithdrawAddrKey(from)); bz != nil {
	//	distrStore.Delete(distrtypes.GetDelegatorWithdrawAddrKey(from))
	//	distrStore.Set(distrtypes.GetDelegatorWithdrawAddrKey(to), bz)
	//}

	//migrate delegate info
	delegateIterator := sdk.KVStorePrefixIterator(stakingStore, stakingtypes.GetDelegationsKey(from))
	defer delegateIterator.Close()
	for ; delegateIterator.Valid(); delegateIterator.Next() {
		info := stakingtypes.MustUnmarshalDelegation(k.cdc, delegateIterator.Value())

		//distribution starting info
		key := distrtypes.GetDelegatorStartingInfoKey(info.GetValidatorAddr(), from)
		startingInfo := distrStore.Get(key)
		distrStore.Delete(key)
		distrStore.Set(distrtypes.GetDelegatorStartingInfoKey(info.GetValidatorAddr(), to), startingInfo)

		//staking delegate
		info.DelegatorAddress = to.String()
		stakingStore.Delete(delegateIterator.Key())
		stakingStore.Set(stakingtypes.GetDelegationKey(to, info.GetValidatorAddr()), stakingtypes.MustMarshalDelegation(k.cdc, info))
	}

	//migrate unbonding delegation
	unbondingDelegationIterator := sdk.KVStorePrefixIterator(stakingStore, stakingtypes.GetUBDsKey(from))
	defer unbondingDelegationIterator.Close()
	for ; unbondingDelegationIterator.Valid(); unbondingDelegationIterator.Next() {
		ubd := stakingtypes.MustUnmarshalUBD(k.cdc, unbondingDelegationIterator.Value())
		ubd.DelegatorAddress = to.String()

		valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		stakingStore.Delete(unbondingDelegationIterator.Key())
		stakingStore.Set(stakingtypes.GetUBDKey(to, valAddr), stakingtypes.MustMarshalUBD(k.cdc, ubd))

		stakingStore.Delete(stakingtypes.GetUBDByValIndexKey(from, valAddr))
		stakingStore.Set(stakingtypes.GetUBDByValIndexKey(to, valAddr), []byte{})

		//migrate unbonding queue
		for _, entry := range ubd.Entries {
			var ubdFlag bool
			UBDQueue := m.stakingKeeper.GetUBDQueueTimeSlice(ctx, entry.CompletionTime)
			for i := range UBDQueue {
				if UBDQueue[i].DelegatorAddress == from.String() {
					UBDQueue[i].DelegatorAddress = to.String()
					ubdFlag = true
				}
			}
			if ubdFlag {
				key := stakingtypes.GetUnbondingDelegationTimeKey(entry.CompletionTime)
				value := k.cdc.MustMarshal(&stakingtypes.DVPairs{Pairs: UBDQueue})
				stakingStore.Set(key, value)
			}
		}
	}

	//migrate redelegate
	redelegateIterator := sdk.KVStorePrefixIterator(stakingStore, stakingtypes.GetREDsKey(from))
	defer redelegateIterator.Close()
	for ; redelegateIterator.Valid(); redelegateIterator.Next() {
		red := stakingtypes.MustUnmarshalRED(k.cdc, redelegateIterator.Value())
		red.DelegatorAddress = to.String()

		valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		stakingStore.Delete(redelegateIterator.Key())
		stakingStore.Set(stakingtypes.GetREDKey(to, valSrcAddr, valDstAddr), stakingtypes.MustMarshalRED(k.cdc, red))

		stakingStore.Delete(stakingtypes.GetREDByValSrcIndexKey(from, valSrcAddr, valDstAddr))
		stakingStore.Set(stakingtypes.GetREDByValSrcIndexKey(to, valSrcAddr, valDstAddr), []byte{})

		stakingStore.Delete(stakingtypes.GetREDByValDstIndexKey(from, valSrcAddr, valDstAddr))
		stakingStore.Set(stakingtypes.GetREDByValDstIndexKey(to, valSrcAddr, valDstAddr), []byte{})

		//migrate redelegate queue
		for _, entry := range red.Entries {
			var redFlag bool
			redQueue := m.stakingKeeper.GetRedelegationQueueTimeSlice(ctx, entry.CompletionTime)
			for i := range redQueue {
				if redQueue[i].DelegatorAddress == from.String() {
					redQueue[i].DelegatorAddress = to.String()
					redFlag = true
				}
			}
			if redFlag {
				key := stakingtypes.GetRedelegationTimeKey(entry.CompletionTime)
				value := k.cdc.MustMarshal(&stakingtypes.DVVTriplets{Triplets: redQueue})
				stakingStore.Set(key, value)
			}
		}
	}
	return nil
}
