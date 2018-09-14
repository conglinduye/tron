## 获得数据同步信息
- url:/api/system/status
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/system/status
```
output:json
```json
      
{
    "network":{
        "type":"mainnet"
    },
    "sync":{
        "progress":99.99988929123896//totalProgress
    },
    "database":{
        "block":2258186,//dbLatestBlock
        "confirmedBlock":2258167//dbconfirmedBlock
    },
    "full":{
        "block":2258188//fullNodeBlock
    },
    "solidity":{
        "block":2258170//solidityBlock
    }
}
```
数据来源：<br>
```
1. network->type: 配置文件中  net.type配置项     
2. sync->progress：    
    progress =（ fullNodeProgress + solidityBlockProgress ） / 2    
> fullNodeProgress = (数据库区块高 / GRPC-fullnode接口getNowBlock返回的块高) * 100      
> solidityBlockProgress = (数据库已确认区块高 / GRPC-solidity接口getNowBlock返回块高) * 100     
3. database->block       
    block = 数据库区块高    
4. database->confirmedBlock      
    confirmedBlock = GRPC-solidity接口getNowBlock返回块高    
5. full->block    
    block = GRPC-fullnode接口getNowBlock返回的块高    
6. solidity->block    
    block = GRPC-solidity接口getNowBlock返回块高
```



## 交易所交易信息
- url:/api/market/markets
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/market/markets
```
output:json
```json
[
    {
        "rank":1,
        "name":"Rfinex",
        "pair":"TRX/ETH",
        "link":"https://rfinex.com/",
        "volume":22144662.8099,
        "volumePercentage":19.6793615403,
        "volumeNative":1194868733.76,
        "price":0.0185331343806
    },
    {
        "rank":2,
        "name":"OKEx",
        "pair":"TRX/USDT",
        "link":"https://www.okex.com/market?product=trx_usdt",
        "volume":15633327.7553,
        "volumePercentage":13.8929145869,
        "volumeNative":848455516,
        "price":0.0184256304078
    },
    {
        "rank":3,
        "name":"Binance",
        "pair":"TRX/USDT",
        "link":"https://www.binance.com/en/trade/TRX_USDT",
        "volume":10097034.2529,
        "volumePercentage":8.97296062956,
        "volumeNative":548286518.938,
        "price":0.0184156164782
    },
    {
        "rank":4,
        "name":"BitForex",
        "pair":"TRX/USDT",
        "link":"https://bitforex.com/trade/spotTrading?commodityCode=TRX¤cyCode=USDT",
        "volume":7443829.771,
        "volumePercentage":6.61512972969,
        "volumeNative":406200833.88,
        "price":0.0183254911121
    },
    {
        "rank":5,
        "name":"Bit-Z",
        "pair":"TRX/BTC",
        "link":"https://www.bit-z.com/exchange/trx_btc",
        "volume":7405919.64419,
        "volumePercentage":6.58144002766,
        "volumeNative":401996343.07,
        "price":0.0184228532718
    },
    {
        "rank":6,
        "name":"OKEx",
        "pair":"TRX/BTC",
        "link":"https://www.okex.com/market?product=trx_btc",
        "volume":7279123.30083,
        "volumePercentage":6.46875955452,
        "volumeNative":393774417,
        "price":0.0184855160381
    },
    {
        "rank":7,
        "name":"Huobi",
        "pair":"TRX/USDT",
        "link":"https://www.huobi.pro/ko-kr/trx_usdt/exchange/",
        "volume":6261512.25012,
        "volumePercentage":5.56443619923,
        "volumeNative":339512535.464,
        "price":0.0184426540881
    },
    {
        "rank":8,
        "name":"Binance",
        "pair":"TRX/BTC",
        "link":"https://www.binance.com/en/trade/TRX_BTC",
        "volume":5874751.45811,
        "volumePercentage":5.2207323677,
        "volumeNative":318883908.558,
        "price":0.0184228532718
    },...
]
```
所有数据均从https://coinmarketcap.com/currencies/tron/爬取 每5s获取一次并加载缓存


## 图表信息所用接口
- url:/api/stats/overview
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/stats/overview
```
output:json
```json  
{
    "success":true,
    "data":[
        {
            "date":1529884800000,
            "totalTransaction":3,
            "avgBlockTime":3,
            "avgBlockSize":198,
            "totalBlockCount":1,
            "newAddressSeen":1946,
            "blockchainSize":534
        },
        {
            "date":1529971200000,
            "totalTransaction":3084,
            "avgBlockTime":3,
            "avgBlockSize":188,
            "totalBlockCount":26548,
            "newAddressSeen":203,
            "blockchainSize":5243100
        },
        {
            "date":1530057600000,
            "totalTransaction":5032,
            "avgBlockTime":3,
            "avgBlockSize":190,
            "totalBlockCount":55340,
            "newAddressSeen":154,
            "blockchainSize":10653964
        },...
    ]  
}
```
日交易数图表，区块链大小，日增地址数，区块平均大小


## 申请测试币-TODO
- url:/api/testnet/request-coins
- method:POST

input:json
```json
{
    "address":"TUePpjwtrHtmj2122h74h7R8UqKAV37DhR",
    "captchaCode":"03AL4dnxo8TLilLfyLINe-Om4GeEnwrTNjIdtg6U4agXHxvKQRTFDtv6TsopZB9dh-CdP-vwaAKpwGi98wgrN_9-8J0W6sR86WA7lh1wZxxi10RcVVHumMQD736APcbt-JJltRpHFi5tDULr-_0GZfLEPAozHjrCufJT_nHdHl7aIFEh_qvrBK508o_CEj9dnok0QSiH7vcx86UN398NjKYimJqURdO-I8G76e29iEZqbG9FH-ugNYvOctYLy86CbxKnllHhYq-jBQj0jdIUSE_JfFMTlYv8EYpA"
}
```
output:json
```json  
{
    "success":false,
    "amount":10000000000,
    "code":"TAPOS_ERROR",
    "message":"Tapos check error"
}
```
获取机器IP，从trxRequest表中获取最近一小时这个ip的申请记录，如果存在，code提示“ALREADY_REQUESTED_IP”， message提示：“Already requested TRX from IP recently"             
如果不存在，从trxRequest表根据传入的address查询        
      如果不存在，校验传入的captchaCode（配置文件中的siteCode，post请求   https://www.google.com/recaptcha/api/siteverify，如果返回success才校验通过）      
      校验通过后，从配置文件中那账户和数量信息，生成key，      buildTrxTransfer->broadcastTransaction,将ip和address插入trxRequest      
      校验不通过，code提示“WRONG_CAPTCHA”， message提示：“Wrong Captcha Code"      
如果存在，code提示“ALREADY_REQUESTED_IP”， message提示：“Already requested TRX from IP recently" 



## 验签-TODO
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
输入输出待确认