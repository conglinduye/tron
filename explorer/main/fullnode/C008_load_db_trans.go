package main

import (
	"fmt"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func loadTransFromDB(blockIDs []int64) []*transaction {

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("load transaction from db create transaction failed:%v\n", err)
		return nil
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
	sqlstr := "select trx_hash, block_id, contract_type, contract_data, result_data, create_time from transactions where block_id = ?"
	stmt, err := txn.Prepare(sqlstr)
	if nil != err {
		fmt.Printf("prepare store transaction SQL failed:%v\n", err)
		return nil
	}
	defer stmt.Close()

	trxList := make([]*transaction, 0, len(blockIDs)*10)

	failedBlockIDs := make([]int64, 0, len(blockIDs))
	for _, id := range blockIDs {

		rows, err := stmt.Query(id)

		if err != nil {
			fmt.Printf("ERROR: load transaction (block_id:%v) failed:%v\n", id, err)
			failedBlockIDs = append(failedBlockIDs, id)
			continue
		}

		for rows.Next() {
			trx := new(transaction)

			err := rows.Scan(&trx.hash, &trx.blockID, &trx.ctxType, &trx.ctxData, &trx.resultData, &trx.createTime)
			if nil != err {
				fmt.Printf("scan transaction failed:%v\n", err)
			}
			trxList = append(trxList, trx)
		}

	}

	return trxList
}

type ownerContract interface {
	GetOwnerAddress() []byte
}

type transaction struct {
	hash       string // hex encoding
	blockID    int64
	ctxType    int32
	ctxData    string      // hex encoding pb []byte
	resultData string      // hex encoding pb []byte
	createTime int64       // timestamp
	ownerAddr  string      // base58encoding
	contract   interface{} // contract real data type
}

// ExtractContract 解析的原始协议对象
func (trx *transaction) ExtractContract() bool {
	_, ctx := utils.GetContractByParamVal(core.Transaction_Contract_ContractType(trx.ctxType), utils.HexDecode(trx.ctxData))

	ownerCtx, ok := ctx.(ownerContract)
	if ok && nil != ownerCtx {
		trx.ownerAddr = utils.Base58EncodeAddr(ownerCtx.GetOwnerAddress())
	}
	trx.contract = ctx

	return nil != ctx
}
