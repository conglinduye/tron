## 查询账户信息
- url:/api/account
- method:get

input:param
```param
&limit=40       //每页40条记录
&count=true     //是否返回总数
&start=0        //记录的起始序号
&sort=-balance       //按照账户余额倒序排列

eg: http://18.216.57.65:20110/api/account?sort=-balance&start=0&limit=10

```
output:json
```json
{
    "total":323497,
    "data":[
        {
            "address":"TDtjQ1JR5UrS92W9kB6BCeAQJwn1dyBEbs",//账户地址
            "name":"Parkseungwan",//账户名称
            "balance":0,//账户余额
            "power":54418725400,//投票权
            "tokenBalances":{//各种token
                "DEX":0,
                "GOD":1000,
                "IGG":1,
                "WIN":4504,
                "IPFS":18,
                "SEED":4,
                "DICKS":2,
                "Durex":2,
                "REcoin":2140,
                "Bitcoin":133,
                "Dislike":50,
                "Pornhub":2,
                "binance":13,
                "ofoBike":14,
                "OnePiece":11,
                "Perogies":1,
                "RingCoin":11,
                "TRONGOLD":1000000,
                "CyberTron":30000,
                "Messenger":10,
                "Skypeople":24,
                "BitTorrent":16,
                "FomoThreeD":4,
                "bittorrent":1,
                "MusicCasper":1,
                "TRXTestCoin":25,
                "MarcusMillichap":10,
                "CommunityNodeToken":100
            },
            "dateCreated":1536151879581,//创建时间
            "dateUpdated":1536151879581//更新时间
        },
        {
            "address":"TC1Y47GVt3bEB4UkAJwXxt284vPNCaTmnf",
            "name":"",
            "balance":561526,
            "power":0,
            "tokenBalances":{
                "IPFS":6,
                "Bitcoin":3,
                "NBACoin":2,
                "binance":3,
                "ofoBike":4,
                "Messenger":2,
                "Skypeople":5,
                "BitTorrent":3,
                "HuobiToken":1,
                "TRXTestCoin":1
            },
            "dateCreated":1536549025630,
            "dateUpdated":1536549025630
        },...
     ]
}
```


## 单个账户信息  
- url:/api/account/:address
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/account/T9ya5cLUd4LUXit2BR5fuG7VCA87RrnTk5
```
output:json
```json
{
    "representative":{
        "enabled":false,
        "lastWithDrawTime":0,
        "allowance":0,
        "url":null
    },
    "name":"",
    "address":"TC1Y47GVt3bEB4UkAJwXxt284vPNCaTmnf",
    "bandwidth":{
        "freeNetUsed":123,
        "freeNetLimit":5000,
        "freeNetRemaining":4877,
        "freeNetPercentage":2.46,
        "netUsed":0,
        "netLimit":0,
        "netRemaining":0,
        "netPercentage":0,
        "assets":{
            "NBACoin":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "TRXTestCoin":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "BitTorrent":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "HuobiToken":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "ofoBike":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "Bitcoin":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "binance":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "Messenger":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "Skypeople":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            },
            "IPFS":{
                "netUsed":0,
                "netLimit":0,
                "netRemaining":0,
                "netPercentage":0
            }
        }
    },
    "balances":[
        {
            "name":"TRX",
            "balance":0.561526
        },
        {
            "name":"NBACoin",
            "balance":2
        },
        {
            "name":"TRXTestCoin",
            "balance":1
        },
        {
            "name":"BitTorrent",
            "balance":3
        },
        {
            "name":"HuobiToken",
            "balance":1
        },
        {
            "name":"ofoBike",
            "balance":4
        },
        {
            "name":"Bitcoin",
            "balance":3
        },
        {
            "name":"binance",
            "balance":3
        },
        {
            "name":"Messenger",
            "balance":2
        },
        {
            "name":"Skypeople",
            "balance":5
        },
        {
            "name":"IPFS",
            "balance":7
        }
    ],
    "balance":561526,
    "tokenBalances":[
        {
            "name":"TRX",
            "balance":0.561526
        },
        {
            "name":"NBACoin",
            "balance":2
        },
        {
            "name":"TRXTestCoin",
            "balance":1
        },
        {
            "name":"BitTorrent",
            "balance":3
        },
        {
            "name":"HuobiToken",
            "balance":1
        },
        {
            "name":"ofoBike",
            "balance":4
        },
        {
            "name":"Bitcoin",
            "balance":3
        },
        {
            "name":"binance",
            "balance":3
        },
        {
            "name":"Messenger",
            "balance":2
        },
        {
            "name":"Skypeople",
            "balance":5
        },
        {
            "name":"IPFS",
            "balance":7
        }
    ],
    "frozen":{
        "total":0,
        "balances":[

        ]
    }
}
```

## 某账户的媒体信息
- url:/api/account/:address/media
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/account/T9ya5cLUd4LUXit2BR5fuG7VCA87RrnTk5/media
```

output:json
```json
{
    "reason":"Could not retrieve file",
    "success":false
}
```

## 验签
- url:/api/auth
- method:POST

input:json
```json
{
    "transaction":"123213242"
}
```
output:json
```json  
{
    "key":"123213242"
}
```
返回的key用途：
调用【修改超级代表github信息】接口时，将key设置请求头【X-Keys】中，用于修改前的校验


## 修改超级代表github信息
- url:/api/account/:address/sr
- method:POST

input:json
```json
{
    "address": "RTFGHJK6GCVHHByui765CVBCVBVB",
    "githubLink": "sesameseed/tronsr-template" 
}
```

output:json
```json

```
修改前需要校验请求头中【X-Key】的key，并校验key解析出来的address与请求参数的address是否一致，如果一致，则继续执行修改逻辑

## 查询超级代表github信息
- url:/api/account/:address/sr
- method:GET

input:param
```param
eg: http://18.216.57.65:20110/api/account/TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp/sr
```
output:json
```json
{
    "address": "RTFGHJK6GCVHHByui765CVBCVBVB",
    "githubLink": "sesameseed/tronsr-template" 
}
```


## 查询用户的交易统计信息
- url:/api/account/:address/stats
- method:GET

input:param
```param
eg: http://18.216.57.65:20110/api/account/TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp/stats
```
output:json
```json
{
    "transactions": "827",
    "transactions_out": "230",
    "transaction_in": "597"
}
```
交易数包括TRX 交易和asset 交易
