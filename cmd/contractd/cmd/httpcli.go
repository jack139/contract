package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/flags"

	httpcli "github.com/jack139/contract/cmd/http"
)

func HttpCliCmd() *cobra.Command {
	cmd := &cobra.Command{	// 启动http服务
		Use:   "http <port>",
		Short: "start http service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need port number")
			}

			// 保存 cmd
			httpcli.HttpCmd = cmd

			httpcli.RunServer(args[0])
			// 不会返回
			return nil
		},
	}

	cmd.Flags().String(flags.FlagChainID, "", "network chain ID")
	cmd.Flags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.Flags().BoolP(flags.FlagSkipConfirmation, "y", true, "Skip tx broadcasting prompt confirmation")

	return cmd
}