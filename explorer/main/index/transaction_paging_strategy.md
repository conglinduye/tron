# mysql transaction 分页查询索引优化

## issue: transaction 分页查询慢 
transaction总量较大,分页时如果 limit 起点较大，查询需要借助block_id做大偏移定位，再小偏移取数据

block_id偏移使用 block_id asc 逻辑，transaction基于block_id asc 的定位固定

如 10w为一个索引节点


```sql
-- 获取第一个索引节点的 hash 和 block_id
mysql> select trx_hash, block_id from transactions where block_id >=0 order by block_id asc limit 100000,1;
+------------------------------------------------------------------+----------+
| trx_hash                                                         | block_id |
+------------------------------------------------------------------+----------+
| 663bf96507eda1c2e2f29b67f3f35c4deeaf11c8a6cdfe67643e1c95996515c5 |  2300789 | 第 100001 条
+------------------------------------------------------------------+----------+
1 row in set (0.67 sec)


-- 获取相对块的偏移
mysql> select trx_hash, block_id from transactions where block_id = 2300789 order by block_id asc;
+------------------------------------------------------------------+----------+
| trx_hash                                                         | block_id |
+------------------------------------------------------------------+----------+
| 0a545a45810ec4afd5d2ef1736385280691fa5e1f8b1e8abc72d74f4a01d7047 |  2300789 | 1
| 12927fb1cde200fb4eea594ddb69ffa8b59589e05ad12a685b4e5e992323eac0 |  2300789 | 2
| 1f1c1f4f137bc56f5cc9decd8e4d8ecc1e2960a92efe21b74b81e4f21712756f |  2300789 | 3
| 4051a9ba252eacb7d79c9677792166492ae17cf75f926b59ce79b7a6ca1490fa |  2300789 | 4
| 512e19f40e97477cd5cf30acdbc7de1aa6a82fbd71d98b1ae0e4340f6dd1f676 |  2300789 | 5
| 536f2640d54370a4aba88be0e669392033a2d1d87270767d38c8c2f4b2ff601f |  2300789 | 6   第 100000 条
| 663bf96507eda1c2e2f29b67f3f35c4deeaf11c8a6cdfe67643e1c95996515c5 |  2300789 | 7 * 第 100001 条
| 9493fd9e5570a9b800131512a7f6253c21f594651f8edf559b9feb39c9519946 |  2300789 |
| 98c00fc30000a69dd0a8d014c4a032749c0ca824240bda828248d81369b568a1 |  2300789 |
| d9db4c96f4ab63b59516542748c8082f382ff842c34a62176c581e9ad4388a21 |  2300789 |
| eb4e640f96247cbae57a08e2f76854fe95b88a30e4451628d2de7f69812dcd85 |  2300789 |
+------------------------------------------------------------------+----------+
11 rows in set (0.01 sec)


-- 下一个 100000 条记录为
mysql> select trx_hash, block_id from transactions where block_id >=0 order by block_id asc limit 200000,1;
+------------------------------------------------------------------+----------+
| trx_hash                                                         | block_id |
+------------------------------------------------------------------+----------+
| 90d64a076b11c9376e695d24adb76873c725a53d8b7f4cb31110f3d5c3e1e3ce |  2315969 | 第 200001 条
+------------------------------------------------------------------+----------+
1 row in set (0.67 sec)


-- 使用相对偏移获取第 200001 条
mysql> select trx_hash, block_id from transactions where block_id >= 2300789 order by block_id asc limit 100006,1;
+------------------------------------------------------------------+----------+
| trx_hash                                                         | block_id |
+------------------------------------------------------------------+----------+
| 90d64a076b11c9376e695d24adb76873c725a53d8b7f4cb31110f3d5c3e1e3ce |  2315969 |
+------------------------------------------------------------------+----------+
1 row in set (0.41 sec)

```

计算逻辑固定, 下一个偏移的相对为止为 本次偏移的 block_id 为最小 block_id 的起点，+ 本次偏移hash在当前块的偏移-1 + 固定偏移量
+ 为什么需要 -2 ？
    + 因为block内偏移从1开始算，因此需要 -1 修正; 同时我们需要的第 step 个的位置，而当前取到的是 step+1 个，所以需要再 -1; 因此整体 -2

首次 block_id == 0, 首次偏移hash在块内的偏移量为 0, 固定偏移为 100000，偏移起点为 第 0 次 * 固定偏移量
select trx_hash, block_id from transactions where block_id >=0 order by block_id asc limit 0, 100000

select trx_hash, block_id from transactions where block_id >= ${idx_block_id} order by block_id asc limit  ${round} * ${step} + ${idx_block_offset} - 2, ${step}






## issue: fullnode 内存用量较大问题
待解决


## question: 如何获取witness得票信息

grpc 接口

> https://github.com/tronprotocol/protocol


1. grpc 获取witness当前轮次得票信息
    + 接口: `https://github.com/tronprotocol/protocol/blob/master/api/api.proto` line: 585
    + 推荐 walletSolidity 接口, ListWitness() 接口获取witness信息，votes应该表示当前轮次witness的得票数
2. 从account获取投票信息
    + 接口：`https://github.com/tronprotocol/protocol/blob/master/api/api.proto` line: 566
    + 获取用户的投票信息, votes 字段为候选人地址和投票数的记录。这个逻辑需要扫描所有用户，以便统计最终结果
3. 从交易中获取本轮次的投票交易
    + 解析本轮次所有block中`VoteWitnessContract`类型交易，统计本轮次的投标变化