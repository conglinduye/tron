## 查询统计信息
- url:/api/stats/overview
- method:get

input:param
```param

eg: 
http://18.216.57.65:20110/api/stats/overview

```

output:json
```json
{
    "success":true,
    "data":[
        {
            "date":1529856000000,     //日期
            "totalTransaction":2850,  //总交易数
            "avgBlockTime":3,         //区块上链平均时间
            "avgBlockSize":207,       //区块平均大小
            "totalBlockCount":16950,  //总区块数
            "newAddressSeen":2218,    //新增地址数
            "blockchainSize":3525089, //区块链大小
            "totalAddress":2248,      //总地址数
            "newBlockSeen":16949,     //新增区块数
            "newTransactionSeen":2847 //新增交易数
        },
        {
            "date":1529942400000,
            "totalTransaction":4658,
            "avgBlockTime":3,
            "avgBlockSize":186,
            "totalBlockCount":45742,
            "newAddressSeen":280,
            "blockchainSize":8906914,
            "totalAddress":2528,
            "newBlockSeen":28792,
            "newTransactionSeen":1808
        }...
    ]
}     
```