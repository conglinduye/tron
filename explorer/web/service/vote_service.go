package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"github.com/wlcy/tron/explorer/lib/util"
	"sort"
	"encoding/json"
	"github.com/wlcy/tron/explorer/lib/config"
	"gopkg.in/redis.v4"
)


const latelyCycleVoteWitnessRankingKey ="lately.cycle.vote.witness.ranking.key"
const latelyDayVoteWitnessRankingKey ="lately.day.vote.witness.ranking.key"

//QueryVoteLiveBuffer 从buffer中获取实时投票数据
func QueryVoteLiveBuffer() (*entity.VoteLiveInfo, error) {
	var voteLive = &entity.VoteLiveInfo{}
	voteBuffer := buffer.GetVoteBuffer()
	votes, _ := voteBuffer.GetVoteLive()
	voteLive.Data = votes
	return voteLive, nil

}

//QueryVoteCurrentCycleBuffer 从buffer中获取上轮投票数据
func QueryVoteCurrentCycleBuffer() (*entity.VoteCurrentCycleResp, error) {
	voteBuffer := buffer.GetVoteBuffer()
	return voteBuffer.GetVoteCurrentCycle(), nil
}

//QueryVoteNextCycleBuffer 本轮投票剩余时长
func QueryVoteNextCycleBuffer() (*entity.VoteNextCycleResp, error) {
	var nextCycle = &entity.VoteNextCycleResp{}
	nextCycle.NextCycle = 0
	currentTime := buffer.GetBlockBuffer().GetMaxBlockTimestamp()
	nextMaintenanceTime := buffer.GetVoteBuffer().GetNextMaintenanceTime()
	if currentTime == 0 || nextMaintenanceTime == 0 {
		return QueryVoteNextCycle()
	}
	nextCycle.NextCycle = nextMaintenanceTime - currentTime
	return nextCycle, nil
}

//QueryVoteNextCycle 本轮投票剩余时长
// 使用旧版scala逻辑
func QueryVoteNextCycle() (*entity.VoteNextCycleResp, error) {
	var nextCycle = &entity.VoteNextCycleResp{}
	nextCycle.NextCycle = 0
	var nextMaintenanceTime, currentTime int64

	client := grpcclient.GetRandomWallet()

	block, err := client.GetNowBlock()
	if err != nil {
		log.Error(err)
		return nextCycle, err
	}
	if block != nil && block.BlockHeader != nil && block.BlockHeader.RawData != nil {
		currentTime = block.BlockHeader.RawData.Timestamp
	}
	nextMaintenanceTime, err = client.GetNextMaintenanceTime()
	if err != nil {
		log.Error(err)
		return nextCycle, err
	}
	nextCycle.NextCycle = nextMaintenanceTime - currentTime
	return nextCycle, nil
}

// QueryVotes
func QueryVotes(req *entity.Votes) (*entity.VotesResp, error) {
	votesResp := &entity.VotesResp{}
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select address, to_address, vote from account_vote_result where 1=1 `)

	if req.Voter != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", req.Voter)
	}
	if req.Candidate != "" {
		filterSQL = fmt.Sprintf(" and to_address='%v'", req.Candidate)
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "votes") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v vote", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}
	}
	if sortSQL != "" {
		if strings.Index(sortSQL, ",") == 0 {
			sortSQL = sortSQL[1:]
		}
		sortSQL = fmt.Sprintf("order by %v", sortSQL)
	}
	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	accountVoteResultRes, err := module.QueryAccountVoteResultRealize(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		log.Errorf("QueryVotes list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	if len(accountVoteResultRes.Data) == 0 {
		votesResp.Data = make([]*entity.VotesInfo , 0)
		return votesResp, nil
	}

	voteInfos := make([]*entity.VotesInfo, 0)
	for _, v := range accountVoteResultRes.Data {
		votesInfo := &entity.VotesInfo{}
		votesInfo.VoterAddress = v.Address
		votesInfo.CandidateAddress = v.ToAddress
		votesInfo.Votes = v.Vote
		voteInfos = append(voteInfos, votesInfo)
	}

	votesResp.Total = accountVoteResultRes.Total
	votesResp.Data = voteInfos

	queryVotesSubHandle(votesResp)

	totalVotes := QueryRealTimeTotalVotes(req)
	votesResp.TotalVotes = totalVotes

	return votesResp, nil
}

// QueryVotesSubHandle
func queryVotesSubHandle(votesResp *entity.VotesResp) {
	votesInfos := votesResp.Data
	for _, votesInfo := range votesInfos {

		strSQLOne := fmt.Sprintf(`
			select acc.address as candidateAddress, acc.account_name as candidateName, wlwit.url as candidateUrl
			from tron_account acc 
			left join wlcy_witness_create_info wlwit on wlwit.address=acc.address 
			where acc.address = '%v'`, votesInfo.CandidateAddress)

		candidateInfo, err := module.QueryCandidateInfo(strSQLOne)
		if err != nil {
			log.Errorf("QueryVotesSubHandle queryCandidateInfo strSQL:%v, err:[%v]",strSQLOne,  err)
		} else {
			votesInfo.CandidateName = candidateInfo.CandidateName
			votesInfo.CandidateURL = candidateInfo.CandidateUrl
		}

		strSQLTwo := fmt.Sprintf(`select frozen from tron_account where address = '%v'`, votesInfo.VoterAddress)

		voterAvailableVotes, err := module.QueryVoterAvailableVotes(strSQLTwo)
		if err != nil {
			log.Errorf("QueryVotesSubHandle queryVoterAvailableVotes strSQL:%v, err:[%v]",strSQLTwo, err)
		} else {
			votesInfo.VoterAvailableVotes = voterAvailableVotes
		}
	}
}


func QueryRealTimeTotalVotes(req *entity.Votes) int64 {
	filterSQL := ""
	strSQL := fmt.Sprintf(`
			select sum(vote) as totalVotes from account_vote_result  where 1=1 `)
	if req.Voter != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", req.Voter)
	}
	if req.Candidate != "" {
		filterSQL = fmt.Sprintf(" and to_address='%v'", req.Candidate)
	}

	totalVotes := module.QueryRealTimeTotalVotes(strSQL + filterSQL)
	return totalVotes
}

// QueryVoteWitness
func QueryVoteWitness(req *entity.VoteWitnessReq) (*entity.VoteWitnessResp, error) {
	var filterSQL, sortSQL, pageSQL string
	strSQL := fmt.Sprintf(`
		select witt.address, witt.vote_count, srac.github_link, acc.account_name,votes.realTimeVotes
		from witness witt
		left join tron_account acc on acc.address=witt.address
		left join wlcy_sr_account srac on witt.address=srac.address
		left join (
			select to_address,sum(vote) as realTimeVotes from account_vote_result  group by to_address 
		) votes on votes.to_address=witt.address
		where 1=1 `)

	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and witt.address='%v'", req.Address)
	}
	sortSQL = "votes.realTimeVotes desc"

	//pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	voteWitnessResp, err := module.QueryVoteWitness(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		log.Errorf("QueryVoteWitness strSQL:%v, err:[%v]",strSQL, err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	totalVotes := module.QueryTotalVotes()
	voteWitnessResp.TotalVotes = totalVotes

	voteWitnessList:= voteWitnessResp.Data
	for index, voteWitness := range voteWitnessList {
		voteWitness.ChangeVotes = voteWitness.RealTimeVotes - voteWitness.LastCycleVotes
		if voteWitness.URL != "" {
			voteWitness.HasPage = true
		}
		if totalVotes != 0 {
			voteWitness.VotesPercentage = float64(voteWitness.LastCycleVotes) / float64(totalVotes) * 100
		}
		voteWitness.RealTimeRanking = int32(index + 1)
	}

	//sort.SliceStable(voteWitnessList, func(i, j int) bool { return voteWitnessList[i].RealTimeVotes > voteWitnessList[j].RealTimeVotes })

	/*for index := range voteWitnessList {
		voteWitnessList[index].RealTimeRanking = int32(index + 1)
	}*/

	// getVoteWitnessRankingChange
	getVoteWitnessRankingChange(voteWitnessList)

	sortList := make([]*entity.VoteWitness, 0, len(voteWitnessList))
	for index := range voteWitnessList {
		voteWitness := new(entity.VoteWitness)
		*voteWitness = *voteWitnessList[index]
		sortList = append(sortList, voteWitness)
	}

	if len(sortList) > 0 {
		sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].ChangeCycle > sortList[j].ChangeCycle })
		voteWitnessResp.FastestRise = sortList[0]
	}

	return voteWitnessResp, nil
}

// queryRealTimeVoteWitness
func queryRealTimeVoteWitnessTotal(toAddress string) int64 {
	strSQL := fmt.Sprintf(`select sum(vote) as realTimeVotes from account_vote_result where to_address = '%v' `, toAddress)
	realTimeVotes := module.QueryRealTimeVoteWitnessTotal(strSQL)
	return realTimeVotes
}

// getVoteWitnessRankingChange
func getVoteWitnessRankingChange(voteWitnessList []*entity.VoteWitness) {
	var latelyDayVoteWitnessRankingValue, latelyCycleVoteWitnessRankingValue string
	var err error
	latelyDayVoteWitnessRankingValue, err = config.RedisCli.Get(latelyDayVoteWitnessRankingKey).Result()
	if err == redis.Nil {
		SyncVoteWitnessRanking()
		latelyDayVoteWitnessRankingValue, _ = config.RedisCli.Get(latelyDayVoteWitnessRankingKey).Result()
	} else if err != nil {
		log.Errorf("getVoteWitnessRankingChange redis get latelyDayVoteWitnessRankingValue error :[%v]\n", err)
		return
	}

	latelyCycleVoteWitnessRankingValue, _ = config.RedisCli.Get(latelyCycleVoteWitnessRankingKey).Result()

	latelyDayVoteWitnessRankingList := make([]*entity.VoteWitnessRanking, 0)
	err = json.Unmarshal([]byte(latelyDayVoteWitnessRankingValue), &latelyDayVoteWitnessRankingList)
	if err != nil {
		log.Errorf("getVoteWitnessRankingChange json.Unmarshal latelyDayVoteWitnessRankingList error :[%v]\n", err)
		return
	}

	latelyCycleVoteWitnessRankingList := make([]*entity.VoteWitnessRanking, 0)
	err = json.Unmarshal([]byte(latelyCycleVoteWitnessRankingValue), &latelyCycleVoteWitnessRankingList)
	if err != nil {
		log.Errorf("getVoteWitnessRankingChange json.Unmarshal latelyCycleVoteWitnessRankingList error :[%v]\n", err)
		return
	}

	for index := range voteWitnessList {
		voteWitness := voteWitnessList[index]
		changeCycle :=  getChangeRanking(voteWitness.Address, voteWitness.RealTimeRanking, latelyCycleVoteWitnessRankingList)
		voteWitness.ChangeCycle = changeCycle
		changeDay := getChangeRanking(voteWitness.Address, voteWitness.RealTimeRanking, latelyDayVoteWitnessRankingList)
		voteWitness.ChangeDay = changeDay
	}

}

// getChangeRanking
func  getChangeRanking(address string, currentRanking int32, voteWitnessRankingList []*entity.VoteWitnessRanking) int32 {
	for index := range voteWitnessRankingList {
		voteWitnessRanking := voteWitnessRankingList[index]
		if voteWitnessRanking.Address == address {
			return voteWitnessRanking.Ranking - currentRanking
		}
	}
	return 0
}


//syncLatelyCycleVoteWitnessRanking
func syncLatelyCycleVoteWitnessRanking() {
	strSQL := fmt.Sprintf(`select address, vote_count from witness order by vote_count desc `)

	voteWitnessRankingList, err := module.QueryVoteWitnessRanking(strSQL)

	if err != nil {
		log.Errorf("syncLatelyCycleVoteWitnessRanking strSQL:%v, err:[%v]",strSQL,  err)
		return
	}

	for index := range voteWitnessRankingList {
		voteWitnessRankingList[index].Ranking = int32(index+1)
	}

	value, err := json.Marshal(voteWitnessRankingList)
	if err != nil {
		log.Errorf("syncLatelyCycleVoteWitnessRanking json.Marshal err:[%v]", err)
		return
	}

	err = config.RedisCli.Set(latelyCycleVoteWitnessRankingKey, string(value), 0).Err()
	if err != nil {
		log.Errorf("syncLatelyCycleVoteWitnessRanking set lately_cycle_vote_witness_ranking_key err:[%v]", err)
	}
}

// SyncVoteWitnessRanking
func SyncVoteWitnessRanking() {
	var latelyDayVoteWitnessRankingValue, latelyCycleVoteWitnessRankingValue string
	var err error
	latelyDayVoteWitnessRankingValue, err = config.RedisCli.Get(latelyDayVoteWitnessRankingKey).Result()
	if err == redis.Nil {
		latelyCycleVoteWitnessRankingValue, err = config.RedisCli.Get(latelyCycleVoteWitnessRankingKey).Result()
		if err == redis.Nil {
			syncLatelyCycleVoteWitnessRanking()
			latelyCycleVoteWitnessRankingValue, _ = config.RedisCli.Get(latelyCycleVoteWitnessRankingKey).Result()
		} else if err != nil  {
			log.Errorf("syncVoteWitnessRanking redis get latelyCycleVoteWitnessRankingValue error :[%v]\n", err)
			return
		}

		latelyDayVoteWitnessRankingValue = latelyCycleVoteWitnessRankingValue

	} else if err != nil {
		log.Errorf("syncVoteWitnessRanking redis get latelyDayVoteWitnessRankingValue error :[%v]\n", err)
		return
	}

	latelyCycleVoteWitnessRankingValue, err = config.RedisCli.Get(latelyCycleVoteWitnessRankingKey).Result()
	if err == redis.Nil {
		syncLatelyCycleVoteWitnessRanking()
		latelyCycleVoteWitnessRankingValue, _ = config.RedisCli.Get(latelyCycleVoteWitnessRankingKey).Result()
	} else if err != nil {
		log.Errorf("syncVoteWitnessRanking redis get latelyCycleVoteWitnessRankingValue error :[%v]\n", err)
		return
	}

	latelyDayVoteWitnessRankingValue = latelyCycleVoteWitnessRankingValue

	err = config.RedisCli.Set(latelyDayVoteWitnessRankingKey, string(latelyDayVoteWitnessRankingValue), 0).Err()
	if err != nil {
		log.Errorf("syncVoteWitnessRanking set lately_day_vote_witness_ranking_key err:[%v]", err)
	}
	syncLatelyCycleVoteWitnessRanking()

}

// QueryVoteWitnessDetail
func QueryVoteWitnessDetail(address string) (*entity.VoteWitnessDetail, error) {
	voteWitnessDetail := &entity.VoteWitnessDetail{}
	voteWitness := &entity.VoteWitness{}
	req := &entity.VoteWitnessReq{}
	voteWitnessResp, err := QueryVoteWitness(req)
	if err != nil {
		log.Errorf("QueryVoteWitnessDetail error :[%v]\n", err)
		voteWitnessDetail.Success = false
		voteWitnessDetail.Data = voteWitness
		return voteWitnessDetail, err
	}

	voteWitnessDetail.Success = true
	if len(voteWitnessResp.Data) != 0 {
		voteWitnessList := voteWitnessResp.Data
		for _, temp := range voteWitnessList {
			if address == temp.Address {
				voteWitness = temp
				break
			}
		}
		voteWitnessDetail.Data = voteWitness
	} else {
		voteWitnessDetail.Success = false
	}
	return voteWitnessDetail, nil

}
