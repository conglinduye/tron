package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

// QueryVoteWitnessBuffer
func QueryVoteWitnessBuffer() (*entity.VoteWitnessResp, error) {
	voteBuffer := buffer.GetVoteBuffer()
	voteWitnessResp := voteBuffer.GetVoteWitness()
	if voteWitnessResp == nil {
		voteWitnessResp = &entity.VoteWitnessResp{}
		voteWitnessResp.Total = 0
		voteWitnessResp.TotalVotes = 0
		voteWitnessList := make([]*entity.VoteWitness, 0)
		voteWitnessResp.Data = voteWitnessList
		voteWitness := &entity.VoteWitness{}
		voteWitnessResp.FastestRise = voteWitness
		return voteWitnessResp, nil
	}
	return voteWitnessResp, nil

}

//QueryVoteNextCycleBuffer 本轮投票剩余时长
func QueryVoteNextCycleBuffer() (*entity.VoteNextCycleResp, error) {
	var nextCycle = &entity.VoteNextCycleResp{}
	nextCycle.NextCycle = 0
	currentTime := buffer.GetBlockBuffer().GetMaxBlockTimestamp()
	nextMaintenanceTime := buffer.GetVoteBuffer().GetNextMaintenanceTime()
	nextCycle.NextCycle = nextMaintenanceTime - currentTime
	if currentTime == 0 || nextMaintenanceTime == 0 || nextCycle.NextCycle < 0 {
		return QueryVoteNextCycle()
	}
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
		votesResp.Data = make([]*entity.VotesInfo, 0)
		return votesResp, nil
	}

	voteInfoList := make([]*entity.VotesInfo, 0)
	for _, v := range accountVoteResultRes.Data {
		votesInfo := &entity.VotesInfo{}
		votesInfo.VoterAddress = v.Address
		votesInfo.CandidateAddress = v.ToAddress
		votesInfo.Votes = v.Vote
		voteInfoList = append(voteInfoList, votesInfo)
	}

	votesResp.Total = accountVoteResultRes.Total
	votesResp.Data = voteInfoList

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
			log.Errorf("QueryVotesSubHandle queryCandidateInfo strSQL:%v, err:[%v]", strSQLOne, err)
		} else {
			votesInfo.CandidateName = candidateInfo.CandidateName
			votesInfo.CandidateURL = candidateInfo.CandidateUrl
		}

		strSQLTwo := fmt.Sprintf(`select frozen from tron_account where address = '%v'`, votesInfo.VoterAddress)

		voterAvailableVotes, err := module.QueryVoterAvailableVotes(strSQLTwo)
		if err != nil {
			log.Errorf("QueryVotesSubHandle queryVoterAvailableVotes strSQL:%v, err:[%v]", strSQLTwo, err)
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

// QueryVoteWitnessDetail
func QueryVoteWitnessDetail(address string) (*entity.VoteWitnessDetail, error) {
	voteWitnessDetail := &entity.VoteWitnessDetail{}
	voteWitness := &entity.VoteWitness{}
	voteWitnessResp, _ := QueryVoteWitnessBuffer()
	voteWitnessList := voteWitnessResp.Data
	if len(voteWitnessList) == 0 {
		voteWitnessDetail.Success = false
		voteWitnessDetail.Data = voteWitness
		return voteWitnessDetail, nil
	}
	voteWitnessDetail.Success = true

	for _, temp := range voteWitnessList {
		if address == temp.Address {
			voteWitness = temp
			break
		}
	}

	voteWitnessDetail.Data = voteWitness

	return voteWitnessDetail, nil
}

// QueryVoterAvailableVotes
func QueryVoterAvailableVotes(address string) *entity.AddressVotes {
	addressVotes := &entity.AddressVotes{}
	votes := make(map[string]float64, 0)
	addressVotes.Votes = votes
	strSQL := fmt.Sprintf(`select frozen from tron_account where address = '%v'`, address)
	voterAvailableVotes, err := module.QueryVoterAvailableVotes(strSQL)
	if err != nil {
		log.Errorf("QueryVoterAvailableVotes strSQL:%v, err:[%v]", strSQL, err)
	} else {
		votes[address] = voterAvailableVotes
	}
	return addressVotes
}
