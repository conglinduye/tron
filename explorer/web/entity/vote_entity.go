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
	VoterAddress        string 	`json:"voterAddress"`        // 投票人地址
	CandidateAddress    string 	`json:"candidateAddress"`    // 被投票人(候选人)地址
	Votes               int64  	`json:"votes"`               // 票数
	CandidateURL        string 	`json:"candidateUrl"`        // 候选人url
	CandidateName       string 	`json:"candidateName"`       // 候选人名称
	VoterAvailableVotes float64 `json:"voterAvailableVotes"` // 投票人可用票数
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

//AddressVotes
type AddressVotes struct {
	Votes					map[string]int64		`json:"votes"`
}


//VotesInfo
type AddressVoteInfo struct {
	VoterAddress        string 	`json:"voterAddress"`        // 投票人地址
	CandidateAddress    string 	`json:"candidateAddress"`    // 被投票人(候选人)地址
	Votes               int64  	`json:"votes"`               // 票数
}

type VoteLiveResp struct {
	Data				map[string]*VoteLive 	`json:"data"`
}

type VoteLive struct {
	Address 			string 	`json:"address"`
	Votes				int64	`json:"votes"`
}

type VoteCurrentCycleResp struct {
	TotalVotes			int64					`json:"total_votes"`
	Candidates    		[]*VoteCurrentCycle		`json:"candidates"`
}


type VoteCurrentCycle struct {
	Address     		string 		`json:"address"`
	Name        		string 		`json:"name"`
	URL         		string 		`json:"url"`
	HasPage     		bool   		`json:"hasPage"`
	Votes      			int64  		`json:"votes"`
	RealTimeVotes		int64 		`json:"realTimeVotes"`
	ChangeCycle 		int32  		`json:"change_cycle"`
	ChangeDay			int32		`json:"change_day"`
}


