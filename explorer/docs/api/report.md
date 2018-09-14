## 查询统计信息
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
            "date":1529856000000,
            "totalTransaction":2850,
            "avgBlockTime":3,
            "avgBlockSize":207,
            "totalBlockCount":16950,
            "newAddressSeen":2218,
            "blockchainSize":3525089,
            "totalAddress":2248,
            "newBlockSeen":16949,
            "newTransactionSeen":2847
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