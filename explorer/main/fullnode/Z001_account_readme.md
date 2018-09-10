# 逻辑
1. 解析新到的 transaction，将其中交易涉及的地址加入 redis 的 set集合
2. account 主逻辑定期从redis的set集合中获取所有发生交易的账户，调用主网的getaccount获取对应账户的最新状态
3. transaction解析将生成对应的contract数据并保存