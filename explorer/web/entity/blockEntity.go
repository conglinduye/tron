package entity

//Blocks 查询区块列表的请求参数
type Blocks struct {
	Sort   string `json:"sort,omitempty"`   // 按照区块高度倒序
	Limit  string `json:"limit,omitempty"`  // 每页记录数
	Count  string `json:"count,omitempty"`  // 是否返回总数
	Start  string `json:"start,omitempty"`  // 记录的起始序号
	Order  string `json:"order,omitempty"`  // 按时间戳倒序
	Number string `json:"number,omitempty"` // 按照区块高度精确查询
}

//BlocksResp 查询区块列表的结果
type BlocksResp struct {
	Total int64        `json:"total"` // 总记录数
	Data  []*BlockInfo `json:"data"`  // 记录详情
}

//BlockInfo 区块信息
type BlockInfo struct {
	Number         int64  `json:"number"`         //:2135998,
	Hash           string `json:"hash"`           //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	Size           int64  `json:"size"`           //:3006,
	CreateTime     int64  `json:"timestamp"`      //:1536314760000,
	TxTrieRoot     string `json:"txTrieRoot"`     //:"JRB1nNvqT6kcRJLdzTnUGyiwvMcnDTAaxYZhTxhvDkjM8kxYh",
	ParentHash     string `json:"parentHash"`     //:"00000000002097bdd482e26710c054eea72280232a9061885dc94c30c3a0f1b5",
	WitnessID      int32  `json:"witnessId"`      //:0,
	WitnessAddress string `json:"witnessAddress"` //:"TSNbzxac4WhxN91XvaUfPTKP2jNT18mP6T",
	NrOfTrx        int64  `json:"nrOfTrx"`        //:11,
	Confirmed      bool   `json:"confirmed"`      //:true
}
