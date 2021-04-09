package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jack139/contract/x/contract/types"
	"strconv"
)

// GetContractCount get the total number of contract
func (k Keeper) GetContractCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractCountKey))
	byteKey := types.KeyPrefix(types.ContractCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetContractCount set the total number of contract
func (k Keeper) SetContractCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractCountKey))
	byteKey := types.KeyPrefix(types.ContractCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateContract creates a contract with a new id and update the count
func (k Keeper) CreateContract(ctx sdk.Context, msg types.MsgCreateContract) {
	// Create the contract
	count := k.GetContractCount(ctx)
	var contract = types.Contract{
		Creator:    msg.Creator,
		Id:         strconv.FormatInt(count, 10),
		ContractNo: msg.ContractNo,
		PartyA:     msg.PartyA,
		PartyB:     msg.PartyB,
		Action:     msg.Action,
		Data:       msg.Data,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	key := types.KeyPrefix(types.ContractKey + contract.Id)
	value := k.cdc.MustMarshalBinaryBare(&contract)
	store.Set(key, value)

	// Update contract count
	k.SetContractCount(ctx, count+1)
}

// SetContract set a specific contract in the store
func (k Keeper) SetContract(ctx sdk.Context, contract types.Contract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	b := k.cdc.MustMarshalBinaryBare(&contract)
	store.Set(types.KeyPrefix(types.ContractKey+contract.Id), b)
}

// GetContract returns a contract from its id
func (k Keeper) GetContract(ctx sdk.Context, key string) types.Contract {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	var contract types.Contract
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.ContractKey+key)), &contract)
	return contract
}


// HasContract checks if the contract exists
func (k Keeper) HasContract(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	return store.Has(types.KeyPrefix(types.ContractKey + id))
}

// GetContractOwner returns the creator of the contract
func (k Keeper) GetContractOwner(ctx sdk.Context, key string) string {
	return k.GetContract(ctx, key).Creator
}

// DeleteContract deletes a contract
func (k Keeper) DeleteContract(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	store.Delete(types.KeyPrefix(types.ContractKey + key))
}

// GetAllContract returns all contract
func (k Keeper) GetAllContract(ctx sdk.Context) (msgs []types.Contract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ContractKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Contract
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}


// GetAllContract returns all contract with spicific contractNo
func (k Keeper) GetContractByNo(ctx sdk.Context, contractNo string) (msgs []types.Contract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ContractKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Contract
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		if msg.ContractNo==contractNo{
			msgs = append(msgs, msg)
		}
	}

	return
}

// GetAllContract returns all contract with spicific user
func (k Keeper) GetContractByUser(ctx sdk.Context, user string) (msgs []types.Contract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ContractKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ContractKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Contract
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		if (msg.PartyA==user) || (msg.PartyB==user) {
			msgs = append(msgs, msg)
		} 
	}

	return
}
