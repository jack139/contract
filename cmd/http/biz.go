package http

import (
	"github.com/jack139/contract/cmd/ipfs"
	cmdclient "github.com/jack139/contract/cmd/client"
	"github.com/jack139/contract/x/contract/types"

	//"strconv"
	"log"
	"bytes"
	"strings"
	"encoding/json"
	"encoding/hex"
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
	address, mnemonic, err := cmdclient.AddUserAccount(HttpCmd, userName, types.RewardRegister)
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
	doContractDelivery(ctx, types.ActionContract)
}


/* 验收 */
func bizDelivery(ctx *fasthttp.RequestCtx) {
	log.Println("biz_delivery")
	doContractDelivery(ctx, types.ActionDelivery)
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


	// 获取 ctx 上下文
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9009, err.Error())
		return
	}

	// 检查 用户地址 是否存在
	_, err = fetchKey(clientCtx.Keyring, pubkeyA)
	if err != nil {
		respError(ctx, 9021, "invalid userkeyA")
		return
	}
	_, err = fetchKey(clientCtx.Keyring, pubkeyB)
	if err != nil {
		respError(ctx, 9021, "invalid userkeyB")
		return
	}


	// data 存 ipfs
	var cid string
	if len(data)>0 {
		cid, err = ipfs.Add([]byte(data))
		if err!=nil {
			respError(ctx, 9013, err.Error())
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
	msg := types.NewMsgCreateContract(clientCtx.GetFromAddress().String(), 
						assetsId, pubkeyA, pubkeyB, action, string(loadBytes))
	if err := msg.ValidateBasic(); err != nil {
		respError(ctx, 9010, err.Error())
		return
	}

	// 设置 接收输出
	buf := new(bytes.Buffer)
	clientCtx.Output = buf

	err = tx.GenerateOrBroadcastTxCLI(clientCtx, HttpCmd.Flags(), msg)
	if err != nil {
		respError(ctx, 9011, err.Error())
		return		
	}

	// 结果输出
	respBytes := []byte(buf.String())

	log.Println("output: ", buf.String())

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9012, err.Error())
		return
	}

	// code==0 提交成功
	if respData["code"].(float64)!=0 { 
		respError(ctx, 9099, buf.String())  ///  提交失败
		return
	}

	//fmt.Printf("%s %s\n", respData["height"], respData["data"])

	// 从 data 中解析出 id
	// >>> bytearray.fromhex("0A170A0E437265617465436F6E7472616374120569643A3137").decode()
	// '\n\x17\n\x0eCreateContract\x12\x05id:17'

    bs, err := hex.DecodeString(respData["data"].(string))
    if err != nil {
		respError(ctx, 9013, err.Error())
		return
    }

	slice1 := strings.Split(string(bs), "\n")
	//log.Println(slice1)
	slice2 := strings.Split(slice1[2], "\x12")
	//log.Println(slice2)
	slice3 := strings.Split(slice2[1], ":")
	log.Println("new: ", slice3)


	// 返回区块 信息 
	resp := map[string] interface{} {
		"block_a" : slice3[1],  // 兼容旧接口
		"block_b" : slice3[1],  // 兼容旧接口
		"height" : respData["height"].(string),
	}

	respJson(ctx, &resp)

}

