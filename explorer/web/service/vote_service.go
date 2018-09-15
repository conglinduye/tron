package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryVotes 条件查询  	//?sort=-number&limit=1&count=true&number=2135998   TODO: cache
func QueryVotes(req *entity.Votes) (*entity.VotesResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false
	//默认查询得票列表
	reportSQL := fmt.Sprint(`
	left join (
		select to_address,sum(vote) as votes from tron.account_vote_result 
		 group by to_address
	) outvoter on outvoter.to_address=acc.address`)

	//按照Voter过滤，获取该Voter投给别的人列表
	if req.Voter != "" {
		reportSQL = fmt.Sprintf(`
	left join (
		select to_address,sum(vote) as votes from tron.account_vote_result 
		where 1=1 and address='%v'
		 group by to_address
	) outvoter on outvoter.to_address=acc.address`, req.Voter)
	}
	//按照Candidate过滤，获取谁投给Candidate的列表
	if req.Candidate != "" {
		reportSQL = fmt.Sprintf(`
	left join (
		select address,sum(vote) as votes from tron.account_vote_result  
		where 1=1 and to_address='%v'
		group by address
	) outvoter on outvoter.address=acc.address`, req.Candidate)
	}
	strSQL := fmt.Sprintf(`
	SELECT acc.address as voteraddress,outvoter.votes,
	       acc.frozen,acc.account_name,wlwit.url
	FROM tron.tron_account acc 
	left join tron.wlcy_witness_create_info wlwit on wlwit.address=acc.address
	%v
     where 1=1 and outvoter.votes>0 `, reportSQL)

	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "votes") > 0 {
			if mutiFilter {
				sortTemp = fmt.Sprintf("%v ,", sortTemp)
			}
			sortTemp = fmt.Sprintf("%v outvoter.votes", sortTemp)
			if strings.Index(v, "-") == 0 {
				sortTemp = fmt.Sprintf("%v desc", sortTemp)
			}
			mutiFilter = true
		}
	}
	if sortTemp != "" {
		if strings.Index(sortTemp, ",") == 0 {
			sortTemp = sortTemp[1:]
		}
		sortTemp = fmt.Sprintf("order by %v", sortTemp)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryVotesRealize(strSQL, filterSQL, sortSQL, pageSQL, req)
}

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
	currentTime := buffer.GetBlockBufferInstance().GetMaxBlockTimestamp()
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
