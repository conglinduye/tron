## 查询交易信息
- url:/api/transaction
- method:get

input:param
```param
&limit=40       //每页40条记录
&count=true     //是否返回总数
&start=0        //记录的起始序号
&sort=-timestamp       //按照时间戳倒序排列
&total=123421   //分页查询时，传入上一次接口返回的数据总量，用作分页查询数据校准

eg: http://18.216.57.65:20110/api/transaction?sort=-number&limit=40&start=0
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


## 单个交易信息
- url:/api/transaction/:hash
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/transaction/0000fa99d223907c660722465ab78306302f99c3423d3abe9b3f9fb086338df3
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

## 创建交易
- url:/api/transaction
- method:POST

input:param
```json
{
    "transaction":"0A84010A025006220880DDBCA411E6159840E8F7E1BDDC2C5204484148415A67080112630A2D747970652E676F6F676C65617069732E636F6D2F70726F746F636F6C2E5472616E73666572436F6E747261637412320A1541E552F6487585C2B58BC2C9BB4492BC1F17132CD012154190919CBA90CE96F9B9A63AFDE5AC66453D3F690E18C0843D124183A239CD8B1A3998B56DF45667E230B6BED13C10889BED456F40716DC5558F58715DD1E31DDF57419C6FD5F4864C5BD8995E20306C44609FA3C4ED556DA6DCE600"
}
```
output:json
```json
{
    "success":true,
    "code":"SUCCESS",
    "message":"",
    "transaction":{
        "hash":"0afa11cbfa9b4707b1308addc48ea31201157a989db92fe75750c068f0cc14e0",
        "timestamp":0,
        "contracts":[
            {
                "contractType":"TransferContract",
                "contractTypeId":1,
                "owner_address":"TWsm8HtU2A5eEzoT8ev8yaoFjHsXLLrckb",
                "to_address":"TP9cjznZ4rGL3Hav7d8zPg6HX6cXdWecMc",
                "amount":1000000
            }
        ],
        "data":"HAHA",
        "signatures":[
            {
                //"bytes":"G4OiOc2LGjmYtW30VmfiMLa+0TwQiJvtRW9AcW3FVY9YcV3R4x3fV0Gcb9X0hkxb2JleIDBsRGCfo8TtVW2m3OY=",
                "bytes":"g6I5zYsaOZi1bfRWZ+Iwtr7RPBCIm+1Fb0BxbcVVj1hxXdHjHd9XQZxv1fSGTFvYmV4gMGxEYJ+jxO1Vbabc5gA=",
                "address":"TWsm8HtU2A5eEzoT8ev8yaoFjHsXLLrckb"
            }
        ]
    }
}
```