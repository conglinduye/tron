package module

import (
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryWitnessRealize 操作数据库
func QueryWitnessRealize(strSQL string) ([]*entity.WitnessInfo, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryWitnessRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryWitnessRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	witnessInfos := make([]*entity.WitnessInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var witness = &entity.WitnessInfo{}
		witness.Votes = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_count"))
		witness.Address = dataPtr.GetField("address")
		witness.ProducedTotal = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_produced"))
		witness.URL = dataPtr.GetField("url")
		witness.Name = dataPtr.GetField("account_name")
		witness.LatestBlockNumber = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_block_num"))
		witness.MissedTotal = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_missed"))
		witness.LatestSlotNumber = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_slot_num"))
		isJob := dataPtr.GetField("is_job")
		if isJob == "1" {
			witness.Producer = true
		}
		witnessInfos = append(witnessInfos, witness)
	}

	return witnessInfos, nil

}

//QueryTotalBlocks 查询总block
func QueryTotalBlocks(curTime int64) (int64, error) {
	var totalBlock = int64(0)
	strSQL := fmt.Sprintf(`
    select ifnull(count(block_id),0) as totalBlock
    from tron.blocks blk
	where 1=1 and blk.create_time>=%v `, curTime)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalBlocks error :[%v]\n", err)
		return totalBlock, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalBlocks dataPtr is nil ")
		return totalBlock, util.NewErrorMsg(util.Error_common_internal_error)
	}

	//填充数据
	for dataPtr.NextT() {
		totalBlock = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBlock"))
	}
	return totalBlock, nil

}

//QueryWitnessStatisticRealize 操作数据库
func QueryWitnessStatisticRealize(strSQL string, totalBlocks int64) ([]*entity.WitnessStatisticInfo, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryWitnessStatisticRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryWitnessStatisticRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	witnessInfos := make([]*entity.WitnessStatisticInfo, 0)
	//填充数据
	for dataPtr.NextT() {
		var witness = &entity.WitnessStatisticInfo{}
		witness.Address = dataPtr.GetField("address")
		witness.Name = dataPtr.GetField("account_name")
		witness.URL = dataPtr.GetField("url")
		witness.BlockProduced = mysql.ConvertDBValueToInt64(dataPtr.GetField("blockproduce"))
		witness.Total = totalBlocks
		witness.Percentage = mysql.ConvertDBValueToFloat64(dataPtr.GetField("blockRate"))
		witnessInfos = append(witnessInfos, witness)
	}
	return witnessInfos, nil

}
