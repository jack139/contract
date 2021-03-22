package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/jack139/contract/x/contract/types"
	"github.com/spf13/cobra"
)

func CmdListContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllContractRequest{
			}

			res, err := queryClient.ContractAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [id]",
		Short: "shows a contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetContractRequest{
				Id: args[0],
			}

			res, err := queryClient.Contract(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowByNoContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-by-no [contractNo]",
		Short: "shows a contract by contract No.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetContractByNoRequest{
				ContractNo: args[0],
			}

			res, err := queryClient.ContractByNo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
