package service

import (
	"fmt"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryWitnessBuffer 从QueryWitnessBuffer中获取witness信息
func QueryWitnessBuffer() ([]*entity.WitnessInfo, error) {
	var witnessList = make([]*entity.WitnessInfo, 0)
	var err error
	witnessBuffer := buffer.GetWitnessBuffer()
	if witnessBuffer == nil {
		witnessList, err = QueryWitness()
	} else {
		witnessList = witnessBuffer.GetWitness()
		if witnessList == nil {
			witnessList, err = QueryWitness()
		}
	}
	return witnessList, err
}

//QueryWitness ...
func QueryWitness() ([]*entity.WitnessInfo, error) {

	strSQL := fmt.Sprintf(`
			select witt.address,witt.vote_count,witt.public_key,witt.url,
			witt.total_produced,witt.total_missed,acc.account_name,
			witt.latest_block_num,witt.latest_slot_num,witt.is_job
			from tron.witness witt
			left join tron.tron_account acc on acc.address=witt.address
			where 1=1 `)

	return module.QueryWitnessRealize(strSQL)
}

//QueryWitnessStatisticBuffer  从buffer获取
func QueryWitnessStatisticBuffer() ([]*entity.WitnessStatisticInfo, error) {
	var err error
	witnessBuffer := buffer.GetWitnessBuffer()
	witnessInfo := witnessBuffer.GetWitnessStatistic()
	if witnessInfo == nil {
		log.Debug("no buffer for witness StatisticBuffer, get them form db")
		witnessInfo, err = QueryWitnessStatistic()
	}
	return witnessInfo, err
}

//QueryWitnessStatistic  ...
func QueryWitnessStatistic() ([]*entity.WitnessStatisticInfo, error) {
	var blocks int64
	curMaintenanceTime, err := getMaintenanceTimeStamp()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	totalBlocks, err := module.QueryTotalBlocks(curMaintenanceTime)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if totalBlocks == 0 {
		log.Errorf("total blocks is 0, replace it with large number")
		blocks = 10000000
	} else {
		blocks = totalBlocks
	}
	strSQL := fmt.Sprintf(`
	select acc.address, acc.account_name,witt.url
		   ,ifnull(blocks.blockproduce,0) as blockproduce , 
		   ifnull(blocks.blockproduce,0)/%v as blockRate
    from  tron.tron_account acc
    left join tron.witness witt on witt.address=acc.address 
    left join (
	    select witness_address,count(block_id) as blockproduce
        from tron.blocks blk
        where 1=1 and blk.create_time>%v 
        group by witness_address
    ) blocks on blocks.witness_address=acc.address
    where 1=1 and acc.is_witness=1`, blocks, curMaintenanceTime)

	return module.QueryWitnessStatisticRealize(strSQL, totalBlocks)
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
