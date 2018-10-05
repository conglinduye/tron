package entity

//CreateAccount 创建账户的返回结果
type CreateAccount struct {
	Key     string `json:"key"`     // 私钥
	Address string `json:"address"` // 地址
}
