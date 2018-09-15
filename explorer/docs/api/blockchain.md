## 查询区块信息
- url:/api/block
- method:get

input:param
```param
sort=-number    //按照区块高度倒序
&limit=40       //每页40条记录
&count=true     //是否返回总数
&start=0        //记录的起始序号
&order=-timestamp       //按照时间戳倒序排列
&number=2170015         //按照区块高度精确查询

eg: http://18.216.57.65:20110/api/block?sort=-number&limit=40&start=0
```
output:json
```json
{
    "total":2169998,
    "data":[
        {
            "number":2170015,//区块高度
            "hash":"0000000000211c9fb87d9cf193db8326349148d32ede34d7c5ac0bee92b22374",//区块hash
            "size":2815,//区块大小
            "timestamp":1536416862000,//时间戳
            "txTrieRoot":"2NT3va4sQGdoHShLFZdyW9jApfCj3AmvyiYLEfnn1NkvMvvewf",//验证数根的hash值
            "parentHash":"0000000000211c9e71b4f0d007d5a6eb508d05dd515c0a82373d685c6325e7ba",//父hash
            "witnessId":0,//所属候选人ID
            "witnessAddress":"TFA1qpUkQ1yBDw4pgZKx25wEZAqkjGoZo1",//所属候选人地址
            "nrOfTrx":12,//交易数
            "confirmed":false//是否经过确认
        },
        {
            "number":2170014,
            "hash":"0000000000211c9e71b4f0d007d5a6eb508d05dd515c0a82373d685c6325e7ba",
            "size":2337,
            "timestamp":1536416859000,
            "txTrieRoot":"22d9oSJenHi5ktjkhpChxdaRxcj3B39du4D37CkNySDGSYtEd5",
            "parentHash":"0000000000211c9df251e69b9e1341d89df9744c8febd6695565fec60bfc9738",
            "witnessId":0,
            "witnessAddress":"TEKUPpjTMKWw9LJZ9YJ4enhCjAmVXSL7M6",
            "nrOfTrx":10,
            "confirmed":false
        },...
     ]
}
```
可与首页的区块高度，出块记录列表，搜索区块功能复用

## 单个区块信息
- url:/api/block/:number
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/block/2341043
```
output:json
```json
{
    "number":2170015,
    "hash":"0000000000211c9fb87d9cf193db8326349148d32ede34d7c5ac0bee92b22374",
    "size":2815,
    "timestamp":1536416862000,
    "txTrieRoot":"2NT3va4sQGdoHShLFZdyW9jApfCj3AmvyiYLEfnn1NkvMvvewf",
    "parentHash":"0000000000211c9e71b4f0d007d5a6eb508d05dd515c0a82373d685c6325e7ba",
    "witnessId":0,
    "witnessAddress":"TFA1qpUkQ1yBDw4pgZKx25wEZAqkjGoZo1",
    "nrOfTrx":12,
    "confirmed":true
}
```
