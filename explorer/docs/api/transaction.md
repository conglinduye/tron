## 查询交易信息
- url:/api/transaction
- method:get

input:param
```param
&limit=40       //每页40条记录
&count=true     //是否返回总数
&start=0        //记录的起始序号
&sort=-timestamp       //按照时间戳倒序排列
```
output:json
```json
{
    "total":2169998,
    "data":[
        {
            "hash":"109917ca3ccd1452557604d2616f387ce78341708d903de58e33f52807f2eba8",//交易hash
            "block":2214873,//所属区块
            "timestamp":1536551478000,//时间戳
            "confirmed":false,//是否确认
            "ownerAddress":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",//交易发起人
            "toAddress":"TTs8B82fuxtDKpYx2qhC8THm32Lm9Ng4jv",//交易接受人
            "contractData":{//
                "to":"TTs8B82fuxtDKpYx2qhC8THm32Lm9Ng4jv",//交易接受人
                "from":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",//交易发起人
                "token":"IPFS",//token名称
                "amount":1//数量
            },
            "contractType":2,//交易类型
            "data":""
        },
        {
            "hash":"5503ecc358572d464ae2553ef5c49c218f0d242dff34673366528906f104b456",
            "block":2214873,
            "timestamp":1536551478000,
            "confirmed":false,
            "ownerAddress":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",
            "toAddress":"TVq5Ayig6EgmsYXtBsSXNSbXfxCNcZFuGf",
            "contractData":{
                "to":"TVq5Ayig6EgmsYXtBsSXNSbXfxCNcZFuGf",
                "from":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",
                "token":"IPFS",
                "amount":1
            },
            "contractType":2,
            "data":""
        },...
     ]
}
```


## 单个区块信息
- url:/api/transaction/:hash
- method:get

input:param
```param

```
output:json
```json
{
    "hash":"109917ca3ccd1452557604d2616f387ce78341708d903de58e33f52807f2eba8",//交易hash
    "block":2214873,//所属区块
    "timestamp":1536551478000,//时间戳
    "confirmed":false,//是否确认
    "ownerAddress":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",//交易发起人
    "toAddress":"TTs8B82fuxtDKpYx2qhC8THm32Lm9Ng4jv",//交易接受人
    "contractData":{//
        "to":"TTs8B82fuxtDKpYx2qhC8THm32Lm9Ng4jv",//交易接受人
        "from":"TV3NmH1enpu4X5Hur8Z16eCyNymTqKXQDP",//交易发起人
        "token":"IPFS",//token名称
        "amount":1//数量
    },
    "contractType":2,//交易类型
    "data":""
}
```
