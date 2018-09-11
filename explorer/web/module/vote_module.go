package module

import (
	"encoding/json"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryVotesRealize 操作数据库
func QueryVotesRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.VotesResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryVotesRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryVotesRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	votesResp := &entity.VotesResp{}
	accountFrozenBalance := make([]*entity.BalanceInfoDB, 0)
	voteInfos := make([]*entity.VotesInfo, 0)
	var totalFrozen = int64(0)

	//填充数据
	for dataPtr.NextT() {
		var vote = &entity.VotesInfo{}
		vote.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		vote.Transaction = dataPtr.GetField("trx_hash")
		vote.VoterAddress = dataPtr.GetField("voter_address")
		vote.CreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		vote.CandidateAddress = dataPtr.GetField("candidate_address")
		vote.CandidateName = dataPtr.GetField("account_name")
		vote.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_num"))
		vote.CandidateURL = dataPtr.GetField("url")
		outerVote := mysql.ConvertDBValueToInt64(dataPtr.GetField("outVotes"))
		frozen := dataPtr.GetField("frozen")
		if err := json.Unmarshal([]byte(frozen), accountFrozenBalance); err != nil {
			log.Errorf("Unmarshal data failed:[%v]", err)
		}
		for _, blanceFrozen := range accountFrozenBalance {
			totalFrozen += blanceFrozen.Amount
		}
		vote.VoterAvailableVotes = totalFrozen - outerVote

		voteInfos = append(voteInfos, vote)
	}

	//查询该语句所查到的数据集合
	var total = int64(len(voteInfos))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	votesResp.Total = total
	votesResp.Data = voteInfos

	return votesResp, nil

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
		vote.Address = dataPtr.GetField("voter_address")
		vote.Name = dataPtr.GetField("account_name")
		vote.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_num"))
		vote.URL = dataPtr.GetField("url")
		outerVote := mysql.ConvertDBValueToInt64(dataPtr.GetField("getVotes"))
		log.Debugf("sum vote from contact table is:[%v]", outerVote)
		votesMap[vote.Address] = vote
	}
	votesResp.Data = votesMap

	return votesResp, nil

}
