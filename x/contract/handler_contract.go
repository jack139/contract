package contract

import (
	"fmt"
	"strconv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jack139/contract/x/contract/keeper"
	"github.com/jack139/contract/x/contract/types"
)

func handleMsgCreateContract(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateContract) (*sdk.Result, error) {
	id := k.CreateContract(ctx, *msg)

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(), 
		Data: []byte("id:"+strconv.FormatInt(id, 10)),
	}, nil
}

func handleMsgUpdateContract(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateContract) (*sdk.Result, error) {
	var contract = types.Contract{
		Creator:    msg.Creator,
		Id:         msg.Id,
		ContractNo: msg.ContractNo,
		PartyA:     msg.PartyA,
		PartyB:     msg.PartyB,
		Action:     msg.Action,
		Data:       msg.Data,
	}

	// Checks that the element exists
	if !k.HasContract(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetContractOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetContract(ctx, contract)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteContract(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteContract) (*sdk.Result, error) {
	if !k.HasContract(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetContractOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.DeleteContract(ctx, msg.Id)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
