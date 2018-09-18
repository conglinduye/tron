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

eg: 
http://18.216.57.65:20110/api/token?start=0&limit=10

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

eg:
http://18.216.57.65:20110/api/token/VIP

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
    "dateCreated":1534304985000,
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
## 根据通证名称查询通证持有人信息
- url:/api/token/:name/address
- method:get

input:param
```param

eg:
http://18.216.57.65:20110/api/token/VIP/address
```
output:json
```json
{
    "total":2,
    "data":[
        {
            "address":"TWWBvgxBteNn1KcQXMrx4yCZTcLfZd49sB",
            "name":"VIP",
            "balance":1
        },
        {
            "address":"TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh",
            "name":"VIP",
            "balance":2069328
        }
    ]
}
```

## 上传Token图片
- url:/api/uploadLogo
- method:post
input:param
```param
{
  "imageData":"data:image/png;base64,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",
  "owner_address":"TVVGvh3DrRrUCuZVy58Ha4QRqF4gGMS7L3"
}

eg:
curl -XPOST -H "Content-Type:application/json" http://18.216.57.65:20110/api/uploadLogo
-d'{"imageData":"data:image/png;base64,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","owner_address":"TVVGvh3DrRrUCuZVy58Ha4QRqF4gGMS7L3"}'

```
ouput:json
```json
{
    "success":true,
    "data":"http://coin.top/tokenLogo/tokenLogo_20180915113703.615207.jpeg"
}
```

## 获取下载TokenTemplateFile地址
- url:/api/download/tokenInfo
- method:get

input:param
```param

eg:
http://18.216.57.65:20110/api/download/tokenInfo

```
output:json
```json
{
    "success":true,
    "data":"http://coin.top/tokenTemplate/TronscanTokenInformationSubmissionTemplate.xlsx"
}
```