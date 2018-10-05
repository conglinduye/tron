## 创建账户
- url:/api/account
- method:POST

input:param
```param

eg: http://18.216.57.65:20110/api/account POST

```
output:json
```json
{
    "key":"cc9fd97198c6072729fa5df0159d5607d4e3da03a92d0c24eb89a9f07f43539d", "address":"TCWkTNVAV5QErki3ERxkYsMjkzWWXaQCB6"
}
```


## 获取账户余额
- url:/api/account/:address/balance
- method:GET

input:param
```param

eg: http://18.216.57.65:20110/api/account/TCWkTNVAV5QErki3ERxkYsMjkzWWXaQCB6/balance

```
output:json
```json
{
    "allowance":43173868649,
    "entropy":0,
    "balances":[
        {
            "name":"TRX",
            "balance":1660.428855
        },
        {
            "name":"Tarquin",
            "balance":1
        },
        {
            "name":"SexTronsPartys",
            "balance":6669
        },
        {
            "name":"GoodLuck",
            "balance":8
        },
        {
            "name":"NBACoin",
            "balance":10
        },...
    ],
    "frozen":{
        "total":10000000,
        "balances":[
            {
                "amount":10000000,
                "expires":1533339042000
            }
        ]
    }
}
```

## 获取账户sr网页
- url:/api/account/:address/sr-pages
- method:GET

input:param
```param

eg: http://18.216.57.65:20110/api/account/TCWkTNVAV5QErki3ERxkYsMjkzWWXaQCB6/sr-pagess

```
output:json
```json
{
    "intro": "400: Invalid request\n",
    "communityPlan": "400: Invalid request\n",
    "team": "400: Invalid request\n",
    "budgetExpenses": "400: Invalid request\n",
    "serverConfiguration": "400: Invalid request\n"
}
```
