## 登录获取token信息
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
curl -XGET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/list

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
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/add -d'{"ownerAddress":"test", "AssetName":"test"}'
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

curl -XDELETE -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgyMjMwMzUsImlhdCI6MTUzODIyMTIzNSwiaWQiOjEsIm5iZiI6MTUzODIyMTIzNSwidXNlcm5hbWUiOiJ0cm9uIn0.eemY-FhIM1wCMiRMn7XkzzsV7WKkgekVvYr0U424cBg" -H "Content-Type: application/json" -H "Content-Type: application/json" https://wlcyapi.tronscan.org/api/tokenBlacklist/delete/13
```
ouput:json
```json
{"code":0,"message":"OK","data":null}
```
