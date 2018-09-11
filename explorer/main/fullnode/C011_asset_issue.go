package main

import (
	"fmt"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

func startAssetDaemon() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if witnessList, ok := getWitness(); ok {
				storeWitness(witnessList)
				time.Sleep(30 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}

			if needQuit() {
				break
			}
		}
		fmt.Printf("Asset Daemon QUIT\n")
	}()
}

func getAssets() ([]*core.AssetIssueContract, bool) {
	client := grpcclient.GetRandomSolidity()

	assetList, err := client.GetAssetIssueList()
	if nil != err || len(assetList) == 0 {
		return nil, false
	}

	return assetList, true
}

func storeAsset(assetList []*core.AssetIssueContract) (iCnt int64, uCnt int64, eCnt int64, err error) {
	if len(assetList) == 0 {
		return
	}

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return
	}
	/*
		CREATE TABLE `asset_issue` (
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `asset_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'asset_name',
		  `asset_abbr` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'asset_abbr',
		  `total_supply` bigint NOT NULL DEFAULT '0' COMMENT '发行量',
		  `frozen_supply` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '冻结量',
		  `trx_num` bigint NOT NULL DEFAULT '0' COMMENT 'price 分子',
		  `num` bigint not null default '0' comment 'price 分母',
		  `start_time` bigint not null default '0' comment '',
		  `end_time` bigint not null default '0' comment '',
		  `order_num` bigint not null default '0' comment '',
		  `vote_score` int not null default '0' comment '',
		  `description` text not null comment '',
		  `url` varchar(500) not null default '' comment '',
		  `free_asset_net_limit` bigint not null default '0' comment '',
		  `public_free_asset_net_limit` bigint not null default '0' comment '',
		  `public_free_asset_net_usage` bigint not null default '0' comment '',
		  `public_latest_free_net_time` bigint not null default '0' comment '',
		  PRIMARY KEY (`owner_address`,`asset_name`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	sqlI := `insert into asset_issue 
	(owner_address, asset_name, asset_abbr, total_supply, frozen_supply, trx_num, num, start_time, end_time, order_num, 
		vote_score, asset_desc, url, free_asset_net_limit, 
		public_free_asset_net_limit, public_free_asset_net_usage, public_latest_free_net_time) 
	values 
		 (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
		 ?, ?, ?, ?,
		 ?, ?, ?)`
	stmt, err := txn.Prepare(sqlI)
	if nil != err {
		fmt.Printf("prepare [%v] failed:%v\n", sqlI, err)
		return
	}
	defer stmt.Close()

	sqlU := `update asset_issue set asset_abbr = ?, total_supply = ?, frozen_supply = ?, trx_num = ?, num = ?, start_time = ?, end_time = ?, order_num = ?,
	vote_score = ?, asset_desc = ?, url = ?, free_asset_net_limit = ?, public_free_asset_net_limit = ?,
	public_free_asset_net_usage = ?, public_latest_free_net_time = ? 
	where owner_address = ? and asset_name = ?`

	stmtU, err := txn.Prepare(sqlU)
	if nil != err {
		fmt.Printf("prepare update witness failed:%v\n", err)
		return
	}
	defer stmtU.Close()

	for _, asset := range assetList {
		_, err = stmt.Exec(
			utils.Base58EncodeAddr(asset.OwnerAddress),
			string(asset.Name),
			string(asset.Abbr),
			asset.TotalSupply,
			utils.ToJSONStr(asset.FrozenSupply),
			asset.TrxNum,
			asset.Num,
			asset.StartTime,
			asset.EndTime,
			asset.Order,
			asset.VoteScore,
			utils.HexEncode(asset.Description),
			string(asset.Url),
			asset.FreeAssetNetLimit,
			asset.PublicFreeAssetNetLimit,
			asset.PublicFreeAssetNetUsage,
			asset.PublicLatestFreeNetTime)
		if nil != err {
			fmt.Printf("insert contract [%#v] failed:%v\n", asset, err)

			_, err = stmtU.Exec(
				string(asset.Abbr),
				asset.TotalSupply,
				utils.ToJSONStr(asset.FrozenSupply),
				asset.TrxNum,
				asset.Num,
				asset.StartTime,
				asset.EndTime,
				asset.Order,
				asset.VoteScore,
				string(asset.Description),
				string(asset.Url),
				asset.FreeAssetNetLimit,
				asset.PublicFreeAssetNetLimit,
				asset.PublicFreeAssetNetUsage,
				asset.PublicLatestFreeNetTime,
				utils.Base58EncodeAddr(asset.OwnerAddress),
				string(asset.Name))
			if nil != err {
				eCnt++
				fmt.Printf("update asset failed:%v --->%v\n", err, utils.ToJSONStr(asset))
			} else {
				uCnt++
			}
		} else {
			iCnt++
		}
	}

	err = txn.Commit()
	if err != nil {
		// fmt.Printf("connit witness failed:%v\n", err)
		return
	}
	return
}
