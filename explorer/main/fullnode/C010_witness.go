package main

import (
	"fmt"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

func startWintnessDaemon() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if witnessList, ok := getWitness(); ok {
				icnt, ucnt, ecnt, err := storeWitness(witnessList)
				fmt.Printf("witness work result:(%v,%v,%v,%v)\n", icnt, ucnt, ecnt, err)
				time.Sleep(30 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}

			if needQuit() {
				break
			}
		}
		fmt.Printf("Witness Daemon QUIT\n")
	}()
}

func getWitness() ([]*core.Witness, bool) {
	client := grpcclient.GetRandomSolidity()

	witnessList, err := client.ListWitnesses()
	if nil != err || len(witnessList) == 0 {
		return nil, false
	}

	return witnessList, true
}

func storeWitness(witnessList []*core.Witness) (iCnt int64, uCnt int64, eCnt int64, err error) {
	if len(witnessList) == 0 {
		return
	}

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return
	}
	/*
			CREATE TABLE `witness` (
		  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
		  `vote_count` bigint(20) DEFAULT '0' COMMENT '得票数',
		  `public_key` varchar(300) DEFAULT '' COMMENT '公钥',
		  `url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
		  `total_produced` bigint(20) DEFAULT '0' COMMENT '生产块数',
		  `total_missed` bigint(20) DEFAULT '0' COMMENT '丢失块数',
		  `latest_block_num` bigint(20) DEFAULT '0',
		  `latest_slot_num` bigint(20) DEFAULT '0',
		  `is_job` tinyint(4) DEFAULT '0' COMMENT '是否为超级候选人 0:false, 1:true',
		  PRIMARY KEY (`address`),
		  UNIQUE KEY `address_UNIQUE` (`address`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	*/

	sqlI := "insert into witness (address, vote_count, public_key, url, total_produced, total_missed, latest_block_num, latest_slot_num, is_job) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtI, err := txn.Prepare(sqlI)
	if nil != err {
		fmt.Printf("prepare insert witness failed:%v\n", err)
		return
	}
	defer stmtI.Close()

	// 修正更新witness逻辑，当total_produced >= 当前值，且 latest_block_num >= 当前值时，才更新数据，否则不更新数据
	sqlU := "update witness set vote_count = ?, public_key = ?, url = ?, total_produced = ?, total_missed = ?, latest_block_num = ?, latest_slot_num = ?, is_job = ? where address = ? and total_produced <= ? and latest_block_num <= ?"
	stmtU, err := txn.Prepare(sqlU)
	if nil != err {
		fmt.Printf("prepare update witness failed:%v\n", err)
		return
	}
	defer stmtU.Close()

	for _, witness := range witnessList {
		if nil == witness {
			eCnt++
			continue
		}

		addr := utils.Base58EncodeAddr(witness.Address)
		_, err = stmtI.Exec(
			addr,
			witness.VoteCount,
			utils.HexEncode(witness.PubKey),
			witness.Url,
			witness.TotalProduced,
			witness.TotalMissed,
			witness.LatestBlockNum,
			witness.LatestSlotNum,
			witness.IsJobs,
		)

		if err != nil {
			// fmt.Printf("new witness failed:%v-->%v\n", err, addr)

			_, err = stmtU.Exec(
				witness.VoteCount,
				utils.HexEncode(witness.PubKey),
				witness.Url,
				witness.TotalProduced,
				witness.TotalMissed,
				witness.LatestBlockNum,
				witness.LatestSlotNum,
				witness.IsJobs,
				addr,
				witness.TotalProduced,
				witness.LatestBlockNum,
			)

			if nil != err {
				eCnt++
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
