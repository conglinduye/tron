package main

import (
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/main/module/account"
	"github.com/wlcy/tron/explorer/main/module/rawmysql"
)

// "github.com/wlcy/tron/explorer/utils"

var _accOnce sync.Once
var accWorker *account.SyncWorker

func startAccountDaemonNew() {
	_accOnce.Do(func() {
		rawmysql.DSN = *gStrMysqlDSN

		accWorker = account.NewAccountWorker(*gIntMaxWorker, *gAccountWorkerQueue, 1, *gAccUniqBufferTime, *gAccRecordPerCommit)
		accWorker.StartAccountWorker()
		accWorker.StartDBWorker()
	})

	wg.Add(1)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if needQuit() {
				break
			}
			accWorker.GetStatus()
		}
		wg.Done()
	}()
}

func syncAddress(addrList [][]byte) {
	accWorker.AppendTask2(addrList)
}

// AddRefreshAddress add address to acc worker
func AddRefreshAddress(addrs ...interface{}) (newLen int, err error) {
	list := make([][]byte, 0, len(addrs))
	for _, addr := range addrs {
		tmp, ok := addr.([]byte)
		if ok && len(tmp) > 0 {
			list = append(list, tmp)
		}
	}
	accWorker.AppendTask2(list)
	return len(list), nil
}

/*
func startAccountDaemon() {
	wc2 = newWorkerCounter(*gIntMaxWorker)
	wc3 = newWorkerCounter(*gIntMaxWorker)

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
		time.Sleep(1 * time.Second)
	}
	cleanAccountBuffer()
	list, err := ClearRefreshAddress() // load all address from redis and prepare handle it
	fmt.Printf("### total account need to synchronze:%-10v, err:%v, start synchronize account info ......\n", len(list), err)

	ts := time.Now()
	accList, restAddr, _ := getAccount(list)
	fmt.Printf("### total account syncrhonzed:%-10v, bad address:%-10v, cost:%v, synchronize to db .....\n", len(accList), len(restAddr), time.Since(ts))

	wg.Add(1)
	go func() {
		defer wg.Done()
		wc3.startOne()
		ts1 := time.Now()
		blukStoreAccount(accList)
		fmt.Printf("### store account size:%-10v to DB cost:%v\n", len(accList), time.Since(ts1))
		wc3.stopOne()
	}()
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

*/
