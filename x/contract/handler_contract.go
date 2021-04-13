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

	// 通证奖励
	var reward string
	switch msg.Action {
		case types.ActionContract:
			reward = types.RewardContract
		case types.ActionDelivery:
			reward = types.RewardDelivery
		default:
			reward = "0credit"
	}

	// 生成 faucet 地址
	faucetAcct, err := sdk.AccAddressFromBech32(types.FaucetAddress)
	if err != nil {
		return nil, err
	}
	userAcctA, err := sdk.AccAddressFromBech32(msg.PartyA)
	if err != nil {
		return nil, err
	}
	userAcctB, err := sdk.AccAddressFromBech32(msg.PartyB)
	if err != nil {
		return nil, err
	}
	// 生成金额
	payment, _ := sdk.ParseCoinsNormalized(reward)
	// 转账
	if err := k.CoinKeeper.SendCoins(ctx, faucetAcct, userAcctA, payment); err != nil {
		return nil, err
	}
	if err := k.CoinKeeper.SendCoins(ctx, faucetAcct, userAcctB, payment); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(), 
		Data: []byte("id:"+strconv.FormatInt(id, 10)), // id 作为data返回
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
