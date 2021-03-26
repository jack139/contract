package http

import (
	"github.com/jack139/contract/cmd/ipfs"
	cmdclient "github.com/jack139/contract/cmd/client"
	"github.com/jack139/contract/x/contract/types"

	//"strconv"
	"log"
	"encoding/json"
	//"encoding/base64"
	"github.com/valyala/fasthttp"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

/* 企业链业务处理 */

/* 用户注册 
	action == 13
*/


func bizRegister(ctx *fasthttp.RequestCtx) {
	log.Println("biz_register")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	userName, ok := (*reqData)["user_name"].(string)
	if !ok {
		respError(ctx, 9001, "need user_name")
		return
	}
	//userType, ok := (*reqData)["user_type"].(string)
	//if !ok {
	//	respError(ctx, 9002, "need user_type")
	//	return
	//}
	//referrer, _ := (*reqData)["referrer"].(string)

	// 生成新用户密钥
	address, mnemonic, err := cmdclient.AddUserAccount(HttpCmd, userName)
	if err != nil{
		respError(ctx, 9009, err.Error())
		return		
	}

	// 返回区块id
	resp := map[string] interface{} {
		"block"   : map[string]interface{}{"id" : ""}, // 为了兼容旧接口，目前无数据返回
		"userkey" : address,
		"mnemonic" : mnemonic,
	}

	respJson(ctx, &resp)
}



/* 签合同 */
func bizContract(ctx *fasthttp.RequestCtx) {
	log.Println("biz_contract")
	doContractDelivery(ctx, "11")
}


/* 验收 */
func bizDelivery(ctx *fasthttp.RequestCtx) {
	log.Println("biz_delivery")
	doContractDelivery(ctx, "12")
}


/*  目前签合同和验收的操作一样，用同一个实现
	action： 11 前合同  12 验收 
*/


func doContractDelivery(ctx *fasthttp.RequestCtx, action string) {
	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkeyA, ok := (*reqData)["userkey_a"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey_a")
		return
	}
	pubkeyB, ok := (*reqData)["userkey_b"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey_b")
		return
	}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9001, "need assets_id")
		return
	}
	data, ok := (*reqData)["data"].(string)
	if !ok {
		respError(ctx, 9002, "need data")
		return
	}
	if len(data)>5242880 {
		respError(ctx, 9003, "data too large: over 5M")
		return		
	}

	/*
	// 获取用户密钥
	meA, ok := SECRET_KEY[pubkeyA]
	if !ok {
		respError(ctx, 9011, "wrong userkey_a")
		return
	}

	// 获取用户密钥
	meB, ok := SECRET_KEY[pubkeyB]
	if !ok {
		respError(ctx, 9011, "wrong userkey_b")
		return
	}
	*/

	// data 存 ipfs
	var cid string
	if len(data)>0 {
		cid, err = ipfs.Add([]byte(data))
		if err!=nil {
			respError(ctx, 9012, err.Error())
			return
		}
	} else {
		cid = ""
	}

	// 准备数据
	var loadData = map[string]interface{}{
		"image" : cid,
	}
	loadBytes, err := json.Marshal(loadData)
	if err != nil {
		respError(ctx, 9008, err.Error())
		return
	}

	// 数据上链
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9009, err.Error())
		return
	}

	msg := types.NewMsgCreateContract(clientCtx.GetFromAddress().String(), 
						assetsId, pubkeyA, pubkeyB, action, string(loadBytes))
	if err := msg.ValidateBasic(); err != nil {
		respError(ctx, 9010, err.Error())
		return
	}
	err = tx.GenerateOrBroadcastTxCLI(clientCtx, HttpCmd.Flags(), msg)
	if err != nil {
		respError(ctx, 9011, err.Error())
		return		
	}


	/*
	// 提交交易, A B 两个用户都提交
	respBytesA, err := meA.Deal(strconv.Itoa(action), assetsId, string(loadBytes), "") 
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}
	respBytesB, err := meB.Deal(strconv.Itoa(action), assetsId, string(loadBytes), "") 
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respDataA map[string]interface{}
	var respDataB map[string]interface{}

	if err := json.Unmarshal(respBytesA, &respDataA); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}
	if err := json.Unmarshal(respBytesB, &respDataB); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}
	*/

	// 返回两个区块id
	resp := map[string] interface{} {
		"block_a" : "",
		"block_b" : "",
	}

	respJson(ctx, &resp)

}

