package http

import (
	"github.com/jack139/contract/cmd/ipfs"
	"github.com/jack139/contract/x/contract/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"

	"strings"
	"log"
	"bytes"
	"encoding/json"
	"context"
	"strconv"
	"github.com/valyala/fasthttp"
)

/*
	{
 		"Contract":[
 			{
 				"creator":"contract1ycd5htx4emg686w05rhpfpqj5hv35g7sy2qmae",
 				"id":"0",
 				"contractNo":"123",
 				"partyA":"contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx",
 				"partyB":"contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx",
 				"action":"\u000b",
 				"data":"{\"image\":\"QmUaLpY8SX68Cop5UDZbYGpaDu3JCK7euFsnDzw5xpN1sz\"}"
 			}
 		]
 	}

 							#
 							#
 						    V

 	[
		{
			"id":"0",
			"assets_id":"123",
			"exchange_id":"contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx",
			"action":11,
			"image":"zzzzzzz",
			"type":"DEAL",
			"refer":"",
		}
 	]
*/

/* 处理一个交易的数据，为了兼容旧的字段名 */
func processData(item0 *map[string]interface{}, user string) (*map[string]interface{}, error) {
	item := *item0

	var data map[string]interface{}
	// 检查query用户是否是相关者
	if (item["partyA"]!=user) && (item["partyB"]!=user){
		// 不相关，不返回data数据
		data = make(map[string]interface{})
		data["image"] = ""
	} else {
		// 相关，解析 data内容			
		if err := json.Unmarshal([]byte(item["data"].(string)), &data); err != nil {
			return nil, err
		}
		
		// 处理image 字段，从ipfs读取
		_, ok := data["image"]
		if ok && len(data["image"].(string))>0 {
			image_data, err := ipfs.Get(data["image"].(string))
			if err!=nil {
				return nil, err
			}
			data["image"] = string(image_data)
		}		
	}

	// 建立返回的数据
	new_item := map[string] interface{} {
		"id": item["id"],
		"exchange_id": user,
		"userkey_a": item["partyA"],
		"userkey_b": item["partyB"],
		"assets_id": item["contractNo"],
		"action": item["action"],
		"type": "DEAL",
		"refer": "",
		"data": data,
	}

	return &new_item, nil
}

/* data字段是已序列化的json串，反序列化一下， 只对一个数据 */
func unmarshalData(reqData *map[string]interface{}, user string) (*map[string]interface{}, error) {
	var respData map[string]interface{}

	data := (*reqData)["Contract"].(map[string]interface{})

	// 检查 data 字段是否正常
	_, ok := data["data"]
	if !ok {
		return &respData, nil
	}
	if !strings.HasPrefix(data["data"].(string), "{") {
		return &respData, nil
	}

	return processData(&data, user)
}

/* data字段是已序列化的json串，反序列化一下，针对数据列表 */
func unmarshalDataList(reqData *map[string]interface{}, user string) (*[]map[string]interface{}, error) {
	var respData []map[string]interface{}

	//log.Printf("id: %v\n", ((*reqData)["Contract"].([]interface{})[0]).(map[string]interface{})["id"])

	dataList := (*reqData)["Contract"].([]interface{})

	// 处理data字段
	for _, item0 := range dataList {
		item := item0.(map[string]interface{})

		// 检查 data 字段是否正常
		_, ok := item["data"]
		if !ok {
			continue
		}
		if !strings.HasPrefix(item["data"].(string), "{") {
			continue
		}

		new_item, err := processData(&item, user)
		if err!=nil {
			return nil, err
		}

		respData = append(respData, *new_item)
	}
	return &respData, nil
}


/* 查询交易， 只允许查询自己的 */
func queryDeals(ctx *fasthttp.RequestCtx) {
	log.Println("query_deals")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}

	// 准备环境
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9002, err.Error())
		return
	}

	// 检查 用户地址 是否存在
	_, err = fetchKey(clientCtx.Keyring, pubkey)
	if err != nil {
		respError(ctx, 9001, "invalid userkey")
		return
	}

	// 准备查询
	queryClient := types.NewQueryClient(clientCtx)

	params := &types.QueryGetContractByUserRequest{
		User: pubkey,
	}

	res, err := queryClient.ContractByUser(context.Background(), params)
	if err != nil {
		respError(ctx, 9003, err.Error())
		return
	}

	//log.Printf("%t\n", res)

	// 设置 接收输出
	buf := new(bytes.Buffer)
	clientCtx.Output = buf

	// 转换输出
	clientCtx.PrintProto(res)

	// 输出的字节流
	respBytes := []byte(buf.String())

	log.Println("output: ", buf.String())

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	respData2, err := unmarshalDataList(&respData, pubkey)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return
	}

	resp := map[string] interface{} {
		"deals" : *respData2,
	}

	respJson(ctx, &resp)
}



/* 按资产id查询交易 */
func queryByAsstes(ctx *fasthttp.RequestCtx) {
	log.Println("query_by_assets")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9008, "need assets_id")
		return
	}

	// 准备环境
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9002, err.Error())
		return
	}

	// 检查 用户地址 是否存在
	_, err = fetchKey(clientCtx.Keyring, pubkey)
	if err != nil {
		respError(ctx, 9001, "invalid userkey")
		return
	}

	// 准备查询
	queryClient := types.NewQueryClient(clientCtx)

	params := &types.QueryGetContractByNoRequest{
		ContractNo: assetsId,
	}

	res, err := queryClient.ContractByNo(context.Background(), params)
	if err != nil {
		respError(ctx, 9003, err.Error())
		return
	}

	//log.Printf("%t\n", res)

	// 设置 接收输出
	buf := new(bytes.Buffer)
	clientCtx.Output = buf

	// 转换输出
	clientCtx.PrintProto(res)

	// 输出的字节流
	respBytes := []byte(buf.String())

	log.Println("output: ", buf.String())

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	respData2, err := unmarshalDataList(&respData, pubkey)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return
	}

	resp := map[string] interface{} {
		"deals" : *respData2,
	}

	respJson(ctx, &resp)
}


/* 指定id查询交易 */
func queryBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	// 准备环境
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9005, err.Error())
		return
	}

	// 检查 用户地址 是否存在
	_, err = fetchKey(clientCtx.Keyring, pubkey)
	if err != nil {
		respError(ctx, 9001, "invalid userkey")
		return
	}

	// 准备查询
	queryClient := types.NewQueryClient(clientCtx)

	params := &types.QueryGetContractRequest{
		Id: blockId,
	}

	res, err := queryClient.Contract(context.Background(), params)
	if err != nil {
		respError(ctx, 9003, err.Error())
		return
	}

	//log.Printf("%t\n", res)

	// 设置 接收输出
	buf := new(bytes.Buffer)
	clientCtx.Output = buf

	// 转换输出
	clientCtx.PrintProto(res)

	// 输出的字节流
	respBytes := []byte(buf.String())

	log.Println("output: ", buf.String())

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	respData2, err := unmarshalData(&respData, pubkey)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return
	}


	resp := map[string] interface{} {
		"deals" : *respData2,
	}

	respJson(ctx, &resp)
}




/* 获取区块数据 */
func getBlock(clientCtx client.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	return legacy.Cdc.MarshalJSON(res)
}

/* 指定区块查询交易 */
func queryRawBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_raw_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	height, ok := (*reqData)["height"].(string)
	if !ok {
		respError(ctx, 9002, "need height")
		return
	}

	height64, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		respError(ctx, 9007, err.Error())
		return		
	}

	// 获取 ctx 上下文
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9005, err.Error())
		return
	}

	// 检查 用户地址 是否存在
	_, err = fetchKey(clientCtx.Keyring, pubkey)
	if err != nil {
		respError(ctx, 9001, "invalid userkey")
		return
	}

	// 准备查询
	respBytes, err := getBlock(clientCtx, &height64)
	if err != nil {
		respError(ctx, 9006, err.Error())
		return
	}

	//log.Printf("%v\n", string(respBytes))

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}
	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}
