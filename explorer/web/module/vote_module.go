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
	SELECT sum(vote_count) as totalVotes FROM witness`)
	log.Sql(strSQL)
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

//QueryRealTimeTotalVotes
func QueryRealTimeTotalVotes(strSQL string) int64 {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryRealTimeTotalVotes error :[%v]\n", err)
		return 0
	}
	if dataPtr == nil {
		log.Errorf("QueryRealTimeTotalVotes dataPtr is nil ")
		return 0
	}
	var votes = int64(0)
	for dataPtr.NextT() {
		votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalVotes"))
	}

	return votes
}

//QueryAccountVoteResultRealize
func QueryAccountVoteResultRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AccountVoteResultRes, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
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
	log.Sql(strSQL)
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
	log.Sql(strSQL)
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
	result = float64(voterAvailableVotes) / 1000000
	return result, nil
}

// QueryVoteWitness
func QueryVoteWitness(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.VoteWitnessResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryVoteWitness error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryVoteWitness dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	voteWitnessResp := &entity.VoteWitnessResp{}
	voteWitnessList := make([]*entity.VoteWitness, 0)

	for dataPtr.NextT() {
		var voteWitness = &entity.VoteWitness{}
		voteWitness.Address = dataPtr.GetField("address")
		voteWitness.LastCycleVotes = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_count"))
		voteWitness.Name = dataPtr.GetField("account_name")
		voteWitness.URL = dataPtr.GetField("github_link")
		voteWitness.RealTimeVotes = mysql.ConvertDBValueToInt64(dataPtr.GetField("realTimeVotes"))

		voteWitnessList = append(voteWitnessList, voteWitness)
	}

	var total = int64(len(voteWitnessList))
	total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}

	voteWitnessResp.Total = total
	voteWitnessResp.Data = voteWitnessList

	return voteWitnessResp, nil
}

// QueryRealTimeVoteWitnessTotal
func QueryRealTimeVoteWitnessTotal(strSQL string) int64 {
	var realTimeVotes = int64(0)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryRealTimeVoteWitnessTotal error :[%v]\n", err)
		return 0
	}
	if dataPtr == nil {
		log.Errorf("QueryRealTimeVoteWitnessTotal dataPtr is nil ")
		return 0
	}
	for dataPtr.NextT() {
		realTimeVotes = mysql.ConvertDBValueToInt64(dataPtr.GetField("realTimeVotes"))
	}
	return realTimeVotes
}

// QueryVoteWitnessRanking
func QueryVoteWitnessRanking(strSQL string) ([]*entity.VoteWitnessRanking, error) {
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryVoteWitnessRanking error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryVoteWitnessRanking dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	voteWitnessRankingList := make([]*entity.VoteWitnessRanking, 0)

	for dataPtr.NextT() {
		var voteWitnessRanking = &entity.VoteWitnessRanking{}
		voteWitnessRanking.Address = dataPtr.GetField("address")
		voteWitnessRankingList = append(voteWitnessRankingList, voteWitnessRanking)
	}

	return voteWitnessRankingList, nil
}

// QueryAddressVoter
func QueryAddressVoter(strSQL string) ([]*entity.AddressVoteInfo, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryAddressVoter error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAddressVoter dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	AddressVoteInfoList := make([]*entity.AddressVoteInfo, 0)

	for dataPtr.NextT() {
		var addressVoteInfo = &entity.AddressVoteInfo{}
		addressVoteInfo.VoterAddress = dataPtr.GetField("address")
		addressVoteInfo.CandidateAddress = dataPtr.GetField("to_address")
		addressVoteInfo.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote"))

		AddressVoteInfoList = append(AddressVoteInfoList, addressVoteInfo)
	}
	return AddressVoteInfoList, nil
}
