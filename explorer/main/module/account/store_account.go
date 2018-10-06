package account

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/core/utils"
)

// StoreAccount 将accountList保存到数据库
func StoreAccount(accountList []*Account, dbb *sql.DB) bool {
	if nil == dbb {
		return false
	}

	ts := time.Now()
	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return false
	}
	/*
		CREATE TABLE `tron_account` (
			`account_name` varchar(300) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'Account name',
			`account_type` integer not null default '0' comment 'account type, 0 for common account, 2 for contract account',
			`address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address',
			`balance` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
			`create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户创建时间',
			`latest_operation_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户最后操作时间',
			`asset_issue_name` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
			`is_witness` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
			`allowance` bigint(20) NOT NULL DEFAULT '0',
			`latest_withdraw_time` bigint(20) NOT NULL DEFAULT '0',
			`latest_consume_time` bigint(20) NOT NULL DEFAULT '0',
			`latest_consume_free_time` bigint(20) NOT NULL DEFAULT '0',
			`frozen` text COLLATE utf8mb4_bin NOT NULL COMMENT '冻结信息',
			`votes` text COLLATE utf8mb4_bin NOT NULL,
			`free_net_used` bigint(20) NOT NULL DEFAULT '0',
			`free_net_limit` bigint(20) NOT NULL DEFAULT '0',
			`net_usage` bigint(20) NOT NULL DEFAULT '0' COMMENT 'bandwidth, get from frozen',
			`net_used` bigint(20) NOT NULL DEFAULT '0',
			`net_limit` bigint(20) NOT NULL DEFAULT '0',
			`total_net_limit` bigint(20) NOT NULL DEFAULT '0',
			`total_net_weight` bigint(20) NOT NULL DEFAULT '0',
			`asset_net_used` text COLLATE utf8mb4_bin NOT NULL,
			`asset_net_limit` text COLLATE utf8mb4_bin NOT NULL,
			`frozen_supply` text COLLATE utf8mb4_bin not NULL,
			`is_committee` tinyint not null default '0',
			`latest_asset_operation_time` text COLLATE utf8mb4_bin not null,
			`account_resource` text COLLATE utf8mb4_bin not null,
			PRIMARY KEY (`address`),
			KEY `idx_tron_account_create_time` (`create_time`),
			KEY `idx_account_name` (`account_name`),
			KEY `idx_account_address` (`address`),
			KEY `idx_account_type` (`account_type`)
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	*/
	sqlI := `insert into tron_account 
		(account_name, address, balance, create_time, latest_operation_time, is_witness, asset_issue_name,
			frozen, allowance, latest_withdraw_time, latest_consume_time, latest_consume_free_time, votes,
			net_usage, free_net_used,
			free_net_limit, net_used, net_limit, total_net_limit, total_net_weight, asset_net_used, asset_net_limit
			, account_type, frozen_supply, is_committee, latest_asset_operation_time, account_resource) values 
		(?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?)`
	stmtI, err := txn.Prepare(sqlI)
	if nil != err {
		fmt.Printf("prepare insert tron_account SQL failed:%v\n", err)
		return false
	}
	defer stmtI.Close()

	sqlU := `update tron_account set account_name = ?, balance = ?, latest_operation_time = ?, is_witness = ?, asset_issue_name = ?,
		frozen = ?, allowance = ?, latest_withdraw_time = ?, latest_consume_time = ?, latest_consume_free_time = ?, votes = ?, net_usage = ?,
		free_net_used = ?, free_net_limit = ?, net_used = ?, net_limit = ?, total_net_limit = ?, total_net_weight = ?, asset_net_used = ?, asset_net_limit = ?
		, account_type = ?, frozen_supply = ?, is_committee = ?, latest_asset_operation_time = ?, account_resource = ? 
		where address = ? and latest_operation_time <= ?`
	stmtU, err := txn.Prepare(sqlU)
	if nil != err {
		fmt.Printf("prepare update tron_account SQL failed:%v\n", err)
		return false
	}
	defer stmtU.Close()

	sqlBI := "insert into account_asset_balance (address, asset_name, balance) values (?, ?, ?)"
	stmtBI, err := txn.Prepare(sqlBI)
	if nil != err {
		fmt.Printf("prepare insert account_asset_balance SQL failed:%v\n", err)
		return false
	}
	defer stmtBI.Close()

	sqlVI := "insert into account_vote_result (address, to_address, vote) values (?, ?, ?)"
	stmtVI, err := txn.Prepare(sqlVI)
	if nil != err {
		fmt.Printf("prepare insert account_vote_result SQL failed:%v\n", err)
		return false
	}
	defer stmtVI.Close()

	insertCnt := 0
	updateCnt := 0
	errCnt := 0

	for _, acc := range accountList {

		_, err := stmtI.Exec(
			acc.Name,
			acc.Addr,
			acc.Raw.Balance,
			acc.Raw.CreateTime,
			acc.Raw.LatestOprationTime,
			acc.IsWitness,
			acc.AssetIssueName,
			acc.Fronzen,
			acc.Raw.Allowance,
			acc.Raw.LatestWithdrawTime,
			acc.Raw.LatestConsumeTime,
			acc.Raw.LatestConsumeFreeTime,
			acc.Votes,
			acc.Raw.NetUsage,
			acc.freeNetUsed,
			acc.freeNetLimit,
			acc.netUsed,
			acc.netLimit,
			acc.totalNetLimit,
			acc.totalNetWeight,
			acc.AssetNetUsed,
			acc.AssetNetLimit,
			acc.Raw.Type,
			utils.ToJSONStr(acc.Raw.FrozenSupply),
			acc.Raw.IsCommittee,
			utils.ToJSONStr(acc.Raw.LatestAssetOperationTime),
			utils.ToJSONStr(acc.Raw.AccountResource))

		if err != nil {
			// fmt.Printf("insert into account failed:%v-->[%v]\n", err, acc.Addr)

			result, err := stmtU.Exec(
				acc.Name,
				acc.Raw.Balance,
				acc.Raw.LatestOprationTime,
				acc.IsWitness,
				acc.AssetIssueName,
				acc.Fronzen,
				acc.Raw.Allowance,
				acc.Raw.LatestWithdrawTime,
				acc.Raw.LatestConsumeTime,
				acc.Raw.LatestConsumeFreeTime,
				acc.Votes,
				acc.Raw.NetUsage,
				acc.freeNetUsed,
				acc.freeNetLimit,
				acc.netUsed,
				acc.netLimit,
				acc.totalNetLimit,
				acc.totalNetWeight,
				acc.AssetNetUsed,
				acc.AssetNetLimit,
				acc.Raw.Type,
				utils.ToJSONStr(acc.Raw.FrozenSupply),
				acc.Raw.IsCommittee,
				utils.ToJSONStr(acc.Raw.LatestAssetOperationTime),
				utils.ToJSONStr(acc.Raw.AccountResource),
				acc.Addr,
				acc.Raw.LatestOprationTime)

			if err != nil {
				errCnt++
				// fmt.Printf("update account failed:%v-->[%v]\n", err, acc.Addr)
			} else {
				_ = result
				// _, err := result.RowsAffected()
				// if nil != err {
				// 	errCnt++
				// 	// fmt.Printf("update failed:%v, affectRow:%v--->%v\n", err, affectRow, acc.Addr)
				// } else {
				updateCnt++
				// }
				// fmt.Printf("update account ok!!!\n")
			}
		} else {
			insertCnt++
			// fmt.Printf("Insert account ok!!!\n")
		}

		result, err := txn.Exec("delete from account_asset_balance where address = ?", acc.Addr)
		_ = result

		for k, v := range acc.AssetBalance {
			_, err := stmtBI.Exec(acc.Addr, k, v)
			if nil != err {
				fmt.Printf("insert account_asset_balance failed:%v\n", err)
			}
		}

		result, err = txn.Exec("delete from account_vote_result where address = ?", acc.Addr)

		for _, vote := range acc.Raw.Votes {
			_, err := stmtVI.Exec(acc.Addr, utils.Base58EncodeAddr(vote.VoteAddress), vote.VoteCount)
			if nil != err {
				fmt.Printf("insert account_asset_balance failed:%v\n", err)
			}
		}

	}

	err = txn.Commit()
	if err != nil {
		fmt.Printf("connit block failed:%v\n", err)
		return false
	}
	fmt.Printf("store account OK, cost:%v, insertCnt:%v, updateCnt:%v, errCnt:%v, total source:%v\n", time.Since(ts), insertCnt, updateCnt, errCnt, len(accountList))

	return true
}
