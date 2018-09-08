package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"

	_ "github.com/go-sql-driver/mysql"
)

func storeTransactions(trans []*core.Transaction) bool {
	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("start transaction for storeTransaction failed:%v\n", err)
		return false
	}
	/*
		CREATE TABLE `transactions` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID，高度',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\nAccountCreateContract = 0;\r\nTransferContract = 1;\r\nTransferAssetContract = 2;\r\nVoteAssetContract = 3;\r\nVoteWitnessContract = 4;\r\nWitnessCreateContract = 5;\r\nAssetIssueContract = 6;\r\nWitnessUpdateContract = 8;\r\nParticipateAssetIssueContract = 9;\r\nAccountUpdateContract = 10;\r\nFreezeBalanceContract = 11;\r\nUnfreezeBalanceContract = 12;\r\nWithdrawBalanceContract = 13;\r\nUnfreezeAssetContract = 14;\r\nUpdateAssetContract = 15;\r\nProposalCreateContract = 16;\r\nProposalApproveContract = 17;\r\nProposalDeleteContract = 18;\r\nSetAccountIdContract = 19;\r\nCustomContract = 20;\r\n// BuyStorageContract = 21;\r\n// BuyStorageBytesContract = 22;\r\n// SellStorageContract = 23;\r\nCreateSmartContract = 30;\r\nTriggerSmartContract = 31;\r\nGetContract = 32;\r\nUpdateSettingContract = 33;\r\nExchangeCreateContract = 41;\r\nExchangeInjectContract = 42;\r\nExchangeWithdrawContract = 43;\r\nExchangeTransactionContract = 44;',
		  `contract_data` varchar(5000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易内容数据,原始数据byte hex encoding',
		  `result_data` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易结果对象byte hex encoding',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
		  `fee` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易花费 单位 sun',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `create_time` bigint NOT NULL DEFAULT 0 COMMENT '交易创建时间',
		  `expire_time` bigint NOT NULL DEFAULT 0 COMMENT '交易过期时间',
		  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_transactions_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/*
	 */
	sqlstr := "insert into transactions (trx_hash, block_id, contract_type, contract_data, result_data, create_time, expire_time) values (?, ?, ?, ?, ?, ?, ?)"
	stmt, err := txn.Prepare(sqlstr)
	if nil != err {
		fmt.Printf("prepare store transaction SQL failed:%v\n", err)
		return false
	}
	defer stmt.Close()

	for _, tran := range trans {
		if nil == tran || nil == tran.RawData {
			continue
		}
		if len(tran.RawData.Contract) > 0 {
			trxHash := utils.HexEncode(utils.CalcTransactionHash(tran))
			trxRetData := []byte{}
			if len(tran.Ret) > 0 {
				trxRetData, _ = proto.Marshal(tran.Ret[0])
			}
			_, err = stmt.Exec(
				trxHash,
				tran.RawData.RefBlockNum,
				tran.RawData.Contract[0].Type,
				utils.HexEncode(tran.RawData.Contract[0].Parameter.Value),
				utils.HexEncode(trxRetData),
				// utils.ConverTimestamp(tran.RawData.Timestamp))
				tran.RawData.Timestamp,
				tran.RawData.Expiration)
		} else {
			fmt.Println("ERROR: transaction contract is empty!")
		}
		if err != nil {
			fmt.Printf("ERROR: store transaction failed!%v, %#v\n", err, utils.ToJSONStr(tran))
			// return false
		}
	}

	err = txn.Commit()
	if err != nil {
		fmt.Printf("commit transaction data failed:%v\n", err)
		return false
	}

	return true
}

func storeBlocks(blocks []*core.Block) (bool, int64, int64, []int64) {
	dbb := getMysqlDB()
	ts := time.Now()
	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return false, 0, 0, nil
	}
	/*
		CREATE TABLE `blocks` (
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID。高度',
		  `block_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '区块hash',
		  `parent_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '区块父级hash',
		  `witness_address` varchar(300) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '代表节点地址',
		  `tx_trie_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '验证数根的hash值',
		  `block_size` int(32) DEFAULT '0' COMMENT '区块大小',
		  `transaction_num` int(32) DEFAULT '0' COMMENT '交易数',
		  `confirmed` tinyint(4) DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '区块创建时间',
		  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
		  PRIMARY KEY (`block_id`),
		  UNIQUE KEY `uniq_blocks_id` (`block_id` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100;

	*/
	sqlstr := "insert into blocks (block_id, block_hash, parent_hash, confirmed, transaction_num, block_size, witness_address, create_time, tx_trie_hash) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := txn.Prepare(sqlstr)
	if nil != err {
		fmt.Printf("prepare insert block SQL failed:%v\n", err)
		return false, 0, 0, nil
	}
	defer stmt.Close()

	tranList := make([]*core.Transaction, 0, len(blocks)*10)

	var succCnt, errCnt int64
	blockIDList := make([]int64, 0, len(blocks))

	for _, block := range blocks {
		if nil == block || nil == block.BlockHeader {
			continue
		}
		if true {
			blockHash := utils.HexEncode(utils.CalcBlockHash(block))
			blockIDList = append(blockIDList, block.BlockHeader.RawData.Number)
			data, _ := proto.Marshal(&api.TransactionList{Transaction: block.Transactions})
			_, err = stmt.Exec(
				block.BlockHeader.RawData.Number,
				blockHash,
				utils.HexEncode(block.BlockHeader.RawData.ParentHash),
				1,
				len(block.Transactions),
				len(data),
				utils.Base58EncodeAddr(block.BlockHeader.RawData.WitnessAddress),
				block.BlockHeader.RawData.Timestamp,
				utils.HexEncode(block.BlockHeader.RawData.TxTrieRoot))
		} else {
			fmt.Println("transaction contract is empty!")
		}
		if err != nil {
			fmt.Printf("insert into block failed:%v-->%v\n", err, utils.ToJSONStr(block))
			// return false
			errCnt++
		} else {
			succCnt++
			for _, tran := range block.Transactions {
				tran.RawData.RefBlockNum = block.BlockHeader.RawData.Number
			}
		}

		tranList = append(tranList, block.Transactions...)
	}

	err = txn.Commit()

	fmt.Printf("store %v blocks cost:%v\n", len(blocks), time.Since(ts))

	// ts = time.Now()
	// // storeTransactions(tranList)
	// fmt.Printf("store %v transactions cost:%v\n", len(tranList), time.Since(ts))

	if err != nil {
		fmt.Printf("connit block failed:%v\n", err)
		return false, succCnt, errCnt, blockIDList
	}
	return true, succCnt, errCnt, blockIDList
}
