## 查询通证列表信息
- url:/api/token
- method:get

input:param
```param
start=0                                     //起始序号
&limit=40                                   //分页数目
&owner=TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh   //创建者地址
&name=VIP                                   //通证名称
&status=ico                                 //ico是否结束标识
```
output:json
```json
{
    "total":10,
    "data":[
        {
            "index":1,
            "price":100000000,
            "issued":1,
            "issuedPercentage":0.00004347826086956522,
            "available":2069328,
            "availableSupply":2069329,
            "remaining":2069328,
            "remainingPercentage":89.97078260869566,
            "percentage":89.97078260869566,
            "frozenTotal":230671,
            "frozenPercentage":10.029173913043477,
            "ownerAddress":"TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh",
            "name":"VIP",
            "totalSupply":2300000,
            "trxNum":100000000,
            "num":1,
            "endTime":1533709512000,
            "startTime":1533709512000,
            "voteScore":0,
            "description":"Token to display exclusive status and gain special privileges.",
            "url":"",
            "frozen":[
                {
                    "days":1533,
                    "amount":57
                },
                {
                    "days":328,
                    "amount":230000
                },
                {
                    "days":520,
                    "amount":557
                },
                {
                    "days":2043,
                    "amount":57
                }
            ],
            "abbr":"VIP",
            "participated":5600,
            "totalTransactions":0,
            "nrOfTokenHolders":0,
            "tokenID":"",
            "reputation":"insufficient_message",
            "imgUrl":"",
            "website":"no_message",
            "white_paper":"no_message",
            "github":"no_message",
            "country":"no_message",
            "social_media":[
                {
                    "name":"Reddit",
                    "url":""
                },
                {
                    "name":"Twitter",
                    "url":""
                },
                {
                    "name":"Facebook",
                    "url":""
                },
                {
                    "name":"Telegram",
                    "url":""
                },
                {
                    "name":"Steem",
                    "url":""
                },
                {
                    "name":"Medium",
                    "url":""
                },
                {
                    "name":"Wechat",
                    "url":""
                },
                {
                    "name":"Weibo",
                    "url":""
                }
            ]
        },...
    ]
}
```

## 根据通证名称查询通证信息
- url:/api/token/:name
- method:get

input:param
```param

```
output:json
```json
{
    "index":0,
    "price":100000000,
    "issued":1,
    "issuedPercentage":0.00004347826086956522,
    "available":2069328,
    "availableSupply":2069329,
    "remaining":2069328,
    "remainingPercentage":89.97078260869566,
    "percentage":89.97078260869566,
    "frozenTotal":230671,
    "frozenPercentage":10.029173913043477,
    "ownerAddress":"TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh",
    "name":"VIP",
    "totalSupply":2300000,
    "trxNum":100000000,
    "num":1,
    "endTime":1533709512000,
    "startTime":1533709512000,
    "voteScore":0,
    "description":"Token to display exclusive status and gain special privileges.",
    "url":"",
    "frozen":[
        {
            "days":1533,
            "amount":57
        },
        {
            "days":328,
            "amount":230000
        },
        {
            "days":520,
            "amount":557
        },
        {
            "days":2043,
            "amount":57
        }
    ],
    "abbr":"VIP",
    "participated":5600,
    "totalTransactions":0,
    "nrOfTokenHolders":2,
    "tokenID":"",
    "reputation":"insufficient_message",
    "imgUrl":"",
    "website":"no_message",
    "white_paper":"no_message",
    "github":"no_message",
    "country":"no_message",
    "social_media":[
        {
            "name":"Reddit",
            "url":""
        },
        {
            "name":"Twitter",
            "url":""
        },
        {
            "name":"Facebook",
            "url":""
        },
        {
            "name":"Telegram",
            "url":""
        },
        {
            "name":"Steem",
            "url":""
        },
        {
            "name":"Medium",
            "url":""
        },
        {
            "name":"Wechat",
            "url":""
        },
        {
            "name":"Weibo",
            "url":""
        }
    ]
}
```

## 获取下载TokenTemplateFile地址
- url:/api/download/tokenInfo
- method:get

input:param
```param

```
output:json
```json
   "http://coin.top/tokenTemplate/TronscanTokenInformationSubmissionTemplate.xlsx"

```