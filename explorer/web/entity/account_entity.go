package entity

//Accounts 查询账户列表的请求参数
type Accounts struct {
	Sort    string `json:"sort,omitempty"`    // 按照区块高度倒序
	Limit   int64  `json:"limit,omitempty"`   // 每页记录数
	Count   string `json:"count,omitempty"`   // 是否返回总数
	Start   int64  `json:"start,omitempty"`   // 记录的起始序号
	Address string `json:"address,omitempty"` // 按照地址精确查询
}

//AccountsResp 查询账户列表的结果
type AccountsResp struct {
	Total int64          `json:"total"` // 总记录数
	Data  []*AccountInfo `json:"data"`  // 记录详情
}

//AccountInfo 账户信息
type AccountInfo struct {
	Address       string           `json:"address"`       //:TDtjQ1JR5UrS92W9kB6BCeAQJwn1dyBEbs,
	Name          string           `json:"name"`          //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	Balance       int64            `json:"balance"`       //:3006,
	Power         int64            `json:"power"`         //:"JRB1nNvqT6kcRJLdzTnUGyiwvMcnDTAaxYZhTxhvDkjM8kxYh",
	TokenBalances map[string]int64 `json:"tokenBalances"` //:"00000000002097bdd482e26710c054eea72280232a9061885dc94c30c3a0f1b5",
	UpdateTime    int64            `json:"dateUpdated"`   //:1536314760000
	CreateTime    int64            `json:"dateCreated"`   //:1536314760000,
}

//AccountDetail 账户详细信息
type AccountDetail struct {
	Representative *Represent     `json:"representative"` //
	Name           string         `json:"name"`           //
	Address        string         `json:"address"`        //
	Bandwidth      *BandwidthInfo `json:"bandwidth"`      //
	Balances       []*Balance     `json:"balances"`       //
	Balance        int64          `json:"balance"`        //
	TokenBalances  []*Balance     `json:"tokenBalances"`  //
	Frozen         *Frozen        `json:"frozen"`         //
}

//Represent 。。。
type Represent struct {
	Enabled          bool   `json:"enabled"`          //:TDtjQ1JR5UrS92W9kB6BCeAQJwn1dyBEbs,
	LastWithDrawTime int64  `json:"lastWithDrawTime"` //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	Allowance        int64  `json:"allowance"`        //:3006,
	URL              string `json:"url"`              //:3006,
}

//BandwidthInfo ...
type BandwidthInfo struct {
	FreeNetUsed       int64                 `json:"freeNetUsed"`       //:
	FreeNetLimit      int64                 `json:"freeNetLimit"`      //:3006,
	FreeNetRemaining  int64                 `json:"freeNetRemaining"`  //:
	FreeNetPercentage float64               `json:"freeNetPercentage"` //:3006,
	NetUsed           int64                 `json:"netUsed"`           //:
	NetLimit          int64                 `json:"netLimit"`          //:3006,
	NetRemaining      int64                 `json:"netRemaining"`      //:
	NetPercentage     float64               `json:"netPercentage"`     //:3006,
	Assets            map[string]*AssetInfo `json:"assets"`            //:
}

//AssetInfo ...
type AssetInfo struct {
	NetUsed       int64   `json:"netUsed"`       //:
	NetLimit      int64   `json:"netLimit"`      //:3006,
	NetRemaining  int64   `json:"netRemaining"`  //:
	NetPercentage float64 `json:"netPercentage"` //:3006,
}

//Balance ...
type Balance struct {
	Name    string  `json:"name"`    //:3006,
	Balance float64 `json:"balance"` //:3006,
}

//Frozen ...
type Frozen struct {
	Total    int64          `json:"total"`    //:
	Balances []*BalanceInfo `json:"balances"` //:
}

//BalanceInfo ...
type BalanceInfo struct {
	Expires int64 `json:"expires"` //:
	Amount  int64 `json:"amount"`  //:
}

//BalanceInfoDB 跟BalanceInfo，数据库存储的跟api返回的不一致，转换
type BalanceInfoDB struct {
	Expires int64 `json:"expire_time"`    //:
	Amount  int64 `json:"frozen_balance"` //:
}

//SuperAccountInfo ...
type SuperAccountInfo struct {
	ID         int64  `json:"id,omitempty"` //:
	Address    string `json:"address"`      //:
	GithubLink string `json:"githubLink"`   //:
}

//AccountMediaInfo ...
type AccountMediaInfo struct {
	Success bool   `json:"sucess"` //:
	Image   string `json:"image"`  //:
	Reason  string `json:"reason"` //:
}

//AccountTransactionNum ...
type AccountTransactionNum struct {
	Transactions    int64 `json:"transactions"`     //:: "827",
	TransactionsOut int64 `json:"transactions_out"` //:: "230",
	TransactionIn   int64 `json:"transaction_in"`   //:: "597"
}
