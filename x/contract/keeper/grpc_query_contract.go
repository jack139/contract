package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jack139/contract/x/contract/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*  具体查询使用contract.go中定义，算法同源 */


func (k Keeper) ContractAll(c context.Context, req *types.QueryAllContractRequest) (*types.QueryAllContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contracts []*types.Contract
	ctx := sdk.UnwrapSDKContext(c)

	r := k.GetAllContract(ctx)
	for i, _ := range r{
		contracts = append(contracts, &r[i])
	}

	return &types.QueryAllContractResponse{Contract: contracts}, nil
}


func (k Keeper) Contract(c context.Context, req *types.QueryGetContractRequest) (*types.QueryGetContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contract types.Contract
	ctx := sdk.UnwrapSDKContext(c)

	contract = k.GetContract(ctx, req.Id)

	return &types.QueryGetContractResponse{Contract: &contract}, nil
}



func (k Keeper) ContractByNo(c context.Context, req *types.QueryGetContractByNoRequest) (*types.QueryGetContractByNoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contracts []*types.Contract
	ctx := sdk.UnwrapSDKContext(c)

	r := k.GetContractByNo(ctx, req.ContractNo)
	for i, _ := range r{
		contracts = append(contracts, &r[i])
	}

	return &types.QueryGetContractByNoResponse{Contract: contracts}, nil
}

func (k Keeper) ContractByUser(c context.Context, req *types.QueryGetContractByUserRequest) (*types.QueryGetContractByUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contracts []*types.Contract
	ctx := sdk.UnwrapSDKContext(c)

	r := k.GetContractByUser(ctx, req.User)
	for i, _ := range r{
		contracts = append(contracts, &r[i])
	}

	return &types.QueryGetContractByUserResponse{Contract: contracts}, nil
}