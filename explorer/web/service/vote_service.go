package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryVotes 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryVotes(req *entity.Votes) (*entity.VotesResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
	SELECT trx_hash,block_id,voter_address,candidate_address,vote_num,wit.create_time,
	       acc.frozen,acc.account_name,wlwit.url,outvoter.outVotes
	FROM tron.contract_vote_witness wit
	left join tron.tron_account acc on acc.address=wit.candidate_address
	left join tron.wlcy_witness_create_info wlwit on wlwit.address=wit.candidate_address
	left join (
		select address,sum(vote) as outVotes from tron.account_vote_result  group by address
	) outvoter on outvoter.address=wit.voter_address
	where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Candidate != "" {
		filterSQL = fmt.Sprintf(" and wit.candidate_address='%v'", req.Candidate)
	}
	if req.Voter != "" {
		filterSQL = fmt.Sprintf(" and wit.voter_address='%v'", req.Voter)
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortTemp = fmt.Sprintf("%v ,", sortTemp)
			}
			sortTemp = fmt.Sprintf("%v wit.create_time", sortTemp)
			if strings.Index(v, "-") == 0 {
				sortTemp = fmt.Sprintf("%v desc", sortTemp)
			}
			mutiFilter = true
		}

		if strings.Index(v, "number") > 0 {
			if mutiFilter {
				sortTemp = fmt.Sprintf("%v ,", sortTemp)
			}
			sortTemp = fmt.Sprintf("%v wit.block_id", sortTemp)
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
	if req.Limit != "" && req.Start != "" {
		pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	}
	return module.QueryVotesRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryVoteLive 实时投票数据
func QueryVoteLive() (*entity.VoteLiveInfo, error) {
	strSQL := fmt.Sprintf(`
	SELECT trx_hash,block_id,voter_address,candidate_address,vote_num,wit.create_time,
	       acc.account_name,wlwit.url,getvoter.getVotes
	FROM tron.contract_vote_witness wit
	left join tron.tron_account acc on acc.address=wit.candidate_address
	left join tron.wlcy_witness_create_info wlwit on wlwit.address=wit.candidate_address
	left join (
		select address,sum(vote) as getVotes from tron.account_vote_result  group by address
	) getvoter on getvoter.address=wit.candidate_address
	where 1=1 `)

	return module.QueryVoteLiveRealize(strSQL)
}

//QueryVoteCurrentCycle 上轮投票数据
func QueryVoteCurrentCycle() (*entity.VoteCurrentCycleResp, error) {
	/*var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,amount,
		token_name,trx_hash,
		contract_type,confirmed,create_time
		from tron.contract_token_transfer
			where 1=1 `)
	*/
	return nil, nil
}

//QueryVoteNextCycle 本轮投票剩余时长
// 使用旧版scala逻辑
func QueryVoteNextCycle() (*entity.VoteNextCycleResp, error) {
	var nextCycle = &entity.VoteNextCycleResp{}
	nextCycle.NextCycle = 0
	var nextMaintenanceTime, currentTime int64
	client := grpcclient.NewWallet(fmt.Sprintf("%s:50051", utils.GetRandFullNodeAddr()))
	err := client.Connect()
	if nil != err {
		log.Error(err)
		return nextCycle, err
	}

	log.Debugf(client.GetState(), client.Target())
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
