## 查询智能合约
- url:/api/contracts
- method:get

input:param
```param
&limit=40       //每页40条记录
&start=0        //记录的起始序号
&sort=-timestamp       //按照时间戳倒序排列

http://18.216.57.65:20110/api/contracts?sort=-name&limit=40&start=0
```
output:json
```json
{
    "total":2169998,
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":[
        {
            "address":"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi", //合约地址
            "name":"xxxxxxx",//合约名称
            "compiler":"v0.0.1",//编译器版本
            "balance":100,//余额
            "trxCount":6000, //交易数量
            "isSetting":true,//是否优化
            "dateVerified":1531711638107,//合约验证时间
        },
        {
            "address":"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi", //合约地址
            "name":"xxxxxxx",//合约名称
            "compiler":"v0.0.1",//编译器版本
            "balance":100,//余额
            "trxCount":6000, //交易数量
            "isSetting":true,//是否优化
            "dateVerified":1531711638107,//合约验证时间
        },...
     ]
}
```
返回所有已经验证的智能合约


## 合约详情信息
- url:/api/contract/:address
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/contract/TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX
```
output:json
```json
{
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":{
        "address":"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi", //合约地址
        "balance":100,//余额
        "trxCount":6000, //交易数量
        "tokenContract":"Local Token (LOT)",//
        "creator":{
            "address":"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi",//创建者地址
            "txHash":"******************",//合约创建hash
        }
    }  
}
```
返回合约详情页的合约概况

## 查询智能合约交易列表
- url:/api/contracts/transaction
- method:get

input:param
```param
&limit=40       //每页40条记录
&start=0        //记录的起始序号
&count=true     //是否返回总数
&sort=-timestamp       //按照时间戳倒序排列
&contract=************     //合约地址
&type=internal/token //交易内合约   主网暂不支持

http://18.216.57.65:20110/api/contracts/transaction?sort=-timestamp&count=true&limit=40&start=0&contract=TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX
```
output:json
```json
{
    "total":2169998,
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":[
        {
            "txHash":"******************", //交易hash
            "parentHash":"******************", //交易父hash only for internal交易
            "block":66000,//区块高度
            "timestamp":1531711638107,//交易时间
            "ownAddress":"*************",//发起人地址  owner_address???
            "toAddress":"************", //接收人地址   contract_address???
            "value":10,//
            "txFee":1,//交易费 ？？？ transaction tab  如何判断out of gas 和转入转出交易？？？？
            "token":"ATRON",//token only for token 交易
        },
        {
            "txHash":"******************", //交易hash
            "parentHash":"******************", //交易父hash only for internal交易
            "block":66000,//区块高度
            "timestamp":1531711638107,//交易时间
            "ownAddress":"*************",//发起人地址
            "toAddress":"************", //接收人地址
            "value":10,//
            "txFee":1,//交易费
            "token":"ATRON",//token only for token 交易
        },...
     ]
}
```
合约详情页-【Transactions】页       
当type=internal：合约详情页-【Internal Txns】页     
当type=token：合约详情页-【Token Txns】页         
返回合约详情页transaction信息


## 查询智能合约code信息
- url:/api/contracts/code
- method:get

input:param
```param
&contract=************     //合约地址

http://18.216.57.65:20110/api/contracts/code?contract=TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX
```
output:json
```json
{
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":{
        "address":"***********",//合约地址
        "name":"*****",//合约名称
        "compiler":"v4.0.0",//编译器版本
        "isSetting":true,//是否优化
        "source":"***********",//合约源代码
        "byteCode":"*****",//编译生成的二进制代码
        "abi":"*****",//编译生成的abi
        "abiEncoded":"********",//编译所需参数
        "librarys":[{
            "index":1,//库序号，最大五个library
            "name":"name",//名称
            "address":"address",//地址
        }]
    }
}
```
合约详情页-【Code】页     
返回合约详情页的code信息


## 查询智能合约event信息
- url:/api/contracts/event
- method:get

input:param
```param
&contract=************     //合约地址

http://18.216.57.65:20110/api/contracts/event?contract=TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX
```
output:json
```json
{
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":[{
        "txHash":"******************", //交易hash
        "block":66000,//区块高度
        "timestamp":1531711638107,//交易时间
        "method":"transfer(address,uint256)",//方法名
        "eventLog":"************"//event log 日志
    }...
    ]
}
```
合约详情页-【Events】页     
返回合约详情页的event信息

## 查询智能合约内部交易
- url:/api/contracts/internalTxs
- method:get

input:param
```param
&limit=40       //每页40条记录
&start=0        //记录的起始序号
&sort=-timestamp       //按照时间戳倒序排列
&contract=************     //合约地址

http://18.216.57.65:20110/api/contracts/internalTxs?sort=-name&limit=40&start=0
```
output:json
```json
{
    "total":2169998,
    "status":{
        "code":"0",
        "message":"success"
    },
    "data":[
        {
            "block":21232, //区块id
            "timestamp":1531711638107,//交易时间
            "parentHash":"******************",//父hash
            "txType":"call",//交易类型
            "ownerAddress":"**************", //发起人地址
            "toAddress":"**************",//接收人地址
            "value":100,//交易额
            "txFee":1,//交易费
        },
        {
            "block":21232, //区块id
            "timestamp":1531711638107,//交易时间
            "parentHash":"******************",//父hash
            "txType":"call",//交易类型
            "ownerAddress":"**************", //发起人地址
            "toAddress":"**************",//接收人地址
            "value":100,//交易额
            "txFee":1,//交易费
        },...
     ]
}
```
返回区块内全部合约内部交易

## 验证智能合约
- url:/api/contracts/verify
- method:POST

input:json
```json
{
    "address":"***********",//合约地址
    "name":"*****",//合约名称
    "compiler":"v4.0.0",//编译器版本
    "isSetting":true,//是否优化
    "source":"***********",//合约源代码
    "byteCode":"*****",//编译生成的二进制代码
    "abi":"*****",//编译生成的abi
    "abiEncoded":"********",//编译所需参数
    "librarys":[{
        "index":1,//库序号，最大五个library
        "name":"name",//名称
        "address":"address",//地址
    }]
}

```
output:json
```json
{
    "status":{
        "code":"0",
        "message":"success"
    }
}
```
根据前端传入地址，source code以及编译后生成的abi，byte code，去主网校验，如果abi和byte code相同，则校验成功存入数据库，否则校验失败         
返回code：0 验证成功，写数据库成功<br>
         1 验证失败，跟主网合约信息不匹配<br>
         2 验证成功，写数据库或内部逻辑错误