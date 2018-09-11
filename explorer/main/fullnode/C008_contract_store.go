package main

import (
	"database/sql"
	"fmt"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func storeContractDetail(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction) {
	if nil == txn || nil == trx || nil == trx.RawData || len(trx.RawData.Contract) == 0 {
		return
	}

	_, ctx := utils.GetContract(trx.RawData.Contract[0])

	switch v := ctx.(type) {
	case *core.AccountCreateContract:
		storeAccountCreateContract(txn, confirmed, trxHash, trx, v)

	case *core.AccountUpdateContract:
		storeAccountUpdateContract(txn, confirmed, trxHash, trx, v)

	case *core.SetAccountIdContract:
		storeSetAccountIDContract(txn, confirmed, trxHash, trx, v)

	case *core.TransferContract:
		storeTransferContract(txn, confirmed, trxHash, trx, v)

	case *core.TransferAssetContract:
		storeTransferAssetContract(txn, confirmed, trxHash, trx, v)

	case *core.VoteAssetContract:
		storeVoteAssetContract(txn, confirmed, trxHash, trx, v)

	case *core.VoteWitnessContract:
		storeVoteWitnessContract(txn, confirmed, trxHash, trx, v)

	case *core.VoteWitnessContract_Vote: // no api use this type, no OwnerAddress
	case *core.UpdateSettingContract:
		storeUpdateSettingContract(txn, confirmed, trxHash, trx, v)
	case *core.WitnessCreateContract:
		storeWitnessCreateContract(txn, confirmed, trxHash, trx, v)

	case *core.WitnessUpdateContract:
		storeWitnessUpdateContract(txn, confirmed, trxHash, trx, v)
	case *core.AssetIssueContract:
		storeAssetIssueContract(txn, confirmed, trxHash, trx, v)

	case *core.AssetIssueContract_FrozenSupply: // no api use this type, no OwnerAddress
	case *core.ParticipateAssetIssueContract:
		storeParticipateAssetIssueContract(txn, confirmed, trxHash, trx, v)

	case *core.FreezeBalanceContract:
		storeFreezeBalanceContract(txn, confirmed, trxHash, trx, v)

	case *core.UnfreezeBalanceContract:
		storeUnfreezeBalanceContract(txn, confirmed, trxHash, trx, v)

	case *core.UnfreezeAssetContract:
		storeUnfreezeAssetContract(txn, confirmed, trxHash, trx, v)

	case *core.WithdrawBalanceContract:
		storeWithdrawBalanceContract(txn, confirmed, trxHash, trx, v)

	case *core.UpdateAssetContract:
		storeUpdateAssetContract(txn, confirmed, trxHash, trx, v)

	case *core.ProposalCreateContract:
	case *core.ProposalApproveContract:
	case *core.ProposalDeleteContract:
	case *core.CreateSmartContract:
		storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.TriggerSmartContract:
		storeTriggerSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.BuyStorageContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.BuyStorageBytesContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.SellStorageContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.ExchangeCreateContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.ExchangeInjectContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.ExchangeWithdrawContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	case *core.ExchangeTransactionContract:
		// storeCreateSmartContract(txn, confirmed, trxHash, trx, v)
	default:
		fmt.Printf("new type:%T-->%v\n", v, v)
	}

}

func storeAccountCreateContract(txn *sql.Tx, confiremd int, trxHash string, trx *core.Transaction, ctx *core.AccountCreateContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}

	/*
		CREATE TABLE `contract_account_create` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
		  `block_id` bigint(20) NOT NULL DEFAULT '0',
		  `contract_type` int(11) NOT NULL DEFAULT '0',
		  `create_time` bigint(20) NOT NULL DEFAULT '0',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint not null default 0,
		  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
		  `new_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
		  `account_type` tinyint(4) NOT NULL DEFAULT '0',
		  PRIMARY KEY (`trx_hash`,`block_id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/

	_, err = txn.Exec(`insert into contract_account_create 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, account_address, account_type) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confiremd,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.AccountAddress),
		ctx.Type)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}

	AddRefreshAddress(ctx.OwnerAddress, ctx.AccountAddress)

	return
}

func storeTransferContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.TransferContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_transfer` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
		  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
		  `token_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token名称',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_transfer 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, to_address, amount, token_name) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ToAddress),
		ctx.Amount,
		"")
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.ToAddress)

	return
}

func storeTransferAssetContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.TransferAssetContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_transfer` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
		  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
		  `token_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token名称',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_transfer 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, to_address, amount, token_name) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ToAddress),
		ctx.Amount,
		string(ctx.AssetName))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}

	_, err = txn.Exec(`insert into contract_asset_transfer 
	(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
		owner_address, to_address, amount, token_name) 
	values 
	(?, ?, ?, ?, ?, ?,
		 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ToAddress),
		ctx.Amount,
		string(ctx.AssetName))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.ToAddress)

	return
}

func storeVoteWitnessContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.VoteWitnessContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_vote_witness` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `votes` TEXT COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投票详情，JSON',
		  `support` tinyint NOT NULL DEFAULT '0' COMMENT 'support 0:false, 1:true',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_vote_witness 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, votes, support) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.ToJSONStr(ctx.Votes),
		ctx.Support)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}

	AddRefreshAddress(ctx.OwnerAddress)
	return
}

func storeWitnessCreateContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.WitnessCreateContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_witness_create` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'url',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_witness_create 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, url) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.Url))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)
	return
}

func storeAssetIssueContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.AssetIssueContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_asset_issue` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
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
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_asset_issue 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, asset_name, asset_abbr, total_supply, frozen_supply, trx_num, num, start_time, end_time, order_num, 
			vote_score, asset_desc, url, free_asset_net_limit, 
			public_free_asset_net_limit, public_free_asset_net_usage, public_latest_free_net_time) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			 ?, ?, ?, ?,
			 ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.Name),
		string(ctx.Abbr),
		ctx.TotalSupply,
		utils.ToJSONStr(ctx.FrozenSupply),
		ctx.TrxNum,
		ctx.Num,
		ctx.StartTime,
		ctx.EndTime,
		ctx.Order,
		ctx.VoteScore,
		string(ctx.Description),
		string(ctx.Url),
		ctx.FreeAssetNetLimit,
		ctx.PublicFreeAssetNetLimit,
		ctx.PublicFreeAssetNetUsage,
		ctx.PublicLatestFreeNetTime)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)
	return
}

func storeParticipateAssetIssueContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.ParticipateAssetIssueContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_participate_asset` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `asset_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token名称',
		  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_participate_asset 
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, to_address, asset_name, amount) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ToAddress),
		string(ctx.AssetName),
		ctx.Amount)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.ToAddress)
	return
}

func storeFreezeBalanceContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.FreezeBalanceContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_freeze_balance` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `frozen_balance` bigint(20) NOT NULL DEFAULT '0',
		  `frozen_duration` bigint(20) NOT NULL DEFAULT '0',
		  `resource` int NOT NULL DEFAULT '0',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_freeze_balance
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, frozen_balance, frozen_duration, resource) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		ctx.FrozenBalance,
		ctx.FrozenDuration,
		ctx.Resource)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeUnfreezeBalanceContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.UnfreezeBalanceContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_unfreeze_balance` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `resource` int NOT NULL DEFAULT '0',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_unfreeze_balance
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, resource) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		ctx.Resource)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeWithdrawBalanceContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.WithdrawBalanceContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_freeze_balance` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_freeze_balance
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address)
		values 
		(?, ?, ?, ?, ?, ?,
			 ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeUnfreezeAssetContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.UnfreezeAssetContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_unfreeze_asset` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_unfreeze_asset
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)
	return
}

func storeAccountUpdateContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.AccountUpdateContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_account_update` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `account_name` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'account_name',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_account_update
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, account_name) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.AccountName))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeSetAccountIDContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.SetAccountIdContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_set_account_id` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `account_id` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'account_id',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_set_account_id
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.AccountId))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)
	return
}

func storeVoteAssetContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.VoteAssetContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_vote_asset` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `vote_address` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投票地址',
		  `support` tinyint NOT NULL DEFAULT '0' COMMENT 'support',
		  `count` bigint NOT NULL DEFAULT '0' COMMENT 'count',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_vote_asset
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, vote_address, support, count) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.ToJSONStr(ctx.VoteAddress),
		ctx.Support,
		ctx.Count)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeUpdateSettingContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.UpdateSettingContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_update_setting` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
		  `consume_user_resource_percent` bigint NOT NULL DEFAULT '0' COMMENT '',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_update_setting
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, contract_address, consume_user_resource_percent) 
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ContractAddress),
		ctx.ConsumeUserResourcePercent)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.ContractAddress)

	return
}

func storeWitnessUpdateContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.WitnessUpdateContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_witness_update` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `update_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_witness_update
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, update_url)
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.UpdateUrl))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeUpdateAssetContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.UpdateAssetContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_update_asset` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `asset_desc` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '',
		  `url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '',
		  `new_limit` bigint NOT NULL DEFAULT '0' COMMENT '',
		  `new_public_limit` bigint  NOT NULL DEFAULT '0' COMMENT '',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_update_asset
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, asset_desc, url, new_limit, new_public_limit)
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		string(ctx.Description),
		string(ctx.Url),
		ctx.NewLimit,
		ctx.NewPublicLimit)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress)

	return
}

func storeCreateSmartContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.CreateSmartContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_create_smart` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
		  `abi` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '',
		  `byte_code` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '',
		  `call_value` bigint(20) NOT NULL DEFAULT '0',
		  `consume_user_resource_percent` bigint(20) NOT NULL DEFAULT '0',
		  `name` varchar(500)  COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_create_smart
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, contract_address, abi, byte_code, call_value, consume_user_resource_percent, name)
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.NewContract.ContractAddress),
		utils.ToJSONStr(ctx.NewContract.Abi),
		utils.HexEncode(ctx.NewContract.Bytecode),
		ctx.NewContract.CallValue,
		ctx.NewContract.ConsumeUserResourcePercent,
		ctx.NewContract.Name)
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.NewContract.ContractAddress)
	return
}

func storeTriggerSmartContract(txn *sql.Tx, confirmed int, trxHash string, trx *core.Transaction, ctx *core.TriggerSmartContract) (err error) {
	if nil == txn || nil == trx || nil == ctx {
		return
	}
	/*
		CREATE TABLE `contract_trigger_smart` (
		  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
		  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
		  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
		  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
		  `expire_time` bigint(20) NOT NULL DEFAULT '0',
		  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
		  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
		  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
		  `call_value` bigint(20) NOT NULL DEFAULT '0',
		  `call_data` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '',
		  PRIMARY KEY (`trx_hash`,`block_id`),
		  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		/*!50100 PARTITION BY HASH (`block_id`)
		PARTITIONS 100 */
	/**/
	_, err = txn.Exec(`insert into contract_trigger_smart
		(trx_hash, block_id, contract_type, create_time, expire_time, confirmed, 
			owner_address, contract_address, call_value, call_data)
		values 
		(?, ?, ?, ?, ?, ?,
			 ?, ?, ?, ?)`,
		trxHash,
		trx.RawData.RefBlockNum,
		trx.RawData.Contract[0].Type,
		trx.RawData.Timestamp,
		trx.RawData.Expiration,
		confirmed,
		utils.Base58EncodeAddr(ctx.OwnerAddress),
		utils.Base58EncodeAddr(ctx.ContractAddress),
		ctx.CallValue,
		string(ctx.Data))
	if nil != err {
		fmt.Printf("insert contract [%v] failed:%v\n", utils.ToJSONStr(ctx), err)
	}
	AddRefreshAddress(ctx.OwnerAddress, ctx.ContractAddress)

	return
}
