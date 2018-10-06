package entity

//CreateAccount 创建账户的返回结果
type CreateAccount struct {
	Key     string `json:"key"`     // 私钥
	Address string `json:"address"` // 地址
}

//AccountBalance 账户余额信息
type AccountBalance struct {
	Allowance int64       `json:"allowance"` //
	Entropy   int64       `json:"entropy"`   //
	Balances  []*Balance  `json:"balances"`  //
	Frozen    *FrozenInfo `json:"frozen"`    //
}

//Balance 账户信息
type Balance struct {
	Name    string  `json:"name"`    //
	Balance float64 `json:"balance"` //
}

//FrozenInfo 冻结信息
type FrozenInfo struct {
	Total    int64            `json:"total"`    //
	Balances []*FrozenBalance `json:"balances"` //
}

//FrozenBalance  冻结详情
type FrozenBalance struct {
	Amount  int64 `json:"amount"`  //
	Expires int64 `json:"expires"` //
}

//BalanceInfoDB 跟FrozenBalance，数据库存储的跟api返回的不一致，转换
type BalanceInfoDB struct {
	Expires int64 `json:"expire_time"`    //:
	Amount  int64 `json:"frozen_balance"` //:
}
