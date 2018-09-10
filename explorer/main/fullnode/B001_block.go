package main

import (
	"flag"
	"time"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var gIntMaxWorker = flag.Int("worker", 10, "maximum worker for fetch blocks")
var gStrMysqlDSN = flag.String("dsn", "tron:tron@tcp(172.16.21.224:3306)/tron", "mysql connection string(DSN)")
var gInt64MaxWorkload = flag.Int64("workload", 10000, "maximum workload for each worker")
var gStartBlokcID = flag.Int64("start_block", 0, "block num start to synchronize")

func main() {
	flag.Parse()

	initDB(*gStrMysqlDSN)

	initWorkerChan()

	getAllBlocks()
}

func initWorkerChan() {
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

func getAllBlocks() {
	ts := time.Now()
	getBlock(0, *gStartBlokcID, 0)
	fmt.Printf("get all blocks cost:%v\n", time.Since(ts))
}

// var maxWorkload = int64(10000) // getBlock任务最大工作负荷

var maxWorker chan struct{}

// func getBlock(id int, numStart int64, numEnd int64) int64 {
// 	startWorker()

// 	ts := time.Now()

// 	servAddr := fmt.Sprintf("%s:50051", utils.GetRandFullNode())

// 	taskID := fmt.Sprintf("[%04v|%v~%v|%v]", id, numStart, numEnd, servAddr)

// 	client := grpcclient.NewWallet(servAddr)
// 	client.Connect()
// 	dbc := grpcclient.NewDatabase(servAddr)
// 	dbc.Connect()

// 	blockCnt := int64(0)
// 	latestNum := int64(0)
// 	limit := int64(100)

// 	if latestNum == 0 || numStart >= latestNum {
// 		latestNum = getLatestNum(dbc)
// 		fmt.Printf("%v latestNum is [%v]\n", taskID, latestNum)
// 	}

// 	if latestNum == 0 {
// 		stopWorker()
// 		return getBlock(id, numStart, numEnd)
// 	}

// 	numStart = checkForkTask(id, taskID, latestNum, numStart, numEnd)
// 	taskID = fmt.Sprintf("[%04v|%v~%v|%v]", id, numStart, numEnd, servAddr)

// 	for {
// 		if latestNum == 0 || numStart >= latestNum {
// 			latestNum = getLatestNum(dbc)
// 			fmt.Printf("%v latestNum is [%v]\n", taskID, latestNum)
// 		}

// 		if numEnd == 0 && id == 0 { // 特殊的任务，不退出，需要读取最新块
// 			if numStart >= latestNum {
// 				time.Sleep(10 * time.Second)
// 				workingTask := workingTaskCnt()
// 				fmt.Printf("current working task:[%v]--max task:[%v]\n", workingTask, *gIntMaxWorker)
// 				if workingTask == 1 {
// 					fmt.Printf("Sync all data cost:%v\n", time.Since(ts))
// 					break
// 				}
// 			}
// 		} else {
// 			if numStart >= latestNum || numStart >= numEnd {
// 				break
// 			}
// 		}

// 		numEndNow := numStart + limit

// 		if numEndNow > numEnd && numEnd > 0 {
// 			numEndNow = numEnd
// 		}
// 		if numEnd == 0 && numEndNow > latestNum {
// 			numEndNow = latestNum
// 		}

// 		// tss := time.Now()
// 		blocks, err := client.GetBlockByLimitNext(numStart, numEndNow)
// 		_ = err
// 		blockCnt += int64(len(blocks))
// 		// fmt.Printf("%v get block:[%v ~ %v], got:[%v], err:[%v], cost:[%v]\n", taskID, numStart, numEndNow, len(blocks), err, time.Since(tss))

// 		storeBlocks(blocks)
// 		numStart += int64(len(blocks))
// 	}

// 	fmt.Printf("%v Finish work, total cost:%v, total block:%v\n", taskID, time.Since(ts), blockCnt)

// 	stopWorker()
// 	return blockCnt
// }

// func getLatestNum(dbc *grpcclient.Database) int64 {
// 	prop, err := dbc.GetDynamicProperties()
// 	if nil == err && nil != prop {
// 		return prop.LastSolidityBlockNum
// 	}
// 	return 0
// }

// func checkForkTask(id int, taskID string, latestE, b, e int64) (newB int64) {
// 	newB = b
// 	if e == 0 {
// 		if id != 0 { // e == 0 only for task id == 0
// 			return
// 		}

// 		if latestE-b > *gInt64MaxWorkload {
// 			newB = latestE - *gInt64MaxWorkload
// 			forkBlockTask(id+1, b, newB)
// 		}
// 	} else {
// 		if e-b > *gInt64MaxWorkload {
// 			newB = e - *gInt64MaxWorkload
// 			forkBlockTask(id+1, b, newB)
// 		}
// 	}
// 	return
// }

// func forkBlockTask(id int, b, e int64) {
// 	go getBlock(id, b, e)
// }
