## 获得数据同步信息-TODO
- url:/api/system/status
- method:get

input:param
```param

```
output:json
```json
      
{
    "network":{
        "type":"mainnet"
    },
    "sync":{
        "progress":99.99988929123896
    },
    "database":{
        "block":2258186,
        "unconfirmedBlock":2258167
    },
    "full":{
        "block":2258188
    },
    "solidity":{
        "block":2258170
    }
}
```


## 交易所交易信息-TODO
- url:/api/market/markets
- method:get

input:param
```param

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

## 图表信息所用接口-TODO
- url:/api/stats/overview
- method:get

input:param
```param

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

```
output:json
```json  


```
输入输出待确认

## 验签-TODO
- url:/api/auth
- method:POST

input:json
```json

```
output:json
```json  


```
输入输出待确认