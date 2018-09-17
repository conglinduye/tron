package module

import (
	"encoding/json"
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryTotalVotes
func QueryTotalVotes() int64 {
	strSQL := fmt.Sprintf(`
	SELECT sum(vote_count) as totalVotes FROM tron.witness`)
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryVoteLiveRealize error :[%v]\n", err)
		return 0
	}
	if dataPtr == nil {
		log.Errorf("QueryVoteLiveRealize dataPtr is nil ")
		return 0
	}
	var votes = int64(0)
	for dataPtr.NextT() {
		votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalVotes"))
	}

	return votes
}

//QueryVoteLiveRealize 操作数据库
func QueryVoteLiveRealize(strSQL string) (*entity.VoteLiveInfo, error) {
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryVoteLiveRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryVoteLiveRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	votesResp := &entity.VoteLiveInfo{}
	votesMap := make(map[string]*entity.LiveInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var vote = &entity.LiveInfo{}
		vote.Address = dataPtr.GetField("voteraddress")
		vote.Name = dataPtr.GetField("account_name")
		vote.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("votes"))
		vote.URL = dataPtr.GetField("url")
		votesMap[vote.Address] = vote
	}
	votesResp.Data = votesMap

	return votesResp, nil

}

//QueryVoteCurrentCycleRealize 操作数据库
func QueryVoteCurrentCycleRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.VoteCurrentCycleResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryVoteCurrentCycleRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryVoteCurrentCycleRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	votesResp := &entity.VoteCurrentCycleResp{}
	voteCycle := make([]*entity.VoteCurrentCycle, 0)
	var totalVotes = int64(0)

	//填充数据
	for dataPtr.NextT() {
		var vote = &entity.VoteCurrentCycle{}
		githubLink := dataPtr.GetField("github_link")
		if githubLink != "" {
			vote.HasPage = true
		}

		vote.Address = dataPtr.GetField("voteraddress")
		vote.Name = dataPtr.GetField("account_name")
		vote.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("votes"))
		vote.URL = dataPtr.GetField("url")

		voteCycle = append(voteCycle, vote)
	}
	totalVotes = QueryTotalVotes()

	votesResp.TotalVotes = totalVotes
	votesResp.Candidates = voteCycle

	return votesResp, nil

}

//QueryAccountVoteResultRealize
func QueryAccountVoteResultRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AccountVoteResultRes, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountVoteResultRealize error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountVoteResultRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	accountVoteResultRes := &entity.AccountVoteResultRes{}
	accountVoteResults := make([]*entity.AccountVoteResult, 0)

	for dataPtr.NextT() {
		accountVoteResult := &entity.AccountVoteResult{}
		accountVoteResult.Address = dataPtr.GetField("address")
		accountVoteResult.ToAddress = dataPtr.GetField("to_address")
		accountVoteResult.Vote = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote"))

		accountVoteResults = append(accountVoteResults, accountVoteResult)
	}

	var total = int64(len(accountVoteResults))
	total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}

	accountVoteResultRes.Total = total
	accountVoteResultRes.Data = accountVoteResults

	return accountVoteResultRes, nil
}

// QueryCandidateInfo
func QueryCandidateInfo(strSQL string) (*entity.CandidateInfo, error) {
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryCandidateInfo error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryCandidateInfo dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	candidateInfo := &entity.CandidateInfo{}

	for dataPtr.NextT() {
		candidateInfo.CandidateAddress = dataPtr.GetField("candidateAddress")
		candidateInfo.CandidateName = dataPtr.GetField("candidateName")
		candidateInfo.CandidateUrl = dataPtr.GetField("candidateUrl")
	}

	return candidateInfo, nil
}

// QueryVoterAvailableVotes
func QueryVoterAvailableVotes(strSQL string) (float64, error) {
	var voterAvailableVotes = int64(0)
	var result = float64(0)
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryCandidateInfo error:[%v]\n", err)
		return result, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryCandidateInfo dataPtr is nil ")
		return result, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		frozen := dataPtr.GetField("frozen")
		if frozen != "" {
			accountFrozenBalance := make([]*entity.BalanceInfoDB, 0)
			if err := json.Unmarshal([]byte(frozen), &accountFrozenBalance); err != nil {
				log.Errorf("Unmarshal data failed:[%v]-%v", err, frozen)
			}
			for _, v := range accountFrozenBalance {
				voterAvailableVotes += v.Amount
			}
		}
	}
	result = float64(voterAvailableVotes)/1000000
	return result, nil
}



