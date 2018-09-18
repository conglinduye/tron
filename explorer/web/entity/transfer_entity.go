package entity

import (
	"time"
)

//Transfers 查询转账列表的请求参数
type Transfers struct {
	Sort    string `json:"sort,omitempty"`    // 按时间戳倒序
	Limit   int64  `json:"limit,omitempty"`   // 每页记录数
	Count   string `json:"count,omitempty"`   // 是否返回总数
	Start   int64  `json:"start,omitempty"`   // 记录的起始序号
	Number  string `json:"number,omitempty"`  // 按照区块高度精确查询
	Hash    string `json:"hash,omitempty"`    // 按照交易hash精确查询
	Address string `json:"address,omitempty"` // 按照交易所属人精确查询
}

//TransfersResp 查询转账列表的结果
type TransfersResp struct {
	Total int64           `json:"total"` // 总记录数
	Data  []*TransferInfo `json:"data"`  // 记录详情
}

//TransferInfo 转账信息
type TransferInfo struct {
	ID                  string    `json:"id"`                  //uuid
	Block               int64     `json:"block"`               //:2135998,
	TransactionHash     string    `json:"transactionHash"`     //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	CreateTime          int64     `json:"timestamp"`           //:1536314760000,
	TransferFromAddress string    `json:"transferFromAddress"` //:"JRB1nNvqT6kcRJLdzTnUGyiwvMcnDTAaxYZhTxhvDkjM8kxYh",
	TransferToAddress   string    `json:"transferToAddress"`   //:"00000000002097bdd482e26710c054eea72280232a9061885dc94c30c3a0f1b5",
	Amount              int64     `json:"amount"`              //:11,
	TokenName           string    `json:"tokenName"`           //:"TRX",
	Confirmed           bool      `json:"confirmed"`           //:true
	LoadTime            time.Time `json:"-"`
}
