package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

var wc2 *workerCounter

// getAccount addrs from redis which is the raw []byte, need convert to base58
func getAccount(addrs []string) ([]*account, []string, error) {
	wc2.startOne()

	totalTask := len(addrs)
	result := make([]*account, 0, len(addrs))
	badAddr := make([]string, 0, len(addrs))
	lock := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	if len(addrs) > getAccountWorkerLimit {
		go getAcoountF(addrs[0:len(addrs)-getAccountWorkerLimit], &result, &badAddr, lock, wg)
		addrs = addrs[len(addrs)-getAccountWorkerLimit:]
	}

	client := grpcclient.GetRandomWallet()
	defer func() {
		if nil != client {
			client.Close()
		}
	}()
	// client1 := grpcclient.GetRandomWallet()

	errCnt := 0

	restAddr := make([]string, 0, len(addrs))
	accountList := make([]*account, 0, len(addrs))
	bad := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		if !utils.VerifyTronAddrByte([]byte(addr)) {
			bad = append(bad, addr)
			continue
		}

		acc, err := client.GetAccountRawAddr(([]byte(addr)))
		if nil != err || nil == acc || len(acc.Address) == 0 {
			errCnt++
			restAddr = append(restAddr, addr)
			if errCnt >= maxErrCnt {
				client.Close()
				client = grpcclient.GetRandomWallet()
				errCnt = 0
			}
			continue
		}

		acct := new(account)
		acct.SetRaw(acc)
		accountList = append(accountList, acct)
	}

	if len(restAddr) > 0 {
		go getAcoountF(restAddr, &result, &badAddr, lock, wg)
	}

	waitCnt := 3
	lock.Lock()
	result = append(result, accountList...)
	badAddr = append(badAddr, bad...)
	lock.Unlock()
	// fmt.Printf("*** account, working task:%-05v, finished:%-06v, badAddr:%-06v, waitCnt:%v\n", wc2.currentWorker(), len(result), len(badAddr), waitCnt)

	for {
		workCnt := wc2.currentWorker()
		lock.Lock()
		fmt.Printf("*** account, working task:%-5v, finished:%-10v, total:%-10v, badAddr:%-10v, waitCnt:%v\n", workCnt, len(result), totalTask, len(badAddr), waitCnt)
		lock.Unlock()

		if workCnt == 1 && len(result)+len(badAddr) >= totalTask {
			waitCnt--
		}
		if waitCnt <= 0 {
			break
		}
		time.Sleep(3 * time.Second)
	}

	// storeAccount(accountList)

	wc2.stopOne()

	process := int64(0)
	getAccountNet(result, &process, lock)

	return result, badAddr, nil
}

func getAccountNet(accc []*account, process *int64, lock *sync.Mutex) {
	if len(accc) == 0 {
		return
	}
	fmt.Printf("*** accountNet start to syncrhonize accountNet info, total account:%v......\n", len(accc))
	wc2.startOne()
	totalTask := int64(len(accc))
	client := grpcclient.GetRandomWallet()
	defer func() {
		if nil != client {
			client.Close()
		}
	}()
	// ts := time.Now()
	errCnt := 0

	addrsLen := len(accc)
	restLen := addrsLen - getAccountWorkerLimit
	if restLen > 0 {
		// fmt.Printf("fork task %v~%v\n", 0, restLen)
		go getAccountNetF(accc[0:restLen], process, lock)
		accc = accc[restLen:]
	}

	restAcc := make([]*account, 0, len(accc))
	for idx, acc := range accc {

		accNet, err := client.GetAccountNetRawAddr(acc.raw.Address)
		if nil != err || nil == accNet {
			errCnt++
			if errCnt > maxErrCnt {
				restAcc = append(restAcc, accc[idx:]...)
				break
			} else {
				restAcc = append(restAcc, acc)
			}
			continue
		}

		acc.SetNetRaw(accNet)
	}
	if len(restAcc) > 0 {
		go getAccountNetF(restAcc, process, lock)
	}

	lock.Lock()
	*process = *process + int64(len(accc)-len(restAcc))
	// fmt.Printf("submit accountNet count:%v, current account result count:%v, restAddr:%v, error count:%v, cost:%v\n", len(accc)-len(restAcc), *process, len(restAcc), errCnt, time.Since(ts))
	lock.Unlock()

	waitCnt := 3

	for {
		workCnt := wc2.currentWorker()
		lock.Lock()
		fmt.Printf("*** accountNet, working task:%-05v, finished:%-06v, total:%-06v, waitCnt:%v\n", workCnt, *process, totalTask, waitCnt)
		lock.Unlock()

		if workCnt == 1 && *process >= totalTask {
			waitCnt--
		}
		if waitCnt <= 0 {
			break
		}
		time.Sleep(3 * time.Second)
	}

	wc2.stopOne()
	return
}

func getAccountNetF(accc []*account, process *int64, lock *sync.Mutex) {
	wc2.startOne()
	client := grpcclient.GetRandomWallet()
	defer func() {
		if nil != client {
			client.Close()
		}
	}()
	// ts := time.Now()
	errCnt := 0

	addrsLen := len(accc)
	restLen := addrsLen - getAccountWorkerLimit
	if restLen > 0 {
		// fmt.Printf("fork task %v~%v\n", 0, restLen)
		go getAccountNetF(accc[0:restLen], process, lock)
		accc = accc[restLen:]
	}

	restAcc := make([]*account, 0, len(accc))
	for idx, acc := range accc {

		accNet, err := client.GetAccountNetRawAddr(acc.raw.Address)
		if nil != err || nil == accNet {
			errCnt++
			if errCnt > maxErrCnt {
				restAcc = append(restAcc, accc[idx:]...)
				break
			} else {
				restAcc = append(restAcc, acc)
			}
			continue
		}

		acc.SetNetRaw(accNet)
	}

	lock.Lock()
	*process = *process + int64(len(accc)-len(restAcc))
	// fmt.Printf("submit accountNet count:%v, current account result count:%v, restAddr:%v, error count:%v, cost:%v\n", len(accc)-len(restAcc), *process, len(restAcc), errCnt, time.Since(ts))
	lock.Unlock()
	if len(restAcc) > 0 {
		go getAccountNetF(restAcc, process, lock)
	}
	// storeAccount(accountList)
	// fmt.Printf("getaccount handle address count:%v, cost:%v\n", len(accountList), time.Since(ts))

	wc2.stopOne()
	return
}

func getAcoountF(addrs []string, result *[]*account, badAddr *[]string, lock *sync.Mutex, wg *sync.WaitGroup) {
	wc2.startOne()
	client := grpcclient.GetRandomSolidity()
	defer func() {
		if nil != client {
			client.Close()
		}
	}()
	// client1 := grpcclient.GetRandomWallet()
	// fmt.Printf("getAccountFork task, address count:%v, client:%v\n", len(addrs), client.Target())

	// ts := time.Now()
	errCnt := 0

	restAddr := make([]string, 0, len(addrs))
	accountList := make([]*account, 0, len(addrs))

	addrsLen := len(addrs)
	restLen := addrsLen - getAccountWorkerLimit
	if restLen > 0 {
		// fmt.Printf("fork task %v~%v\n", 0, restLen)
		go getAcoountF(addrs[0:restLen], result, badAddr, lock, wg)
		addrs = addrs[restLen:]
	}

	bad := make([]string, 0, len(addrs))

	for idx, addr := range addrs {
		if !utils.VerifyTronAddrByte([]byte(addr)) {
			bad = append(bad, addr)
			continue
		}
		acc, err := client.GetAccountRawAddr(([]byte(addr)))
		if nil != err || nil == acc || len(acc.Address) == 0 {
			errCnt++
			if errCnt > maxErrCnt {
				restAddr = append(restAddr, addrs[idx:]...)
				break
			} else {
				restAddr = append(restAddr, addr)
			}
			continue
		}
		// accNet, err := client1.GetAccountNetRawAddr([]byte(addr))
		// if nil != err || nil == accNet {
		// 	errCnt++
		// 	restAddr = append(restAddr, addr)
		// 	if errCnt > maxErrCnt {
		// 		restAddr = append(restAddr, addrs[idx:]...)
		// 		break
		// 	}
		// }

		acct := new(account)
		acct.SetRaw(acc)
		// acct.SetNetRaw(accNet)
		accountList = append(accountList, acct)
	}

	lock.Lock()
	*result = append(*result, accountList...)
	*badAddr = append(*badAddr, bad...)
	// fmt.Printf("submit account count:%v, current account result count:%v, badAddr:%v, resetAddr:%v, error count:%v, cost:%v\n", len(accountList), len(*result), len(*badAddr), len(restAddr), errCnt, time.Since(ts))
	lock.Unlock()
	if len(restAddr) > 0 {
		go getAcoountF(restAddr, result, badAddr, lock, wg)
	}
	// storeAccount(accountList)
	// fmt.Printf("getaccount handle address count:%v, cost:%v\n", len(accountList), time.Since(ts))

	wc2.stopOne()
	return
}

func storeAccount(accountList []*account, dbb *sql.DB) bool {
	if nil == dbb {
		dbb = getMysqlDB()
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
			acc.raw.Balance,
			acc.raw.CreateTime,
			acc.raw.LatestOprationTime,
			acc.IsWitness,
			acc.AssetIssueName,
			acc.Fronzen,
			acc.raw.Allowance,
			acc.raw.LatestWithdrawTime,
			acc.raw.LatestConsumeTime,
			acc.raw.LatestConsumeFreeTime,
			acc.Votes,
			acc.raw.NetUsage,
			acc.freeNetUsed,
			acc.freeNetLimit,
			acc.netUsed,
			acc.netLimit,
			acc.totalNetLimit,
			acc.totalNetWeight,
			acc.AssetNetUsed,
			acc.AssetNetLimit,
			acc.raw.Type,
			utils.ToJSONStr(acc.raw.FrozenSupply),
			acc.raw.IsCommittee,
			utils.ToJSONStr(acc.raw.LatestAssetOperationTime),
			utils.ToJSONStr(acc.raw.AccountResource))

		if err != nil {
			// fmt.Printf("insert into account failed:%v-->[%v]\n", err, acc.Addr)

			result, err := stmtU.Exec(
				acc.Name,
				acc.raw.Balance,
				acc.raw.LatestOprationTime,
				acc.IsWitness,
				acc.AssetIssueName,
				acc.Fronzen,
				acc.raw.Allowance,
				acc.raw.LatestWithdrawTime,
				acc.raw.LatestConsumeTime,
				acc.raw.LatestConsumeFreeTime,
				acc.Votes,
				acc.raw.NetUsage,
				acc.freeNetUsed,
				acc.freeNetLimit,
				acc.netUsed,
				acc.netLimit,
				acc.totalNetLimit,
				acc.totalNetWeight,
				acc.AssetNetUsed,
				acc.AssetNetLimit,
				acc.raw.Type,
				utils.ToJSONStr(acc.raw.FrozenSupply),
				acc.raw.IsCommittee,
				utils.ToJSONStr(acc.raw.LatestAssetOperationTime),
				utils.ToJSONStr(acc.raw.AccountResource),
				acc.Addr,
				acc.raw.LatestOprationTime)

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

		for _, vote := range acc.raw.Votes {
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

// func updateTrxOwner(trxList []*transaction) bool {
// 	dbb := getMysqlDB()

// 	ts := time.Now()
// 	txn, err := dbb.Begin()
// 	if err != nil {
// 		fmt.Printf("get db failed:%v\n", err)
// 		return false
// 	}

// 	sqlU := "update transactions set owner_address = ? where block_id = ? and trx_hash = ?"
// 	stmt, err := txn.Prepare(sqlU)
// 	if nil != err {
// 		fmt.Printf("prepare update transaction owner address SQL failed:%v\n", err)
// 		return false
// 	}
// 	defer stmt.Close()

// 	for _, trx := range trxList {

// 		_, err := stmt.Exec(trx.ownerAddr, trx.blockID, trx.hash)

// 		if nil != err {
// 			fmt.Printf("update transaction owner failed:%v, blockID:%v, trx hash:[%v]\n", err, trx.blockID, trx.hash)
// 		}
// 	}

// 	err = txn.Commit()
// 	if nil != err {
// 		fmt.Printf("commit update transaction owner failed:%v\n", err)
// 		return false
// 	}
// 	fmt.Printf("update transaction owner count:%v, cost:%v\n", len(trxList), time.Since(ts))

// 	return true
// }

func getDBMaxBlockID() int64 {
	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if nil != err {
		fmt.Printf("start db transaction failed:%v\n", err)
		return 10000000
	}

	row, err := txn.Query("select max(block_id) from blocks")
	if nil != err {
		fmt.Printf("getDBMaxBlockID failed:%v, return 10000000 as default!\n", err)
		if nil != row {
			row.Close()
		}
		return 10000000
	}
	defer row.Close()

	for row.Next() {
		var blockID int64
		err = row.Scan(&blockID)
		if nil == err {
			return blockID
		}
		return 10000000
	}

	return 10000000
}
