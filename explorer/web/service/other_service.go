package service

import (
	"strconv"
	"strings"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QuerySystemStatus ...
func QuerySystemStatus() (*entity.SystemStatusResp, error) {
	var systemStatusResp = &entity.SystemStatusResp{}
	var latestBlockIDDb int64
	netType := &entity.Network{Type: "mainnet"} //TODO 从配置文件中获取

	blockBuffer := buffer.GetBlockBuffer()

	// 查询数据库按时间倒序排列的最近区块高度
	/*latestBlockDb, err := QueryBlocks(&entity.Blocks{Order: "-timestamp", Limit: "1", Start: "0"})
	if err != nil {
		log.Errorf("QueryBlocks in QuerySystemStatus err:%v", err)
		return nil, err
	}
	for _, blockinfo := range latestBlockDb.Data {
		latestBlockIDDb = blockinfo.Number
		break
	}*/
	//查询soliditynode最近的区块高度
	solidityClient := grpcclient.GetRandomSolidity()
	solidityBlock, err := solidityClient.GetNowBlock()
	if nil != err || nil == solidityBlock || nil == solidityBlock.BlockHeader || nil == solidityBlock.BlockHeader.RawData {
		solidityClient = grpcclient.GetRandomSolidity()
		log.Errorf("reset wallet connection, new client:%v!!!\n", solidityClient.Target())
		log.Errorf("getNowBlock from solidity in QuerySystemStatus err:%v", err)
		return nil, err
	}

	solidityNowBlockID := solidityBlock.BlockHeader.RawData.Number

	//查询fullnode最近的区块高度
	fullnodeClient := grpcclient.GetRandomWallet()
	fullnodeBlock, err := fullnodeClient.GetNowBlock()
	if nil != err || nil == fullnodeBlock || nil == fullnodeBlock.BlockHeader || nil == fullnodeBlock.BlockHeader.RawData {
		fullnodeClient = grpcclient.GetRandomWallet()
		log.Errorf("reset wallet connection, new client:%v!!!\n", fullnodeClient.Target())
		log.Errorf("getNowBlock from fullnode in QuerySystemStatus err:%v", err)
		return nil, err
	}
	fullnodeNowBlockID := fullnodeBlock.BlockHeader.RawData.Number

	//从buffer中获取fullnode块高
	latestBlockIDDb = blockBuffer.GetMaxBlockID()
	//从buffer中获取数据库中的确认过的块高
	comfirmedBlockIDDb := blockBuffer.GetMaxConfirmedBlockID()

	//计算总进度
	solidityProcess := float64(comfirmedBlockIDDb) / float64(solidityNowBlockID) * 100
	fullnodeProcess := float64(latestBlockIDDb) / float64(fullnodeNowBlockID) * 100
	totalProcess := (solidityProcess + fullnodeProcess) / 2
	//拼接返回数据
	systemStatusResp.Network = netType
	systemStatusResp.Sync = &entity.Sync{Progress: totalProcess}
	systemStatusResp.Full = &entity.BlockNode{Block: fullnodeNowBlockID}
	systemStatusResp.Solidity = &entity.BlockNode{Block: solidityNowBlockID}
	systemStatusResp.Database = &entity.DataBase{Block: latestBlockIDDb, ConfirmedBlock: solidityNowBlockID}

	return systemStatusResp, nil
}

//QueryMarkets 查询交易所信息  爬虫
func QueryMarkets() ([]*entity.MarketInfo, error) {
	marketInfos := make([]*entity.MarketInfo, 0)
	marketURL := "https://coinmarketcap.com/currencies/tron/"
	_, body, errs := gorequest.New().Get(marketURL).End()
	if errs != nil && len(errs) > 0 {
		log.Error(errs)
		return nil, errs[0]
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	doc.Find("#markets-table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		marketInfo := &entity.MarketInfo{}
		node := strconv.Itoa(i + 1)
		rank, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
		name, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Attr("data-sort")
		pair, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3)").Attr("data-sort")
		link, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3) > a").Attr("href")
		volume, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-usd")
		volumeNative, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-native")
		price, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(5) > span[class=price]").Attr("data-usd")
		volumePercentage, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(6)").Attr("data-sort")
		marketInfo.Rank = mysql.ConvertStringToInt64(rank, 0)
		marketInfo.Name = name
		marketInfo.Pair = pair
		marketInfo.Link = link
		marketInfo.Volume = mysql.ConvertStringToFloat(volume, 0)
		marketInfo.VolumeNative = mysql.ConvertStringToFloat(volumeNative, 0)
		marketInfo.VolumePercentage = mysql.ConvertStringToFloat(volumePercentage, 0)
		marketInfo.Price = mysql.ConvertStringToFloat(price, 0)
		marketInfos = append(marketInfos, marketInfo)
	})

	log.Debugf("market : parse page data done.")

	return marketInfos, nil
}

//QueryAuth ...
func QueryAuth() (*entity.TransfersResp, error) {

	return nil, nil
}

//QueryTestRequestCoin ...
func QueryTestRequestCoin() (*entity.TransfersResp, error) {

	return nil, nil
}
