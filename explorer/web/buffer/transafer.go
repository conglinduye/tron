package buffer

/*
## transaction 查询方式
1. 最新列表，分页
2. block关联 --> blockID
3. 用户关联
	1. owner --> owner_address
	2. to (only for transfer) --> to_address
	3. 用户交易数量 (transfer 统计: owner_address, to_address count)
	4. 交易总数 (transaction count, owner_address)

*/
