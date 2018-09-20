## 查询节点投票列表信息
- url:/api/vote/witness
- method:get

input:param

eg:
http://18.216.57.65:20110/api/vote/witness
```
output:json
```json

{
    "total":143,
    "totalVotes":7852792991,
     "fastestRise":{ //排名上升变化最快的节点
            "address":"TJgmwx9TYaqujmdthJkjaLyWXrwTCmmTan",
            "name":"ZADEA-MadeInItaly",
            "url":"www",
            "hasPage":true,
            "lastCycleVotes":84350,
            "realTimeVotes":2240028,
            "changeVotes":2155678,
            "votesPercentage":0.00001074783688738554,
            "change_cycle":56,//前端页面拿该值显示变化位数
            "change_day":56
        }
    "data":[
        {
            "address":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//地址
            "name":"Sesameseed",//名称
            "url":"https://www.sesameseed.org",//url
            "hasPage":true,//是否有page
            "lastCycleVotes":510451563,//上一轮投票数
            "realTimeVotes":508640696,//实时投票数
            "changeVotes":-1810867,//变化票数
            "votesPercentage":0.06500254923121276,//投票占比
            "change_cycle":1,//6小时票数排位变化
            "change_day":1 //12小时票数排位变化
        },
        {
            "address":"TV6qcwSp38uESiDczxxb7zbJX1h2LfDs78",
            "name":"TronsTronics",
            "url":"https://tronstronics.com",
            "hasPage":false,
            "lastCycleVotes":466345134,
            "realTimeVotes":463860453,
            "changeVotes":-2484681,
            "votesPercentage":0.059385894233360416,
            "change_cycle":0,
            "change_day":0
        },...
    ]
}
```

## 查询节点投票信息
- url:/api/vote/witness/:address
- method:get

input:param
```param

eg:
http://18.216.57.65:20110/api/vote/witness/TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp
```
output:json
```json
      
{
    "success":true,
    "data":{
        "address":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",
        "name":"Sesameseed",
        "url":"https://www.sesameseed.org",
        "hasPage":true,
        "lastCycleVotes":505712168,
        "realTimeVotes":506759341,
        "changeVotes":1047173,
        "votesPercentage":0.06443760395530662,
        "change_cycle":0,
        "change_day":0
    }
}
```

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

eg:
http://18.216.57.65:20110/api/vote?sort=-number&limit=40&start=0

http://18.216.57.65:20110/api/vote?candidate=TAEw4zwwYMiDcWFC9xQrLP9moMi34YAZbz

http://18.216.57.65:20110/api/vote?voter=TAEw4zwwYMiDcWFC9xQrLP9moMi34YAZbz
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

缓存策略：
程序运行初次加载数据到内存；     
缓存数据每隔30s更新一次；    
如果缓存中无数据，则触发重新加载


## 所有代表的实时投票信息
- url:/api/vote/live
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/vote/live
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
缓存策略：
程序运行初次加载数据到内存；     
缓存数据每隔60s更新一次；    
如果缓存中无数据，则触发重新加载

## 上一个投票周期的投票情况
- url:/api/vote/current-cycle 
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/vote/current-cycle
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
缓存策略：
程序运行初次加载数据到内存；     
缓存数据每隔60s更新一次；    
如果缓存中无数据，则触发重新加载

## 返回倒计时时间
- url:/api/vote/next-cycle
- method:get

input:param
```param
eg: http://18.216.57.65:20110/api/vote/next-cycle
```
output:json
```json
{
    "nextCycle": 9891000 // 剩余时长，单位ms
}
```
从主网调用接口获取下次投票周期时间nextMaintenanceTime
从主网调用接口获取最新块的时间戳 curTimeStamp
则倒计时时间=nextMaintenanceTime-curTimeStamp
缓存策略：
程序运行初次加载主网数据到内存；     
缓存数据每隔60s更新一次；    
如果缓存中无数据，则触发重新加载


