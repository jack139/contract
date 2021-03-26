package client

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/tx"
)


// 建新用户 user， 建key，建account
// 返回： address, mnemonic
func AddUserAccount(cmd *cobra.Command, name string) (string, string, error) {
	
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return "", "", err
	}

	// 获取 keyring 环境
	var kb keyring.Keyring

	buf := bufio.NewReader(cmd.InOrStdin())
	keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
	if err != nil {
		return "", "", err
	}
	kb, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.KeyringDir, buf)

	// 获取 user0的地址
	keyref := "faucet"
	info0, err := kb.Key(keyref)
	if err != nil {
		return "", "", err
	}
	//addr0 := info0.GetAddress()

	// 参考cosmos-sdk/client/keys/show.go 中 getBechKeyOut()
	ko, err := keyring.Bech32KeyOutput(info0)  
	if err != nil {
		return "", "", err
	}

	// 取得地址字符串： 例如 contract1zfqgxtujvpy92prtzgmzs3ygta9y2cl3w8hxlh
	addr0 := ko.Address
	//fmt.Println(addr0)

	cmd.Flags().Set(flags.FlagFrom, addr0)
	clientCtx, err = client.GetClientTxContext(cmd) // 设置了addr0, 重新获取一次context
	if err != nil {
		return "", "", err
	}

	// 注册新的 key
	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return "", "", err
	}

	hdPath := hd.CreateHDPath(sdk.GetConfig().GetCoinType(), 0, 0).String()

	// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
	mnemonicEntropySize := 256
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return "", "", err
	}

	// Get bip39 mnemonic
	var mnemonic, bip39Passphrase string

	mnemonic, err = bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", "", err
	}

	info, err := kb.NewAccount(name, mnemonic, bip39Passphrase, hdPath, algo)
	if err != nil {
		return "", "", err
	}

	fmt.Println("mnemonic: ", mnemonic)
	//fmt.Println(info)

	// 新用户的 地址
	toAddr := info.GetAddress()

	//fmt.Println("from ", clientCtx.GetFromAddress())
	//fmt.Println("to ", toAddr)

	// 转账 1credit， 会自动建立auth的账户
	coins, err := sdk.ParseCoinsNormalized("20credit")
	if err != nil {
		return "", "", err
	}

	msg := banktypes.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins)
	if err := msg.ValidateBasic(); err != nil {
		return "", "", err
	}


	// 参考cosmos-sdk/client/keys/show.go 中 getBechKeyOut()
	ko_new, err := keyring.Bech32KeyOutput(info)  
	if err != nil {
		return "", "", err
	}

	// 取得地址字符串： 例如 contract1zfqgxtujvpy92prtzgmzs3ygta9y2cl3w8hxlh
	addr_new := ko_new.Address

	// 调用 send 的 RPC 服务
	return addr_new, mnemonic, tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

}
