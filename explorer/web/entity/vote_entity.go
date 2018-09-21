package entity

//Votes 查询投票列表的请求参数
type Votes struct {
	Sort      string `json:"sort,omitempty"`      // 按时间戳倒序
	Limit     int64  `json:"limit,omitempty"`     // 每页记录数
	Count     string `json:"count,omitempty"`     // 是否返回总数
	Start     int64  `json:"start,omitempty"`     // 记录的起始序号
	Candidate string `json:"candidate,omitempty"` // 按照候选人精确查询
	Voter     string `json:"voter,omitempty"`     // 按照投票人精确查询
}

//VotesResp 查询投票列表的结果
type VotesResp struct {
	Total      int64        `json:"total"`      // 总记录数
	TotalVotes int64        `json:"totalVotes"` // 总投票数
	Data       []*VotesInfo `json:"data"`       // 记录详情
}

//VotesInfo  投票信息
type VotesInfo struct {
	ID                  string 	`json:"id"`                  //uuid
	Block               int64  	`json:"block"`               //:2135998,
	Transaction         string 	`json:"transaction"`         //:"00000000002097beb4b9ceabbff396bf788a8d9ee8c09de37e5e0da039a6a87f",
	CreateTime          int64  	`json:"timestamp"`           //:1536314760000,
	VoterAddress        string 	`json:"voterAddress"`        //:"JRB1nNvqT6kcRJLdzTnUGyiwvMcnDTAaxYZhTxhvDkjM8kxYh",
	CandidateAddress    string 	`json:"candidateAddress"`    //:"00000000002097bdd482e26710c054eea72280232a9061885dc94c30c3a0f1b5",
	Votes               int64  	`json:"votes"`               //:11,
	CandidateURL        string 	`json:"candidateUrl"`        //:"TRX",
	CandidateName       string 	`json:"candidateName"`       //:"TRX",
	VoterAvailableVotes float64 `json:"voterAvailableVotes"` //:10
}

//VoteLiveInfo 实时投票数据
type VoteLiveInfo struct {
	Data map[string]*LiveInfo `json:"data"` // 记录详情
}

//LiveInfo 实时信息
type LiveInfo struct {
	Address string `json:"address"` //:"TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir",
	Name    string `json:"name"`    //:"trongalaxy",
	URL     string `json:"url"`     //:"http://www.trongalaxy.io",
	Votes   int64  `json:"votes"`   //:100006481
}

//VoteCurrentCycleResp 上轮投票信息
type VoteCurrentCycleResp struct {
	TotalVotes int64               `json:"total_votes"` //:"TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir",
	Candidates []*VoteCurrentCycle `json:"candidates"`  //:"trongalaxy",
}

//VoteCurrentCycle 上轮投票信息
type VoteCurrentCycle struct {
	Address     string `json:"address"`      //:"TFuC2Qge4GxA2U9abKxk1pw3YZvGM5XRir",
	Name        string `json:"name"`         //:"trongalaxy",
	URL         string `json:"url"`          //:"http://www.trongalaxy.io",
	HasPage     bool   `json:"hasPage"`      //:100006481
	Votes       int64  `json:"votes"`        //:100006481
	ChangeCycle int32  `json:"change_cycle"` //:1,
	ChangeDay   int32  `json:"change_day"`   //:1,
}

//VoteNextCycleResp 返回倒计时时间
type VoteNextCycleResp struct {
	NextCycle int64 `json:"nextCycle"` //:,毫秒
}

// AccountVoteResultRes
type AccountVoteResultRes struct {
	Total      int64        `json:"total"`      		// total
	Data       []*AccountVoteResult `json:"data"`       // data
}

//AccountVoteResult
type AccountVoteResult struct {
	Address		string 	`json:"address"`	// address
	ToAddress	string 	`json:"toAddress"`  // toAddress
	Vote		int64   `json:"vote"`		// vote
}

// CandidateInfo
type CandidateInfo struct {
	CandidateAddress	string 	`json:"candidateAddress"`		// candidateAddress
	CandidateName		string 	`json:"candidateName"`			// candidateName
	CandidateUrl		string 	`json:"candidateUrl"`			// candidateUrl
}

//VoteWitnessReq
type VoteWitnessReq struct {
	Limit     string  `json:"limit,omitempty"`     	// 每页记录数
	Start     string  `json:"start,omitempty"`     	// 记录的起始序号
	Address   string  `json:"address,omitempty"` 	// 地址

}


//VoteWitnessResp
type VoteWitnessResp struct {
	Total				int64					`json:"total"`
	TotalVotes			int64					`json:"totalVotes"`
	Data    			[]*VoteWitness			`json:"data"`
	FastestRise			*VoteWitness			`json:"fastestRise"`
}

//VoteWitness 节点投票信息
type VoteWitness struct {
	RealTimeRanking		int32		`json:"realTimeRanking"`
	Address     		string 		`json:"address"`
	Name        		string 		`json:"name"`
	URL         		string 		`json:"url"`
	HasPage     		bool   		`json:"hasPage"`
	LastCycleVotes      int64  		`json:"lastCycleVotes"`
	RealTimeVotes		int64 		`json:"realTimeVotes"`
	ChangeVotes			int64		`json:"changeVotes"`
	VotesPercentage		float64		`json:"votesPercentage"`
	ChangeCycle 		int32  		`json:"change_cycle"`
	ChangeDay   		int32  		`json:"change_day"`
}

// VoteWitnessRanking
type VoteWitnessRanking struct {
	Address     		string 		`json:"address"`
	Ranking				int32		`json:ranking`
}

//VoteWitnessDetail
type VoteWitnessDetail struct {
	Success 				bool					`json:"success"`
	Data 					*VoteWitness			`json:"data"`

}
