package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/client/tx"
)

const (
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			vestingStart, err := cmd.Flags().GetInt64(flagVestingStart)
			if err != nil {
				return err
			}
			vestingEnd, err := cmd.Flags().GetInt64(flagVestingEnd)
			if err != nil {
				return err
			}
			vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
			if err != nil {
				return err
			}

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)

				if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
					baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
					return errors.New("vesting amount cannot be greater than total amount")
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)

				default:
					return errors.New("invalid vesting parameters; must supply start and end time or end time")
				}
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(depCdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Int64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Int64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}




// 添加用户： add key and account
func AddUserCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-user [user_name]",
		Short: "Add a user to chain",
		Long: `Add a user to chain. `,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return AddUserAccount(cmd, args[0])
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}


// 建新用户 user， 建key，建account
func AddUserAccount(cmd *cobra.Command, name string) error {
	
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	// 获取 keyring 环境
	var kb keyring.Keyring

	buf := bufio.NewReader(cmd.InOrStdin())
	keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
	if err != nil {
		return err
	}
	kb, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.KeyringDir, buf)

	// 获取 user0的地址
	keyref := "user0"
	info0, err := kb.Key(keyref)
	if err != nil {
		return err
	}
	//addr0 := info0.GetAddress()

	// 参考cosmos-sdk/client/keys/show.go 中 getBechKeyOut()
	ko, err := keyring.Bech32KeyOutput(info0)  
	if err != nil {
		return err
	}

	// 取得地址字符串： 例如 contract1zfqgxtujvpy92prtzgmzs3ygta9y2cl3w8hxlh
	addr0 := ko.Address
	//fmt.Println(addr0)

	cmd.Flags().Set(flags.FlagFrom, addr0)
	clientCtx, err = client.GetClientTxContext(cmd) // 设置了addr0, 重新获取一次context
	if err != nil {
		return err
	}

	// 注册新的 key
	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return err
	}

	hdPath := hd.CreateHDPath(sdk.GetConfig().GetCoinType(), 0, 0).String()

	// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
	mnemonicEntropySize := 256
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return err
	}

	// Get bip39 mnemonic
	var mnemonic, bip39Passphrase string

	mnemonic, err = bip39.NewMnemonic(entropySeed)
	if err != nil {
		return err
	}

	info, err := kb.NewAccount(name, mnemonic, bip39Passphrase, hdPath, algo)
	if err != nil {
		return err
	}

	fmt.Println("mnemonic: ", mnemonic)
	//fmt.Println(info)

	// 新用户的 地址
	toAddr := info.GetAddress()

	//fmt.Println("from ", clientCtx.GetFromAddress())
	//fmt.Println("to ", toAddr)

	// 转账 1credit， 会自动建立auth的账户
	coins, err := sdk.ParseCoinsNormalized("1credit")
	if err != nil {
		return err
	}

	msg := banktypes.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins)
	if err := msg.ValidateBasic(); err != nil {
		fmt.Print("3", err)
		return err
	}

	// 调用 send 的 RPC 服务
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

}
