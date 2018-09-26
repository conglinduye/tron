## 查询节点投票列表信息
- url:/api/vote/witness
- method:get

input:param
```param
eg:
http://18.216.57.65:20110/api/vote/witness
```
output:json
```json

{
    "total":143,
    "totalVotes":7852792991,
     "fastestRise":{ //排名上升变化最快的节点
            "realTimeRanking":46,
            "address":"TJgmwx9TYaqujmdthJkjaLyWXrwTCmmTan",
            "name":"ZADEA-MadeInItaly",
            "url":"www",
            "hasPage":true,
            "lastCycleVotes":84350,
            "realTimeVotes":2240028,
            "changeVotes":2155678,
            "votesPercentage":0.001074783688738554,
            "change_cycle":56//前端页面拿该值显示变化位数
        }
    "data":[
        { 
            "realTimeRanking":1,//实时投票排名
            "address":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//地址
            "name":"Sesameseed",//名称
            "url":"https://www.sesameseed.org",//url
            "hasPage":true,//是否有page
            "lastCycleVotes":510451563,//上一轮投票数
            "realTimeVotes":508640696,//实时投票数
            "changeVotes":-1810867,//变化票数
            "votesPercentage":6.500254923121276,//投票占比
            "change_cycle":1,//6小时票数排位变化
        },
        {
            "address":"TV6qcwSp38uESiDczxxb7zbJX1h2LfDs78",
            "name":"TronsTronics",
            "url":"https://tronstronics.com",
            "hasPage":false,
            "lastCycleVotes":466345134,
            "realTimeVotes":463860453,
            "changeVotes":-2484681,
            "votesPercentage":5.9385894233360416,
            "change_cycle":0
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
        "realTimeRanking":1,
        "address":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",
        "name":"Sesameseed",
        "url":"https://www.sesameseed.org",
        "hasPage":true,
        "lastCycleVotes":505712168,
        "realTimeVotes":506759341,
        "changeVotes":1047173,
        "votesPercentage":0.06443760395530662,
        "change_cycle":0
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
            "voterAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//投票人地址
            "candidateAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//被投票人(候选人)地址
            "votes":10,//票数
            "candidateUrl":"https://www.sesameseed.org",//候选人url
            "candidateName":"Sesameseed",//候选人名称
            "voterAvailableVotes":10// 投票人可用票数
        },
        {
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

## 根据地址查看投票(投给谁)信息
- url:/api/account/:address/votes
- method:GET

input:param
```param
eg: 
http://18.216.57.65:20110/api/account/TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp/votes
```
output:json
```json
{
    "votes":{
        "TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp":10
    }
}
```



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


