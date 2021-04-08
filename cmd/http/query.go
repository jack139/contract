package http

import (
	"github.com/jack139/contract/cmd/ipfs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/jack139/contract/x/contract/types"

	"strings"
	"log"
	"bytes"
	"encoding/json"
	"context"
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

/* data字段是已序列化的json串，反序列化一下 */
func unmarshalData(reqData *map[string]interface{}) ([]map[string]interface{}, error) {
	var respData []map[string]interface{}

	//log.Printf("id: %v\n", ((*reqData)["Contract"].([]interface{})[0]).(map[string]interface{})["id"])

	dataList := (*reqData)["Contract"].([]interface{})

	// 处理data字段
	for _, item0 := range dataList {
		item := item0.(map[string]interface{})
		_, ok := item["data"]
		if !ok {
			continue
		}
		if !strings.HasPrefix(item["data"].(string), "{") {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(item["data"].(string)), &data); err != nil {
			return nil, err
		}
		
		// 处理image 字段，从ipfs读取
		_, ok = data["image"]
		if ok && len(data["image"].(string))>0 {
			image_data, err := ipfs.Get(data["image"].(string))
			if err!=nil {
				return nil, err
			}
			data["image"] = string(image_data)
		}		

		//item["data"] = data

		new_item := map[string] interface{} {
			"id": item["id"],
			"exchange_id": item["partyA"],
			"assets_id": item["contractNo"],
			"action": item["action"],
			"type": "DEAL",
			"refer": "",
			"data": data,
		}

		respData = append(respData, new_item)
	}
	return respData, nil
}


/* 查询交易， 只允许查询自己的 */
/*
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

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 只查询当前用户的交易
	respBytes, err := me.Query("deal", "_")
	if err!=nil {
		respError(ctx, 9001, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	err = unmarshalData(&respData)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	// 返回结果
	resp := map[string] interface{} {
		"deals" : respData,
	}

	respJson(ctx, &resp)
}
*/


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
	//pubkey, ok := (*reqData)["userkey"].(string)
	//if !ok {
	//	respError(ctx, 9009, "need userkey")
	//	return
	//}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9001, "need assets_id")
		return
	}

	// 准备查询
	clientCtx, err := client.GetClientTxContext(HttpCmd)
	if err != nil {
		respError(ctx, 9002, err.Error())
		return
	}

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
	respData2, err := unmarshalData(&respData)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return
	}


	resp := map[string] interface{} {
		"deals" : respData2,
	}

	respJson(ctx, &resp)
}


/* 指定区块查询交易 */
/*
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

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	respBytes, err := me.QueryTx(pubkey, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}


	// 处理data字段
	temp := []map[string]interface{}{ respData }
	err = unmarshalData(&temp)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}
*/

/* 指定区块查询交易 */
/*
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
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	respBytes, err := me.QueryRawBlock(pubkey, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

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
*/