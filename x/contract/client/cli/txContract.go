package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/jack139/contract/x/contract/types"
)

func CmdCreateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-contract [contractNo] [partyA] [partyB] [action] [data]",
		Short: "Creates a new contract",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsContractNo := string(args[0])
			argsPartyA := string(args[1])
			argsPartyB := string(args[2])
			argsAction := string(args[3])
			argsData := string(args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateContract(clientCtx.GetFromAddress().String(), string(argsContractNo), string(argsPartyA), string(argsPartyB), string(argsAction), string(argsData))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-contract [id] [contractNo] [partyA] [partyB] [action] [data]",
		Short: "Update a contract",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			argsContractNo := string(args[1])
			argsPartyA := string(args[2])
			argsPartyB := string(args[3])
			argsAction := string(args[4])
			argsData := string(args[5])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateContract(clientCtx.GetFromAddress().String(), id, string(argsContractNo), string(argsPartyA), string(argsPartyB), string(argsAction), string(argsData))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-contract [id] [contractNo] [partyA] [partyB] [action] [data]",
		Short: "Delete a contract by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteContract(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
