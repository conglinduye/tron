package service

import (
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
	"sort"
)

//QueryWitnessBuffer 从QueryWitnessBuffer中获取witness信息
func QueryWitnessBuffer() ([]*entity.WitnessInfo, error) {

	witnessBuffer := buffer.GetWitnessBuffer()

	witnessList := witnessBuffer.GetWitness()

	return witnessList, nil
}

//QueryWitness
func QueryWitness() ([]*entity.WitnessInfo, error) {
	strSQL := fmt.Sprintf(`
			select witt.address,witt.vote_count,witt.public_key,witt.url,
			witt.total_produced,witt.total_missed,acc.account_name,
			witt.latest_block_num,witt.latest_slot_num,witt.is_job
			from witness witt
			left join tron_account acc on acc.address=witt.address
			where 1=1 `)

	witnessInfoList, err := module.QueryWitnessRealize(strSQL)
	if nil != err {
		log.Errorf("load witness from db failed:%v\n", err)
		return make([]*entity.WitnessInfo, 0), err
	}

	totalVotes := module.QueryTotalVotes()
	for index := range witnessInfoList {
		witnessInfo := witnessInfoList[index]
		if witnessInfo.ProducedTotal != 0 {
			witnessInfo.ProducePercentage = float64(witnessInfo.ProducedTotal-witnessInfo.MissedTotal)/float64(witnessInfo.ProducedTotal) * 100
		} else {
			witnessInfo.ProducePercentage = 0
		}
		if totalVotes != 0 {
			witnessInfo.VotesPercentage = float64(witnessInfo.Votes)/float64(totalVotes) * 100
		} else {
			witnessInfo.VotesPercentage = 0
		}
	}

	addrMap := make(map[string]*entity.WitnessInfo, len(witnessInfoList))
	sortList := make([]*entity.WitnessInfo, 0, len(witnessInfoList))
	for _, witness := range witnessInfoList {
		addrMap[witness.Address] = witness
		sortList = append(sortList, witness)
	}
	// votes 大的排在前面
	sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].Votes > sortList[j].Votes })

	return sortList, nil
}

//QueryWitnessStatisticBuffer  从buffer获取
func QueryWitnessStatisticBuffer() ([]*entity.WitnessStatisticInfo, error) {
	witnessBuffer := buffer.GetWitnessBuffer()
	witnessInfo := witnessBuffer.GetWitnessStatistic()
	return witnessInfo, nil
}

//获取当前轮开始时间戳
func getMaintenanceTimeStamp() (int64, error) {

	client := grpcclient.GetRandomWallet()

	nextMaintenanceTime, err := client.GetNextMaintenanceTime()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	curMaintenanceTime := nextMaintenanceTime - 6*60*60*1000 //6小时
	return curMaintenanceTime, nil
}
