package entity

//SystemStatusResp 查询转账列表的结果
type SystemStatusResp struct {
	Network  *Network   `json:"network"`  //
	Sync     *Sync      `json:"sync"`     //
	Database *DataBase  `json:"database"` //
	Full     *BlockNode `json:"full"`     //
	Solidity *BlockNode `json:"solidity"` //
}

//Network 主网信息
type Network struct {
	Type string `json:"type"` //:"mainnet"
}

//Sync 同步进度信息
type Sync struct {
	Progress float64 `json:"progress"` //:99.99
}

//DataBase 数据库block
type DataBase struct {
	Block          int64 `json:"block"`          //:"2258186"
	ConfirmedBlock int64 `json:"confirmedBlock"` //:"2258167"
}

//BlockNode 主网block
type BlockNode struct {
	Block int64 `json:"block"` //:"2258188"
}

//MarketInfo 交易所信息
type MarketInfo struct {
	Rank             int64   `json:"rank"`             //:1,
	Name             string  `json:"name"`             //:"Rfinex",
	Pair             string  `json:"pair"`             //:"TRX/ETH",
	Link             string  `json:"link"`             //:"https://rfinex.com/",
	Volume           float64 `json:"volume"`           //:22144662.8099,
	VolumePercentage float64 `json:"volumePercentage"` //:19.6793615403,
	VolumeNative     float64 `json:"volumeNative"`     //:1194868733.76,
	Price            float64 `json:"price"`            //:0.0185331343806
}
