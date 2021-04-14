##  应用层API



| 修改日期   | 修改内容                                                     |
| ---------- | ------------------------------------------------------------ |
| 2021-01-15 | 增加企业链业务处理接口 biz_*                                 |
| 2021-01-26 | 增加ipfs支持                                                 |
| 2021-03-26 | biz_register返回增加密码字符串                               |
| 2021-04-13 | query_* 查询结果返回新增内容；**query_raw_block查询输入参数改为height** |
| 2021-04-14 | 增加查询通证余额 query_balance                               |



###  一、 说明

​		应用层API与区块链节点一起部署，提供给客户端调用，进行基础的区块链功能操作。



### 二、 概念和定义

#### 1. 节点

​		节点是区块链上的一个业务处理和存储的单元，是一个具有独立处理区块链业务的服务程序。节点可以是一台物理服务器，也可以是多个节点共用一个物理服务器，通过不同端口提供各自节点的功能。

#### 2. 链用户

​		链用户是具有提交区块链交易权限的用户，线下可定义为交易所。每个链用户通过一对密钥识别（例如下例中的PubKey），同时使用此密钥进行数据的加密解密操作，因此链用户的密钥需要妥善保管。密钥类似如下形式：
```json
{
	"sign_key":{
		"type":"ed25519/privkey",
		"value":"UgM13IPx/BkwfQo8jce6TMR5bRuAv7ZLdBooTZWm2ixLaNitCW91NHW06h8VQw=="
	},
	"CryptoPair":{
		"PrivKey":"tgNfUoYkh9xKs1hVKs+5uXNetCxvDRRHBNmLMs5/NKk=",
		"PubKey":"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="
	}
}
```

#### 3. 交易区块
​		链上数据存储在区块链的区块中，区块目前分两类：（1）交易区块；（2）授权区块。交易区块用于存储买入卖出交易的交易信息和交易数据。交易区块中的部分数据是公开的，部分数据是加密的。链用户只能查看自己提交的区块上的加密数据。如果要查看其他链用户的区块加密数据，需要向区块所有者（即区块的提交者）进行请求授权。当区块所有者同意并授权后，请求方才能看到相应加密区块的数据。同时，请求和授权过程也会记录在区块链上，用于追溯。

**交易区块内容：**

| 名称       | 类型   | 说明                                       |
| ---------- | ------ | ------------------------------------------ |
| ID         | uuid   | 交易ID，自动生成                           |
| ExchangeID | string | 交易所ID（即，链用户公钥）                 |
| AssetsID   | string | 资产ID，唯一标示交易资产，由客户端定义     |
| Data       | string | 加密交易数据（只有链用户ExchangeID可解密） |
| Refer      | string | 参考数据，可用于检索                       |
| Action     | byte   | 交易类型：1 买入， 2 卖出， 3 变更所有权   |

**授权区块内容：**

| 名称           | 类型   | 说明                                   |
| -------------- | ------ | -------------------------------------- |
| ID             | uuid   | 授权ID，自动生成                       |
| ExchangeID     | string | 数据原始提交者的交易所ID（链用户公钥） |
| AuthExchangeID | string | 请求授权的交易所ID（链用户公钥）       |
| Data           | string | 加密交易数据（AuthExchangeID可以解密） |
| Action         | byte   | 交易类型：4 请求授权， 5 响应授权      |

> 说明：
>
> 1. 上述字段中，AssetsID、Data、Refer均没有长度限制，但不建议放很大的数据块
> 2. 如果需要存储大型数据，请使用IPFS存储，然后在Data字段保存IPFS的文件哈希值
> 3. AssetsID必须是可显示字符（32<ASCII<127）



### 三、 API提供的区块链功能

| 序号 | 接口名称        | 接口功能                     |
| :--: | :-------------- | ---------------------------- |
|  1   | biz_register    | 用户注册                     |
|  2   | biz_contract    | 签合同                       |
|  3   | biz_delivery    | 验收                         |
|  4   | query_deals     | 查询自己的历史交易           |
|  5   | query_by_assets | 按合同编号进行查询历史交易   |
|  6   | query_block     | 按区块ID查询指定区块         |
|  7   | query_raw_block | 按区块ID查询指定区块原始数据 |
|  8   | query_balance   | 查询通证余额                 |




### 四、接口定义

#### 1. 全局接口定义

输入参数

| 参数      | 类型   | 说明                          | 示例        |
| --------- | ------ | ----------------------------- | ----------- |
| appid | string | 应用渠道编号                  |             |
| version   | string | 版本号                        | 1 |
| sign_type | string | 签名算法，目前使用SHA256算法 | SHA256 |
| sign_data | string | 签名数据，具体算法见下文      |             |
| timestamp | int    | unix时间戳（秒）              |             |
| data      | json   | 接口数据，详见各接口定义      |             |

> 签名/验签算法：
>
> 1. appid和app_secret均从线下获得。
> 2. 筛选，获取参数键值对，剔除sign_data参数。data参数按key升序排列进行json序列化。
> 3. 排序，按key升序排序；data中json也按key升序排序。
> 4. 拼接，按排序好的顺序拼接请求参数。
>
> ```key1=value1&key2=value2&...&key=appSecret```，key=app_secret固定拼接在参数串末尾。
>
> 4. 签名，使用制定的算法进行加签获取二进制字节，使用 16进制进行编码Hex.encode得到签名串，然后base64编码。
> 5. 验签，对收到的参数按1-4步骤签名，比对得到的签名串与提交的签名串是否一致。

签名示例：

```json
请求参数：
{
    "appid": "66A095861BAE55F8735199DBC45D3E8E", 
    "version": "1", 
    "data": {
        "test1": "test1", 
        "atest2": "test2", 
        "Atest2": "test2"
    }, 
    "timestamp": 1608904438, 
    "sign_type": "SHA256",  
    "sign_data": "..."
}

密钥：
app_secret="43E554621FF7BF4756F8C1ADF17F209C"

待加签串：
appid=66A095861BAE55F8735199DBC45D3E8E&data={"Atest2":"test2","atest2":"test2","test1":"test1"}&sign_type=SHA256&timestamp=1608948188&version=1&key=43E554621FF7BF4756F8C1ADF17F209C

SHA256加签结果：
"fa72d34eafea3639b0a207bdd7ceb49586f4be92e58ee97b6453b696b0edb781"

base64后结果：
"ZmE3MmQzNGVhZmVhMzYzOWIwYTIwN2JkZDdjZWI0OTU4NmY0YmU5MmU1OGVlOTdiNjQ1M2I2OTZiMGVkYjc4MQ=="
```

返回结果

| 参数      | 类型    | 说明                                                         |
| --------- | ------- | ------------------------------------------------------------ |
| code      | int   | 状态代码，0 表示成功，非0 表示出错                                 |
| msg   | string | 成功时返回success；出错时，返回出错信息                                                     |
| data      | json    | 成功时返回结果数据，详见具体接口                |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
    }
}
```

全局出错代码

| 编码 | 说明                               |
| ---- | ---------------------------------- |
| 9000 | 签名错误                           |



#### 2. 业务处理接口

##### 2.1 注册用户

请求URL

> http://<host>:<port>/api/biz_register

请求方式

> POST

输入参数（data字段下）

| 参数      | 类型   | 说明                       |
| --------- | ------ | -------------------------- |
| user_name | string | 用户名称                   |
| user_type | string | 注册用户类型               |
| referrer  | string | 推荐人的用户公钥（可为空） |

> user_type 取值："office" 事务所；"supplier" 供应商；"buyer" 企业用户。

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 用户公钥、密码字符串                    |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "user_name": "test1", 
        "user_type": "buyer"
    }, 
    "timestamp": 1610692800, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "MTZlODRlNGYyMWNiNTk1MzAxYWUyNjI0ODIzOWQxYWI1MjZmZmQzMDc3ZDU5ZmZiMGEzMWU2Y2QwOGE1NTdhOQ=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'block': {'id': ''}, 
        'mnemonic': 'path basic oblige aware sort prefer logic program differ badge reveal effort evoke fork clown before autumn frozen unusual lottery dawn swim exercise bread', 
        'userkey': 'contract1eq3ppq7n7ukty6gm6v7pyfz2mls63x7n9q5v6r'}, 
    'msg': 'success'
}
```



##### 2.2 签合同

请求URL

> http://<host>:<port>/api/biz_contract

请求方式

> POST

输入参数（data字段下）

| 参数      | 类型   | 说明               |
| --------- | ------ | ------------------ |
| userkey_a | string | 甲方公钥           |
| userkey_b | string | 乙方公钥           |
| assets_id | string | 合同编号           |
| data      | base64 | 合同照片base64编码 |

返回结果

| 参数 | 类型             | 说明                                    |
| ---- | ---------------- | --------------------------------------- |
| code | int              | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string           | 成功时返回success；出错时，返回出错信息 |
| data | json | 交易数据id |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey_a": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex", 
        "userkey_b": "contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx", 
        "assets_id": "12345678", 
        "data": "abcdefghijklmn"
    }, 
    "timestamp": 1618284612, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "ZDYwNjdiOTI1NzBmMzIxYmM2NjZmNjc5ZTY5YjVkYzhlMzk1NjBiYTJmOTlhOGE1ZTBiY2U4ZmFlMWU3MDAxYw=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'block_a': '22', 
        'block_b': '22', 
        'height': '217319', /* 区块高度，用于查询raw block */
    }, 
    'msg': 'success'
}
```



##### 2.3 验收

请求URL

> http://<host>:<port>/api/biz_delivery

请求方式

> POST

输入参数（data字段下）

| 参数      | 类型   | 说明               |
| --------- | ------ | ------------------ |
| userkey_a | string | 甲方公钥           |
| userkey_b | string | 乙方公钥           |
| assets_id | string | 合同编号           |
| data      | base64 | 验收照片base64编码 |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 区块id                                  |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey_a": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex", 
        "userkey_b": "contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx", 
        "assets_id": "12345678", 
        "data": "abcdefghijklmn"
    }, 
    "timestamp": 1618284657, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "ZTk2ZWNjYmQ1Zjk1ZTUzZGFlYzUxODRlOTBjNTg1ZWY3ZjI4YzgxYmQ1MmUwYmRhYjUwNWJlODE4NWExZDgyNA=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'block_a': '23', 
        'block_b': '23', 
        'data_id': '23', 
        'height': '217364', /* 区块高度，用于查询raw block */
    }, 
    'msg': 'success'
}
```



#### 3. 查询接口

##### 3.1 查询所有历史交易

请求URL

> http://<host>:<port>/api/query_deals

请求方式

> POST

输入参数（data字段下）

| 参数    | 类型   | 说明     |
| ------- | ------ | -------- |
| userkey | string | 用户公钥 |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 交易列表                                |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex"
    }, 
    "timestamp": 1618284714, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "ZmI5Y2VmNzE2NjFhYjczODdiZmNlNzU1ZTUxOTA4MmFkYjk4MDI2M2VhYWMzNDkxODUxYTBmYzhmMzA0N2ZkZQ=="
}

```

返回示例

```json
{
    'code': 0, 
    'data': {
        'deals': [
            {
                'action': '11', 
                'assets_id': '1234', 
                'data': {'image': '11111111111111111111'}, 
                'exchange_id': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'id': '13', 
                'refer': '', 
                'type': 'DEAL', 
                'userkey_a': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'userkey_b': 'contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx'
            }, 
            {
                'action': '12', 
                'assets_id': '1234', 
                'data': {
                    'image': '11111111111111111111'}, 
                'exchange_id': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'id': '18', 
                'refer': '', 
                'type': 'DEAL', 
                'userkey_a': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'userkey_b': 'contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx'
            }
        ]
    }, 
    'msg': 'success'
}
```





##### 3.2 按合同编号查询历史交易

请求URL

> http://<host>:<port>/api/query_by_assets

请求方式

> POST

输入参数（data字段下）

| 参数      | 类型   | 说明     |
| --------- | ------ | -------- |
| userkey   | string | 用户公钥 |
| assets_id | string | 合同编号 |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 相同资产ID的交易列表                    |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex", 
        "assets_id": "12345678"
    }, 
    "timestamp": 1618284759, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "M2VlNDhjMjgxMmY0Y2E4YzMwZjUyYmJmZTE0MzY0YWYzNDNkZWU4MGQ5Y2Y1YTg5ZThiNjkyZDI4ZTUxMjE3OA=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'deals': [
            {
                'action': '11', 
                'assets_id': '12345678', 
                'data': {'image': 'abcdefghijklmn'}, 
                'exchange_id': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'id': '22', 
                'refer': '', 
                'type': 'DEAL', 
                'userkey_a': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'userkey_b': 'contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx'
            }, 
            {
                'action': '12', 
                'assets_id': '12345678', 
                'data': {
                    'image': 'abcdefghijklmn'
                }, 
                'exchange_id': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'id': '23', 
                'refer': '', 
                'type': 'DEAL', 
                'userkey_a': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
                'userkey_b': 'contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx'
            }
        ]
    }, 
    'msg': 'success'
}
```



##### 3.3 查询指定区块ID的交易内容

请求URL

> http://<host>:<port>/api/query_block

请求方式

> POST

输入参数（data字段下）

| 参数     | 类型   | 说明     |
| -------- | ------ | -------- |
| userkey  | string | 用户公钥 |
| block_id | string | 区块ID   |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 指定区块的交易/授权数据                 |

> 说明：
>
> 按区块ID查询时没有限制链用户范围。

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex", 
        "block_id": "21"
    }, 
    "timestamp": 1618284796, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "YWY2YmY0ZmYzMjJmOTE5YzVjOTg0YWQ0Zjk3NjQ1MzE5YzMyOWMzOThlN2ZjZWU1YmM2MjQwOWFjNTUxNDg5NQ=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'deals': {
            'action': '12', 
            'assets_id': '1234', 
            'data': {'image': '11111111111111111111'}, 
            'exchange_id': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
            'id': '21', 
            'refer': '', 
            'type': 'DEAL', 
            'userkey_a': 'contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex', 
            'userkey_b': 'contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx'
        }
    }, 
    'msg': 'success'
}
```



##### 3.4 查询指定区块ID的原始区块数据

请求URL

> http://<host>:<port>/api/query_raw_block

请求方式

> POST

输入参数（data字段下）

| 参数    | 类型   | 说明     |
| ------- | ------ | -------- |
| userkey | string | 用户公钥 |
| height  | string | 区块高度 |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 指定区块的原始区块数据                  |

> 说明：
>
> 按区块ID查询时没有限制链用户范围。

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey": "contract1ghfcl0hm5pxu0q0jgnl2nw3hhmrkklgyh3lgvx", 
        "height": '210274'
    }, 
    "timestamp": 1618284344, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "OWMxZDZlMGYxNDY2Y2Q1YWQyN2JlZGQzYzcxY2Y0ZGNlYmNmOTBmODRjNjM5MzA4ZmYyZDg0MWY2Y2FlZTFjYQ=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'blcok': {
            'block': {
                'data': {
                    'txs': ['Co0CCooCCiwvamFjazEzOS5jb250cmFjdC5jb250cmFjdC5Nc2dDcmVhdGVDb250cmFjdBLZAQovY29udHJhY3QxOHpmZHNqem44dDUwMHF1NDNjcHhocjlzenMyMjZyODh1OGhseTkSBDEyMzQaL2NvbnRyYWN0MWxhbnJ2enhkOTl4eTAwempneGZqbTVwZHFoczVqdjZoNXo5bWV4Ii9jb250cmFjdDFnaGZjbDBobTVweHUwcTBqZ25sMm53M2hobXJra2xneWgzbGd2eCoCMTIyOnsiaW1hZ2UiOiJRbVFXdVg3bXdFNUxEdGpja2M3M3ZpV1o3TEJpVTRnRW5YcXRKdkMyN1JSUnlZIn0SWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA/6+qwsJW0bHd0OaCqa2Mfxr1lRQGE9NtS/+66lG9EO7EgQKAggBGAMSBBDAmgwaQAbZAl1RUh8EdAAWDWqx+MKKPrZ9JRW0PxgdTAdVWKT1OE8xw3Cq3wuUiuTCAajsEmCKjKeqcvz3UJC5eL1O93U=']
                }, 
                'evidence': {'evidence': None}, 
                'header': {
                    'app_hash': '767B1E1A1F42FB08187284E3831E4065AE42F2EA04C55B7EE84C8856E95BECD0', 
                    'chain_id': 'contract', 
                    'consensus_hash': '048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F', 
                    'data_hash': '71FC654A49D737F32092B64E32DE7569A072EAAEEB0DC0E4C2C16331D0672414', 
                    'evidence_hash': 'E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855', 
                    'height': '210274', 
                    'last_block_id': {
                        'hash': '4BB4962E7A1462DC665A888652652D652C21E24DF994E9260578093C8F7794A5', 
                        'parts': {
                            'hash': '35C561892C4B3F789D0AA714F8B62EBA1AE97B1A2106523C20F91111021CB3D0', 
                            'total': 1
                        }
                    }, 
                    'last_commit_hash': '3654C6C1A7BB84C44FFDCE14ACD2F3983FA4BE6678E13B99ED890845C1709E66', 
                    'last_results_hash': 'E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855', 
                    'next_validators_hash': '98B670043EDB03D07C4096C7DE5BC389EA11DB382A7E501F4635F7B73482C078', 
                    'proposer_address': '627C9E0096F61C1A40A980781B38B3CFC7B32E93', 
                    'time': '2021-04-09T09:03:37.86117745Z', 
                    'validators_hash': '98B670043EDB03D07C4096C7DE5BC389EA11DB382A7E501F4635F7B73482C078', 
                    'version': {'block': '11'}
                }, 
                'last_commit': {
                    'block_id': {
                        'hash': '4BB4962E7A1462DC665A888652652D652C21E24DF994E9260578093C8F7794A5', 
                        'parts': {
                            'hash': '35C561892C4B3F789D0AA714F8B62EBA1AE97B1A2106523C20F91111021CB3D0', 
                            'total': 1
                        }
                    }, 
                    'height': '210273', 
                    'round': 0, 
                    'signatures': [
                        {
                            'block_id_flag': 2, 
                            'signature': 'SwbBqVdczrfxTrvwCCUykvl2ZBM48MYyaGsYFirj4MksyWLI9hZSOMoHjKQqB7BxeDXaPbjqpyUBSsxHTWfbBw==', 
                            'timestamp': '2021-04-09T09:03:37.86117745Z', 
                            'validator_address': '627C9E0096F61C1A40A980781B38B3CFC7B32E93'
                        }
                    ]
                }
            }, 
            'block_id': {
                'hash': '3691B91FC6B22CB271AC0B20135200112716605F57F8C85C609FE4C2908011B1', 
                'parts': {
                    'hash': 'E2B3B49B97634892FD230C05F2D6116055070BF1043DE57C8C32C19647E19194', 
                    'total': 1
                }
            }
        }
    }, 
    'msg': 'success'
}
```



##### 3.5 查询通证余额

请求URL

> http://<host>:<port>/api/query_balance

请求方式

> POST

输入参数（data字段下）

| 参数    | 类型   | 说明     |
| ------- | ------ | -------- |
| userkey | string | 用户公钥 |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 通证余额信息                            |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "userkey": "contract1lanrvzxd99xy00zjgxfjm5pdqhs5jv6h5z9mex"
    }, 
    "timestamp": 1618295472, 
    "appid": "4fcf3871f4a023712bec9ed44ee4b709", 
    "sign_data": "MzI1YzE5ZWFkM2NmNTMzNjFiMWVmYTMwM2ZhZmU2MDQwMWU0NzJkM2QzMDA1OWM1YWI0ZjY5NjUwODQwMzg0ZA=="
}
```

返回示例

```json
{
    'code': 0, 
    'data': {
        'blcok': {
            'amount': '20',   /* 用户通证数量 */
            'denom': 'credit' /* 通证单位 */
        }
    }, 
    'msg': 'success'
}
```


