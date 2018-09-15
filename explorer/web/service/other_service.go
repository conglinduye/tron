package service

import (
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QuerySystemStatus ...
func QuerySystemStatus() (*entity.SystemStatusResp, error) {
	var systemStatusResp = &entity.SystemStatusResp{}
	var confirmBlockIDDb int64
	netType := &entity.Network{Type: "mainnet"} //暂时写死，用的都是主网数据

	blockBuffer := buffer.GetBlockBufferInstance()

	//查询数据库按时间倒序排列的最近确认区块高度--从buffer获取
	confirmBlockDb, err := blockBuffer.GetBlocks(-1, 0, 1)
	if err != nil || confirmBlockDb == nil {
		log.Errorf("get block data err from buffer, get in db instead")
		confirmBlock, err := QueryBlocks(&entity.Blocks{Order: "-timestamp", Limit: 1, Start: 0})
		if err != nil {
			log.Errorf("QueryBlocks in QuerySystemStatus err:%v", err)
			return nil, err
		}
		confirmBlockDb = confirmBlock.Data
	}
	for _, blockinfo := range confirmBlockDb {
		confirmBlockIDDb = blockinfo.Number
		break
	}
	//查询数据库获取最大块包含非确认--从buffer获取（由于数据库中没存未确认的块，所以从fullnode获取）
	latestdBlockIDDb := blockBuffer.GetMaxBlockID()
	//从buffer中获取fullnode块高
	fullnodeNowBlockID := blockBuffer.GetUnConfirmedBlockID()
	//从buffer中获取db的确认过的块高
	solidityNowBlockID := blockBuffer.GetMaxConfirmedBlockID()

	//计算总进度
	solidityProcess := float64(confirmBlockIDDb) / float64(solidityNowBlockID) * 100
	fullnodeProcess := float64(latestdBlockIDDb) / float64(fullnodeNowBlockID) * 100
	totalProcess := (solidityProcess + fullnodeProcess) / 2
	log.Debugf("QuerySystemStatus,confirmBlockIDDb:[%.f],latestdBlockIDDb:[%.f],solidityNowBlockID:[%.f],fullnodeNowBlockID:[%.f]\n",
		confirmBlockIDDb, latestdBlockIDDb, solidityNowBlockID, fullnodeNowBlockID)
	log.Debugf("QuerySystemStatus,solidityProcess:[%.f],fullnodeProcess:[%.f],totalProcess:[%.f]\n",
		solidityProcess, fullnodeProcess, totalProcess)
	//拼接返回数据
	systemStatusResp.Network = netType
	systemStatusResp.Sync = &entity.Sync{Progress: totalProcess}
	systemStatusResp.Full = &entity.BlockNode{Block: fullnodeNowBlockID}
	systemStatusResp.Solidity = &entity.BlockNode{Block: solidityNowBlockID}
	systemStatusResp.Database = &entity.DataBase{Block: latestdBlockIDDb, ConfirmedBlock: confirmBlockIDDb}

	return systemStatusResp, nil
}

//QueryMarketsBuffer ... 从buffer获取市场信息
func QueryMarketsBuffer() ([]*entity.MarketInfo, error) {

	marketBuffer := buffer.GetMarketBuffer()
	markets := marketBuffer.GetMarket()

	return markets, nil
}

//QueryAuth ...
func QueryAuth() (*entity.TransfersResp, error) {

	return nil, nil
}

//QueryTestRequestCoin ...
func QueryTestRequestCoin() (*entity.TransfersResp, error) {

	return nil, nil
}
