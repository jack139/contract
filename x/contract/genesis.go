package contract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jack139/contract/x/contract/keeper"
	"github.com/jack139/contract/x/contract/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the contract
	for _, elem := range genState.ContractList {
		k.SetContract(ctx, *elem)
	}

	// Set contract count
	k.SetContractCount(ctx, int64(len(genState.ContractList)))

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all contract
	contractList := k.GetAllContract(ctx)
	for _, elem := range contractList {
		elem := elem
		genesis.ContractList = append(genesis.ContractList, &elem)
	}

	return genesis
}
