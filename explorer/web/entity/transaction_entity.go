package entity

//Transactions 查询转账列表的请求参数
type Transactions struct {
	Sort    string `json:"sort,omitempty"`    // 按时间戳倒序
	Limit   int64  `json:"limit,omitempty"`   // 每页记录数
	Count   string `json:"count,omitempty"`   // 是否返回总数
	Start   int64  `json:"start,omitempty"`   // 记录的起始序号
	Number  string `json:"number,omitempty"`  // 按照区块高度精确查询
	Hash    string `json:"hash,omitempty"`    // 按照交易hash精确查询
	Address string `json:"address,omitempty"` // 按照交易精确查询
}

//TransactionsResp 查询转账列表的结果
type TransactionsResp struct {
	Total int64              `json:"total"` // 总记录数
	Data  []*TransactionInfo `json:"data"`  // 记录详情
}

//TransactionInfo 转账信息
type TransactionInfo struct {
	ID           string      `json:"id"`           //uuid
	Block        int64       `json:"block"`        //:2135998,
	Hash         string      `json:"hash"`         //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	CreateTime   int64       `json:"timestamp"`    //:1536314760000,
	OwnerAddress string      `json:"ownerAddress"` //:"JRB1nNvqT6kcRJLdzTnUGyiwvMcnDTAaxYZhTxhvDkjM8kxYh",
	ToAddress    string      `json:"toAddress"`    //:"00000000002097bdd482e26710c054eea72280232a9061885dc94c30c3a0f1b5",
	Data         string      `json:"data"`         //:"", 没用
	ContractType int64       `json:"contractType"` //:1,
	Confirmed    bool        `json:"confirmed"`    //:true
	ContractData interface{} `json:"contractData"` //:原始交易数据
}

//PostTransaction  创建交易
type PostTransaction struct {
	Transaction string `json:"transaction"` // 总记录数
}
