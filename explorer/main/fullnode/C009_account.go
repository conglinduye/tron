package main

import (
	"fmt"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

type account struct {
	raw *core.Account
}

// getAccount addrs from redis which is the raw []byte, need convert to base58
func getAccount(addrs []string) ([]*account, []string, error) {

	client := grpcclient.GetRandomSolidity()

	maxErr := 0

	restAddr := make([]string, 0, len(addrs))
	accountList := make([]*account, 0, len(addrs))

	for _, addr := range addrs {
		acc, err := client.GetAccountRawAddr(([]byte(addr)))
		utils.VerifyCall(acc, err)
		if nil != err {
			maxErr++
			restAddr = append(restAddr, addr)
		}

		acct := new(account)
		acct.raw = acc
		accountList = append(accountList, acct)

	}

	return accountList, restAddr, nil
}

func storeAccount(accountList []*account) bool {
	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return false
	}
	/*
		CREATE TABLE `account` (
		  `account_name` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Account name',
		  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address',
		  `balance` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
		  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '账户创建时间',
		  `latest_operation_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '账户最后操作时间',
		  `is_witness` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
		  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
		  `fronze_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '冻结金额, 投票权',
		  `create_unix_time` int(32) NOT NULL DEFAULT '0' COMMENT '账户创建时间unix时间戳，用于分区',
		  UNIQUE KEY `uniq_account_address` (`address`,`create_unix_time`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	*/
	sqlI := "insert into account (account_name, address, balance, create_time, latest_operation_time, is_witness, fronze_amount, create_unix_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtI, err := txn.Prepare(sqlI)
	if nil != err {
		fmt.Printf("prepare insert account SQL failed:%v\n", err)
		return false
	}
	defer stmtI.Close()

	sqlU := "update account set account_name = ?, balance = ?, latest_operation_time = ?, is_witness = ?, fronze_amount = ? where create_unix_time = ? and address = ?"
	stmtU, err := txn.Prepare(sqlU)
	if nil != err {
		fmt.Printf("prepare update account SQL failed:%v\n", err)
		return false
	}
	defer stmtU.Close()

	for _, acc := range accountList {
		isWitness := 0
		if acc.raw.IsWitness {
			isWitness = 1
		}

		stmtI.Exec(
			acc.raw.AccountName,
			acc.raw.Address,
			acc.raw.Balance,
			acc.raw.CreateTime,
			acc.raw.LatestAssetOperationTime,
			isWitness,
			utils.ToJSONStr(acc.raw.Frozen),
			utils.ConvTimestamp(acc.raw.CreateTime),
		)

		if err != nil {
			fmt.Printf("insert into account failed:%v-->%v\n", err, utils.ToJSONStr(acc.raw))
			// return false

			stmtU.Exec(
				acc.raw.AccountName,
				acc.raw.Balance,
				acc.raw.LatestAssetOperationTime,
				isWitness,
				utils.ToJSONStr(acc.raw.Frozen),
				utils.ConvTimestamp(acc.raw.CreateTime),
				acc.raw.Address)
			if err != nil {
				fmt.Printf("update account failed:%v-->%v\n", err, utils.ToJSONStr(acc.raw))
			}
		} else {

		}

	}

	err = txn.Commit()
	if err != nil {
		fmt.Printf("connit block failed:%v\n", err)
		return false
	}

	return true
}
