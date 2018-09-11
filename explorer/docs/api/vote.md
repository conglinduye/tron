## 查询投票信息
- url:/api/vote
- method:get

input:param
```param
&limit=40       //每页40条记录
&count=true     //是否返回总数
&start=0        //记录的起始序号
&sort=-votes       //按照投票数量倒序排列
&candidate=12 //按候选人地址查询
&voter=13 //按投票人地址查询
```
output:json
```json
      
{
    "total":1,//总记录数
    "totalVotes":15000010,//总票数
    "data":[
        {
            "id":"0ed12410-41a7-4caf-a49f-33dae029053e",// 投票ID
            "block":1836424,//交易块
            "transaction":"a7fe8203ed520bbc12a406f8ec5ed8ec9fb2d262967b2212581971ab8cc95cfd",//交易hash
            "timestamp":"2018-08-27T23:25:12.000Z",//时间戳
            "voterAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//投票人地址
            "candidateAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//被投票人地址
            "votes":10,//票数
            "candidateUrl":"https://www.sesameseed.org",//候选人url
            "candidateName":"Sesameseed",//候选人名称
            "voterAvailableVotes":10// 可用票数
        },
        {
            "id":"1acaaac8-497c-490a-a536-acbb8d18562f",
            "block":1712572,
            "transaction":"b3d2887f14f03f2b6ed3cc41dae382063b55a615ee05e316670ebb24d626bdd2",
            "timestamp":"2018-08-23T16:03:21.000Z",
            "voterAddress":"TR6vbWGkT9ztWkKoG972Uzz529niHAEpGp",
            "candidateAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",
            "votes":15000000,
            "candidateUrl":"https://www.sesameseed.org",
            "candidateName":"Sesameseed",
            "voterAvailableVotes":60000000
        },...
    ]
}
```


## 所有代表的实时投票信息
- url:/api/vote/live
- method:get

input:param
```param

```
output:json
```json
{
    "data":{
        "TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV":{//代表地址
            "address":"TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV",//代表地址
            "name":"trongalaxy",
            "url":"http://www.trongalaxy.io",
            "votes":100001281//票数
        },
        "TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir":{
            "address":"TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir",
            "name":"trongalaxy",
            "url":"http://www.trongalaxy.io",
            "votes":100006481
        },
        "TRXDEXMoaAprSGJSwKanEUBqfQjvQEDuaw":{
            "address":"TRXDEXMoaAprSGJSwKanEUBqfQjvQEDuaw",
            "name":"trongalaxy",
            "url":"http://www.trongalaxy.io",
            "votes":4808866
        },...
}
```

## 上一个投票周期的投票情况
- url:/api/vote/current-cycle 
- method:get

input:param
```param

```
output:json
```json
{
    "total_votes":7664305937,//总投票数
    "candidates":[
        {
            "address":"TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV",//地址
            "name":"",//名字
            "url":"http://TronGr10.com",//rul
            "hasPage":false,//是否有github地址
            "votes":100001281,//得票数
            "change_cycle":1,//6小时投票变化
            "change_day":1//12小时投票变化
        },
        {
            "address":"TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir",
            "name":"",
            "url":"http://TronGr11.com",
            "hasPage":false,
            "votes":100006481,
            "change_cycle":3,
            "change_day":3
        },
        {
            "address":"TRXDEXMoaAprSGJSwKanEUBqfQjvQEDuaw",
            "name":"TrxDexCom",
            "url":"https://TrxDex.com",
            "hasPage":true,
            "votes":4808866,
            "change_cycle":0,
            "change_day":0
        },...
    ]
}
```

## 返回倒计时时间
- url:/api/vote/next-cycle
- method:get

input:param
```param

```
output:json
```json
{
    "nextCycle": 9891000 // 剩余时长，单位ms
}
```
投票每六小时一个周期<br>
投票时段：（北京时间）   2:00  8:00  14:00  20:00


