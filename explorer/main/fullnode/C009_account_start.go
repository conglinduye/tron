package main

import (
	"fmt"
	"time"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func startAccountDaemon() {
	wc2 = newWorkerCounter(*gIntMaxWorker)

	wg.Add(1)
	go func() {
		defer wg.Done()

		ts := time.Now()
		interval := time.Second * time.Duration(*gIntHandleAccountInterval)
		for {
			if needQuit() {
				break
			}
			if time.Since(ts) < interval {
				time.Sleep(interval - time.Since(ts))
			}

			ts = time.Now()
			syncAccount()
		}

		syncAccount()
		fmt.Printf("Account Daemon QUIT\n")
	}()

}

func syncAccount() {
	if nil == wc2 {
		return
	}
	for wc2.currentWorker() > 0 { // wait
		time.Sleep(3 * time.Second)
	}
	cleanAccountBuffer()
	list, err := ClearRefreshAddress() // load all address from redis and prepare handle it
	fmt.Printf("### total account need to synchronze:%-10v, err:%v, start synchronize account info ......\n", len(list), err)

	ts := time.Now()
	accList, restAddr, _ := getAccount(list)
	fmt.Printf("### total account syncrhonzed:%-10v, bad address:%-10v, cost:%v, synchronize to db .....\n", len(accList), len(restAddr), time.Since(ts))

	ts1 := time.Now()
	blukStoreAccount(accList)
	fmt.Printf("### store account size:%-10v to DB cost:%v\n", len(accList), time.Since(ts1))
}

func blukStoreAccount(accList []*account) {
	pos := 0
	remain := len(accList)
	for remain > 0 {
		if remain >= maxTransPerTxn {
			storeAccount(accList[pos:pos+maxTransPerTxn], nil)
			pos += maxTransPerTxn
			remain -= maxTransPerTxn
			continue
		}
		storeAccount(accList[pos:pos+remain], nil)
		pos += remain
		remain -= remain
	}
}

type account struct {
	raw            *core.Account
	netRaw         *api.AccountNetMessage
	Name           string
	Addr           string
	CreateTime     int64
	IsWitness      int8
	Fronzen        string
	AssetIssueName string

	AssetBalance map[string]int64
	Votes        string

	// acccount net info
	freeNetUsed    int64
	freeNetLimit   int64
	netUsed        int64
	netLimit       int64
	totalNetLimit  int64
	totalNetWeight int64
	AssetNetUsed   string
	AssetNetLimit  string

	/*
		`account_name` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Account name',
		`address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address',
		`balance` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
		`create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户创建时间',
		`latest_operation_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户最后操作时间',
		`is_witness` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
		`frozen` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '冻结金额, 投票权',
		`create_unix_time` int(32) NOT NULL DEFAULT '0' COMMENT '账户创建时间unix时间戳，用于分区',
		`allowance` bigint(20) DEFAULT '0',
		`latest_withdraw_time` bigint(20) DEFAULT '0',
		`latest_consume_time` bigint(20) DEFAULT '0',
		`latest_consume_free_time` bigint(20) DEFAULT '0',
		`votes` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
	*/
}

// var maxErrCnt = 10
var getAccountWorkerLimit = 1000

var beginTime, _ = time.Parse("2006-01-02 15:03:04.999999", "2018-06-25 00:00:00.000000")

func (a *account) SetRaw(raw *core.Account) {
	a.raw = raw
	a.Name = string(raw.AccountName)
	a.Addr = utils.Base58EncodeAddr(raw.Address)
	a.AssetIssueName = string(raw.AssetIssuedName)
	a.CreateTime = raw.CreateTime
	if a.CreateTime == 0 {
		a.CreateTime = beginTime.UnixNano()
	}
	a.IsWitness = 0
	if raw.IsWitness {
		a.IsWitness = 1
	}
	if len(raw.Frozen) > 0 {
		a.Fronzen = utils.ToJSONStr(raw.Frozen)

	}
	a.AssetBalance = a.raw.Asset
	if len(raw.Votes) > 0 {
		a.Votes = utils.ToJSONStr(raw.Votes)
	}
}

func (a *account) SetNetRaw(netRaw *api.AccountNetMessage) {
	if nil == netRaw {
		return
	}
	a.netRaw = netRaw
	a.AssetNetUsed = utils.ToJSONStr(netRaw.AssetNetUsed)
	a.AssetNetLimit = utils.ToJSONStr(netRaw.AssetNetLimit)
	a.freeNetUsed = netRaw.FreeNetUsed
	a.freeNetLimit = netRaw.FreeNetLimit
	a.netLimit = netRaw.NetLimit
	a.netUsed = netRaw.NetUsed
	a.totalNetLimit = netRaw.TotalNetLimit
	a.totalNetWeight = netRaw.TotalNetWeight
}
