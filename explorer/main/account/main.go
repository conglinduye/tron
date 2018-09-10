package main

import (
	"flag"
	"fmt"
	"time"
)

var gIntMaxWorker = flag.Int("worker", 10, "maximum worker for fetch blocks")
var gStrMysqlDSN = flag.String("dsn", "tron:tron@tcp(172.16.21.224:3306)/tron", "mysql connection string(DSN)")
var gInt64MaxWorkload = flag.Int64("workload", 10000, "maximum workload for each worker")
var gMaxBlockID = flag.Int64("max_block_id", 0, "max block")

func main() {
	flag.Parse()

	trxBulkBlockNum = *gInt64MaxWorkload

	initDB(*gStrMysqlDSN)
	initRedis([]string{"127.0.0.1:6379"})

	initWorkerChan()

	analyzeTrx(0, *gMaxBlockID)

}

var trxBulkBlockNum = int64(10000)

func analyzeTrx(b, e int64) {
	startWorker()

	if e == 0 {
		e = getDBMaxBlockID()
	}

	if e-b > trxBulkBlockNum {
		go analyzeTrxFrk(b, e-trxBulkBlockNum)

		b = e - trxBulkBlockNum
	}

	ts := time.Now()
	blockIDs := genVerifyBlockIDList(b, e) // 7408
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("block range:[%v] ~ [%v], load %v trans cost:%v\n", b, e, len(trxList), time.Since(ts))

	for _, trx := range trxList {
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}
	updateTrxOwner(trxList)

	waitCnt := 10
	for {
		workerCnt := workingTaskCnt()
		fmt.Printf("main routine, working task:%v, waitCnt:%v\n", workerCnt, waitCnt)
		if workerCnt == 1 {
			waitCnt--
			if waitCnt < 0 {
				break
			}

		}
		time.Sleep(3 * time.Second)

	}
	stopWorker()

	list, err := ClearRefreshAddress() // load all address from redis and prepare handle it
	fmt.Println(err)
	fmt.Printf("total account:%v\n", len(list))

	accList, restAddr, _ := getAccount(list)
	fmt.Printf("total account:%v, rest address:%v\n", len(accList), len(restAddr))
	storeAccount(accList)

	fmt.Printf("accList size:%v, restAddr size:%v\n", len(accList), len(restAddr))
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
	fmt.Printf("block range:[%v] ~ [%v], load %v trans cost:%v\n", b, e, len(trxList), time.Since(ts))

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
