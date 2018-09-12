package service

import (
	"strconv"
	"strings"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QuerySystemStatus ...
func QuerySystemStatus() (*entity.SystemStatusResp, error) {
	var systemStatusResp = &entity.SystemStatusResp{}
	var latestBlockIDDb int64
	netType := &entity.Network{Type: "mainnet"} //TODO 从配置文件中获取
	// 查询数据库按时间倒序排列的最近区块高度
	latestBlockDb, err := QueryBlocks(&entity.Blocks{Order: "-timestamp", Limit: "1", Start: "0"})
	if err != nil {
		log.Errorf("QueryBlocks in QuerySystemStatus err:%v", err)
		return nil, err
	}
	for _, blockinfo := range latestBlockDb.Data {
		latestBlockIDDb = blockinfo.Number
		break
	}

	//查询soliditynode最近的区块高度
	solidityClient := grpcclient.GetRandomSolidity()
	solidityBlock, err := solidityClient.GetNowBlock()
	if err != nil {
		log.Errorf("getNowBlock from solidity in QuerySystemStatus err:%v", err)
		return nil, err
	}
	solidityNowBlockID := solidityBlock.BlockHeader.RawData.Number

	//查询fullnode最近的区块高度
	fullnodeClient := grpcclient.GetRandomWallet()
	fullnodeBlock, err := fullnodeClient.GetNowBlock()
	if err != nil {
		log.Errorf("getNowBlock from fullnode in QuerySystemStatus err:%v", err)
		return nil, err
	}
	fullnodeNowBlockID := fullnodeBlock.BlockHeader.RawData.Number

	//计算总进度
	solidityProcess := float64(latestBlockIDDb) / float64(solidityNowBlockID) * 100
	fullnodeProcess := float64(latestBlockIDDb) / float64(fullnodeNowBlockID) * 100
	totalProcess := (solidityProcess + fullnodeProcess) / 2
	//拼接返回数据
	systemStatusResp.Network = netType
	systemStatusResp.Sync = &entity.Sync{Progress: totalProcess}
	systemStatusResp.Full = &entity.BlockNode{Block: fullnodeNowBlockID}
	systemStatusResp.Solidity = &entity.BlockNode{Block: solidityNowBlockID}

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

//QueryStatsOverview ...
func QueryStatsOverview() (*entity.TransfersResp, error) {

	return nil, nil
}

//QueryAuth ...
func QueryAuth() (*entity.TransfersResp, error) {

	return nil, nil
}

//QueryTestRequestCoin ...
func QueryTestRequestCoin() (*entity.TransfersResp, error) {

	return nil, nil
}
