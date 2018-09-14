package main

import (
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryBlocksRealize 操作数据库
func QueryBlocksRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.BlocksResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryBlocks error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryBlocks dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	blocksResp := &entity.BlocksResp{}
	blockInfos := make([]*entity.BlockInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var block = &entity.BlockInfo{}
		block.Number = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		block.Hash = dataPtr.GetField("block_hash")
		block.Size = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_size"))
		block.CreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		block.TxTrieRoot = dataPtr.GetField("tx_trie_hash")
		block.ParentHash = dataPtr.GetField("parent_hash")
		block.WitnessAddress = dataPtr.GetField("witness_address")
		block.WitnessID = 0
		block.NrOfTrx = mysql.ConvertDBValueToInt64(dataPtr.GetField("transaction_num"))
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			block.Confirmed = true
		}

		blockInfos = append(blockInfos, block)
	}

	//查询该语句所查到的数据集合
	var total = int64(len(blockInfos))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	blocksResp.Total = total
	blocksResp.Data = blockInfos

	return blocksResp, nil

}

//QueryBlockRealize 操作数据库
func QueryBlockRealize(strSQL, filterSQL string) (*entity.BlockInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryBlocks error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryBlocks dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var block = &entity.BlockInfo{}
	//填充数据
	for dataPtr.NextT() {
		block.Number = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		block.Hash = dataPtr.GetField("block_hash")
		block.Size = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_size"))
		block.CreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		block.TxTrieRoot = dataPtr.GetField("tx_trie_hash")
		block.ParentHash = dataPtr.GetField("parent_hash")
		block.WitnessAddress = dataPtr.GetField("witness_address")
		block.WitnessID = 0
		block.NrOfTrx = mysql.ConvertDBValueToInt64(dataPtr.GetField("transaction_num"))
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			block.Confirmed = true
		}
	}

	return block, nil

}
