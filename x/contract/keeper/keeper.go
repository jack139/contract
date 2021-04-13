package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jack139/contract/x/contract/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		CoinKeeper bankkeeper.Keeper
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, coinKeeper bankkeeper.Keeper) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		CoinKeeper: coinKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
