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
http://18.216.57.65:20110/api/token?start=0&limit=20

http://18.216.57.65:20110/api/token?start=0&limit=20&status=ico

http://18.216.57.65:20110/api/token?owner=TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh&name=VIP
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

## 根据通证名称查询通证转账
- url:/api/asset/transfer
- method:get

input:param
```param

eg:
http://18.216.57.65:20110/api/asset/transfer?start=0&limit=20&name=VIP
```
ouput:json
```json

{
    "total":2,
    "Data":[
        {
            "blockId":2385200, //区块ID
            "transactionHash":"00041899b98d38de173431434d23478f0a59c1b52b6cb024e64c08205a3350ec",//交易hash
            "timestamp":1537063485000, //创建时间戳
            "transferFromAddress":"TKjdnbJxP4yHeLTHZ86DGnFFY6QhTjuBv2", //转账from地址
            "transferToAddress":"TFGw3DMJS9NPZnV14wkw6K6aiDP6HgvXA3", //转账to地址
            "amount":1, //转账数量
            "tokenName":"HuobiToken", //通证名称
            "confirmed":true //是否确认
        },
        {
            "blockId":2346000,
            "transactionHash":"000e9ac7cc63ad472aa6bb753878ce6cd863f23a67aafb897bac426c3fbd3987",
            "timestamp":1536945795000,
            "transferFromAddress":"TFMAHNTS35ZNL5Pr9XToVVekyAbE9ts65Y",
            "transferToAddress":"TQoXdQAhvgiXCX6HrfKNFKSUBTmuVTDpoi",
            "amount":1,
            "tokenName":"BitTorrent",
            "confirmed":true
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

## 登录获取访问api的token信息
- url:/api/login
input:param
```param
eg:
curl -XPOST -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/login -d'{"username":"***","password":"***"}'
```
output:json
```json
{
    "code":0,
    "message":"OK",
    "data":{
        "token":"***"
    }
}
```

## 查询通证黑名单列表信息
- url:/api/tokenBlacklist/list
- method:get

input:param
```param
start=0                                             //起始序号
&limit=20                                           //分页数目
&ownerAddress=TXiBuTWoXvWYKz47NcvKe1gfQDdNHnbBmh    //创建者地址
&name=tokenName                                     //通证名称

eg: 
curl -XGET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/list

```
output:json
```json
{
    "code":0,
    "message":"OK",
    "data":{
        "total":12,
        "data":[
            {
                "id":12,
                "ownerAddress":"TKiMWYxYh3EgtTdVd2fQ2dbLXEtAgZ5aiq",
                "assetName":"TWM",
                "createTime":"2018-09-26 02:57:43"
            },
            {
                "id":11,
                "ownerAddress":"TN93Yq91SpCtiHaLGNXq8HCJ4WsvRK5dup",
                "assetName":"binance",
                "createTime":"2018-09-25 18:16:57"
            },
            {
                "id":10,
                "ownerAddress":"TWhRyN6bWEWEgogav9ECApqqAiGiDZzsg1",
                "assetName":"ZZZ",
                "createTime":"2018-09-25 18:16:40"
            },
            {
                "id":9,
                "ownerAddress":"TEQKZ46TWgPw4ZDVHFP8xzoWuuD2sk8mwj",
                "assetName":"ZTX",
                "createTime":"2018-09-25 18:16:28"
            },
            {
                "id":8,
                "ownerAddress":"TWn6LWzMUY4mYW8ohXBiPwPCe6LMQjMtxf",
                "assetName":"XP",
                "createTime":"2018-09-25 18:16:17"
            },
            {
                "id":7,
                "ownerAddress":"TH1Rk4jnH3e8DZhuw6dnYWpVBR5wW4gCXU",
                "assetName":"WWGoneWGA",
                "createTime":"2018-09-25 18:16:06"
            },
            {
                "id":6,
                "ownerAddress":"TD2wyhKFydHzG7evRkkKY391WJqKsbvX8i",
                "assetName":"VBucks",
                "createTime":"2018-09-25 18:15:52"
            },
            {
                "id":5,
                "ownerAddress":"TUBDh44BponJideGUrdM2j1GWjPF4DueVb",
                "assetName":"TronWatchmarket",
                "createTime":"2018-09-25 18:15:42"
            },
            {
                "id":4,
                "ownerAddress":"TEhYJhHs795aP6LBkw2KMyNFAQxJVgWdop",
                "assetName":"TronWatch",
                "createTime":"2018-09-25 18:15:35"
            },
            {
                "id":3,
                "ownerAddress":"TWxezzaF6Lyvh4mNVQRj9okro8hqC3LfJt",
                "assetName":"Skypeople",
                "createTime":"2018-09-25 18:15:19"
            },
            {
                "id":2,
                "ownerAddress":"TAahb2b3TPUL8rMUk6LQLiFFkFYrUX1qDr",
                "assetName":"Fortnite",
                "createTime":"2018-09-25 18:15:08"
            },
            {
                "id":1,
                "ownerAddress":"TVpnuztB8UYaxWAC6MLx1rZn9KXm44BLVb",
                "assetName":"CheapAirGoCoin",
                "createTime":"2018-09-25 18:15:01"
            }
        ]
    }
}


```

## 添加通证黑名单
- url:/api/tokenBlacklist/add
- method:post

input:param
```param

eg:
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/add -d'{"ownerAddress":"test", "AssetName":"test"}'
```
output:json
```json
{"code":0,"message":"OK","data":null}
```

## 删除通证黑名单
- url:/api/tokenBlacklist/delete/:id
- method:delete

input:param
```param

eg:

curl -XDELETE -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/delete/13
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
```

## 添加tokenExt信息
- url:/api/tokenExt/addInfo
- method:post

input:param
```param

eg:
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyOTYyMzQsImlhdCI6MTUzODI5NDQzNCwiaWQiOjEsIm5iZiI6MTUzODI5NDQzNCwidXNlcm5hbWUiOiJ0cm9uIn0.nYBlNxYaL8cxGZxTTYEr5uwiODwJj3_8ilFdbjH2V20" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenExt/addInfo -d'{"address":"testAddress", "tokenName":"testTokenName", "tokenId":"111", "brief":"test2", "website":"test3", "whitePaper":"test4", "github":"test5", "country":"test6", "credit":"1", "reddit":"test7", "twitter":"test8", "facebook":"test9", "telegram":"test10", "steam":"test11", "medium":"test12", "webchat":"test13", "weibo":"test14", "review":"1", "status":"1"}'
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
```

## 修改tokenExt信息
- url:/api/tokenExt/updateInfo
- method:post

input:param
```param

eg:
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyOTYyMzQsImlhdCI6MTUzODI5NDQzNCwiaWQiOjEsIm5iZiI6MTUzODI5NDQzNCwidXNlcm5hbWUiOiJ0cm9uIn0.nYBlNxYaL8cxGZxTTYEr5uwiODwJj3_8ilFdbjH2V20" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenExt/updateInfo -d'{"address":"testAddress", "tokenName":"000testTokenName", "tokenId":"000111", "brief":"000test2", "website":"000test3", "whitePaper":"000test4", "github":"000test5", "country":"000test6", "credit":"0001", "reddit":"000test7", "twitter":"000test8", "facebook":"000test9", "telegram":"000test10", "steam":"000test11", "medium":"000test12", "webchat":"000test13", "weibo":"000test14", "review":"0", "status":"0"}'
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
```

## 添加tokenLogo信息
- url:/api/tokenExt/addLogo
- method:post

input:param
```param

eg:
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyOTYyMzQsImlhdCI6MTUzODI5NDQzNCwiaWQiOjEsIm5iZiI6MTUzODI5NDQzNCwidXNlcm5hbWUiOiJ0cm9uIn0.nYBlNxYaL8cxGZxTTYEr5uwiODwJj3_8ilFdbjH2V20" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenExt/addLogo -d'{"address":"testAddress", "logoUrl":"testurl"}'
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
```

## 修改tokenLogo信息
- url:/api/tokenExt/updateLogo
- method:post

input:param
```param

eg:
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyOTYyMzQsImlhdCI6MTUzODI5NDQzNCwiaWQiOjEsIm5iZiI6MTUzODI5NDQzNCwidXNlcm5hbWUiOiJ0cm9uIn0.nYBlNxYaL8cxGZxTTYEr5uwiODwJj3_8ilFdbjH2V20" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenExt/updateLogo -d'{"address":"testAddress", "logoUrl":"000testurl"}'
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
