package entity


//Token		查询token的请求参数
type Token struct {
	Start  	string `json:"start,omitempty"`  	// 记录的起始序号
	Limit  	string `json:"limit,omitempty"`  	// 每页记录数
	Owner 	string `json:"owner,omitempty"`     // creator_address
	Name 	string `json:"name,omitempty"`		// token_name
	Status  string `json:status,omitempty`		// status
}

//TokenResp	查询token的结果
type TokenResp struct {
	Total	int64				`json:"total"` 	// 总记录数
	Data    []*TokenInfo		`json:"data"`   // 数据
}

// Token 	通证信息
type TokenInfo struct {
	Price 					int64			`json:"price"` 					// 价格
	Issued  				int64			`json:"issued"`					// 流通数量
	IssuedPercentage		float64  		`json:"IssuedPercentage"`  		// 流通数量占比
	Available 				int64			`json:"available"`				// 余额数量
	AvailableSupply			int64			`json:"availableSupply"`		// 流通数量+余额数量
	Remaining				int64			`json:"remaining"`				// 余额数量
	RemainingPercentage 	float64			`json:"remainingPercentage"`	// 余额数量占比
	Percentage				float64  		`json:"percentage"`				// 余额数量占比
	FrozenTotal				int64			`json:"frozenTotal"`			// 冻结数量
	FrozenPercentage		float64			`json:frozenPercentage`			// 冻结数量占比
	TotalParticipateAmount  int64			`json:totalParticipateAmount`	// 已募集金额(sun)

	OwnerAddress			string 			`json:ownerAddress`				// owner_address
	Name  					string 			`json:name`						// asset_name
	TotalSupply				int64			`json:totalSupply`				// 发行总量
	TrxNum					int64			`json:trxNum`					// 通证汇率分子
	Num 					int64			`json:num`						// 通证汇率分母
	EndTime					int64			`json:endTime`					// 结束时间
	StartTime				int64			`json:startTime`				// 开始时间
	VoteScore				int64 			`json:"voteScore"`				// voteScore
	Description				string 			`json:"description"`			// asset_desc
	Url 					string			`json:"url"`					// url
	Frozen				    []TokenFrozen	`json:"Frozen"`					// frozen
	Abbr					string 			`json:"abbr"`					// asset_abbr
	Participated			int64			`json:"participated"`			// participated

}

type TokenFrozen struct {
	Days 	int64
	Amount  int64
}