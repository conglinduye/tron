package entity

//SrAccount  查询超级代表的请求参数
type SrAccount struct {
	Limit  	string `json:"limit,omitempty"`  	// 每页记录数
	Start  	string `json:"start,omitempty"`  	// 记录的起始序号
	Address	string `json:"address,omitempty"`	// 地址查询
}

//SrAccountResp 查询超级代表的结果
type SrAccountResp struct {
	Total	int64				`json:"total"` 	// 总记录数
	Data    []*SrAccountInfo	`json:"data"`   // 数据
}

//SrAccount 超级代表信息
type SrAccountInfo struct {
	Address 		string	`json:"address"`		// 地址
	GithubLink		string 	`json:"githubLink"`		// github链接
	CreateTime 		string	`json:"createTime"`		// 创建时间
	ModifiedTime	string	`json:"modifiedTime"`   // 修改时间
}
