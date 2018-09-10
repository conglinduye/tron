package main

import (
	"flag"
	"fmt"
	"time"
)

var gIntMaxWorker = flag.Int("worker", 10, "maximum worker for fetch blocks")
var gStrMysqlDSN = flag.String("dsn", "tron:tron@tcp(172.16.21.224:3306)/tron", "mysql connection string(DSN)")
var gInt64MaxWorkload = flag.Int64("workload", 10000, "maximum workload for each worker")
var gMaxBlockID = flag.Int64("max_block_id", 0, "max block num, 0 for current block, -1 will continue latest and do analyze until kill")
var gMinBlockID = flag.Int64("start_block", 2200000, "block num start to analyze")

func main() {
	flag.Parse()

	trxBulkBlockNum = *gInt64MaxWorkload

	initDB(*gStrMysqlDSN)
	initRedis([]string{"127.0.0.1:6379"})

	initWorkerChan()

	if -1 == *gMaxBlockID {
		b := *gMinBlockID
		e := getDBMaxBlockID()
		for {
			fmt.Printf("Start account analyze for block range [%v] ~ [%v]\n", b, e)
			ts := time.Now()
			analyzeTrx(b, e)
			tsCost := time.Since(ts)
			if tsCost < time.Second*10 {
				time.Sleep(10*time.Second - tsCost)
			}
			b = e
			e = getDBMaxBlockID()
		}
	} else {
		fmt.Printf("Start account analyze for block range [%v] ~ [%v]\n", *gMinBlockID, *gMaxBlockID)
		analyzeTrx(*gMinBlockID, *gMaxBlockID)
	}

}

var trxBulkBlockNum = int64(10000)

func analyzeTrx(b, e int64) {
	startWorker()

	if e <= 0 {
		e = getDBMaxBlockID()
	}

	if e-b > trxBulkBlockNum {
		go analyzeTrxFrk(b, e-trxBulkBlockNum)

		b = e - trxBulkBlockNum
	}
	ts := time.Now()

	blockIDs := genVerifyBlockIDList(b, e) // 7408
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("block range:[%v] ~ [%v] (%v blocks), load %v trans cost:%v\n", b, e, e-b, len(trxList), time.Since(ts))

	if len(trxList) == 0 {
		stopWorker()
		return
	}

	for _, trx := range trxList {
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}
	updateTrxOwner(trxList)

	waitCnt := 3
	for {
		workerCnt := workingTaskCnt()
		fmt.Printf("main routine, working task:%v, waitCnt:%v\n", workerCnt, waitCnt)
		if workerCnt == 1 {
			waitCnt--
			if waitCnt <= 0 {
				break
			}

		}
		time.Sleep(3 * time.Second)

	}
	stopWorker()

	list, err := ClearRefreshAddress() // load all address from redis and prepare handle it
	fmt.Printf("total account:%v, err:%v\n", len(list), err)

	accList, restAddr, _ := getAccount(list)
	fmt.Printf("total account:%v, rest address:%v, cost:%v, synchronize to db .....\n", len(accList), len(restAddr), time.Since(ts))
	ts = time.Now()
	storeAccount(accList)
	fmt.Printf("accList size:%v, restAddr size:%v, synchronze to DB cost:%v\n", len(accList), len(restAddr), time.Since(ts))

}

func analyzeTrxFrk(b, e int64) {
	startWorker()
	if e-b > trxBulkBlockNum {
		go analyzeTrxFrk(b, e-trxBulkBlockNum)

		b = e - trxBulkBlockNum
	}

	ts := time.Now()
	blockIDs := genVerifyBlockIDList(b, e)
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("block range:[%v] ~ [%v] (%v blocks), load %v trans cost:%v\n", b, e, e-b, len(trxList), time.Since(ts))

	if len(trxList) == 0 {
		stopWorker()
		return
	}

	for _, trx := range trxList {
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}
	updateTrxOwner(trxList)

	stopWorker()
	return
}

var maxWorker chan struct{}

func initWorkerChan() {
	*gIntMaxWorker = 20
	maxWorker = make(chan struct{}, *gIntMaxWorker)
	for i := 0; i < *gIntMaxWorker; i++ {
		maxWorker <- struct{}{}
	}
}

func startWorker() {
	<-maxWorker
}

func stopWorker() {
	maxWorker <- struct{}{}
}

func workingTaskCnt() int {
	return *gIntMaxWorker - len(maxWorker)
}
