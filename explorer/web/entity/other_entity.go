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

//Auth 验证签名请求参数
type Auth struct {
	Transaction string `json:"transaction"` //:99.99
}

//AuthResp 验证签名相应
type AuthResp struct {
	Token string `json:"token"` //:99.99
}

//Address 地址签名结构
type Address struct {
	Address string `json:"address"` //
}

//TestCoin  申请测试币 传入参数
type TestCoin struct {
	Address     string `json:"address"`     //
	CaptchaCode string `json:"captchaCode"` //
}

//TestCoinResp  申请测试币 响应参数
type TestCoinResp struct {
	Success bool   `json:"success"`          //:false,
	Amount  int64  `json:"amount,omitempty"` //:10000000000,
	Code    string `json:"code"`             //:"TAPOS_ERROR",
	Message string `json:"message"`          //:"Tapos check error"
}

//VerifyCode 调用google验证接口结构体
type VerifyCode struct {
	Secret   string `json:"secret"`   //
	Response string `json:"response"` //
}

//VerifyCodeResp 调用google验证接口结构体
type VerifyCodeResp struct {
	Success        bool     `json:"success"`          //
	ErrorCodes     []string `json:"error-codes"`      //
	ChallengeTs    string   `json:"challenge_ts"`     //: timestamp,  // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	ApkPackageName string   `json:"apk_package_name"` //:: string,
}
