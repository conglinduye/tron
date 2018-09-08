package main

import (
	"flag"
	"time"

	"fmt"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"

	_ "github.com/go-sql-driver/mysql"
)

var gIntMaxWorker = flag.Int("worker", 10, "maximum worker for fetch blocks")
var gStrMysqlDSN = flag.String("dsn", "tron:tron@tcp(172.16.21.224:3306)/tron", "mysql connection string(DSN)")
var gInt64MaxWorkload = flag.Int64("workload", 10000, "maximum workload for each worker")

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
	blockCnt := getBlock(0, 0, 0)
	fmt.Printf("total get block count:%v, cost:%v\n", blockCnt, time.Since(ts))
}

// var maxWorkload = int64(10000) // getBlock任务最大工作负荷

var maxWorker chan struct{}

func getBlock(id int, numStart int64, numEnd int64) int64 {
	startWorker()

	ts := time.Now()
	servAddr := fmt.Sprintf("%s:50051", utils.GetRandFullNode())

	taskID := fmt.Sprintf("[%04v|%v~%v|%v]", id, numStart, numEnd, servAddr)

	client := grpcclient.NewWallet(servAddr)
	client.Connect()
	dbc := grpcclient.NewDatabase(servAddr)
	dbc.Connect()

	blockCnt := int64(0)
	latestNum := int64(0)
	limit := int64(100)

	if latestNum == 0 || numStart >= latestNum {
		prop, err := dbc.GetDynamicProperties()
		if nil == err && nil != prop {
			latestNum = prop.LastSolidityBlockNum
		}
		fmt.Printf("%v latestNum is [%v]\n", taskID, latestNum)
	}

	if latestNum == 0 {
		stopWorker()
		return getBlock(id, numStart, numEnd)
	}

	numStart = checkForkTask(id, taskID, latestNum, numStart, numEnd)
	taskID = fmt.Sprintf("[%04v|%v~%v|%v]", id, numStart, numEnd, servAddr)

	for {
		if latestNum == 0 || numStart >= latestNum {
			prop, err := dbc.GetDynamicProperties()
			if nil == err && nil != prop {
				latestNum = prop.LastSolidityBlockNum
			}
			fmt.Printf("%v latestNum is [%v]\n", taskID, latestNum)
		}

		if numEnd == 0 && id == 0 { // 特殊的任务，不退出，需要读取最新块
			if numStart >= latestNum {
				time.Sleep(10 * time.Second)
				workingTask := workingTaskCnt()
				fmt.Printf("current working task:[%v]--max task:[%v]\n", workingTask, *gIntMaxWorker)
				if workingTask == 1 {
					fmt.Printf("Sync all data cost:%v\n", time.Since(ts))
					break
				}
			}
		} else {
			if numStart >= latestNum || numStart >= numEnd {
				break
			}
		}

		numEndNow := numStart + limit

		if numEndNow > numEnd && numEnd > 0 {
			numEndNow = numEnd
		}

		// tss := time.Now()
		blocks, err := client.GetBlockByLimitNext(numStart, numEndNow)
		_ = err
		blockCnt += int64(len(blocks))
		// fmt.Printf("%v get block:[%v ~ %v], got:[%v], err:[%v], cost:[%v]\n", taskID, numStart, numEndNow, len(blocks), err, time.Since(tss))

		storeBlocks(blocks)
		numStart += int64(len(blocks))
	}

	fmt.Printf("%v Finish work, total cost:%v, total block:%v\n", taskID, time.Since(ts), blockCnt)

	stopWorker()
	return blockCnt
}

func checkForkTask(id int, taskID string, latestNum, numStart, numEnd int64) int64 {
	if latestNum == 0 {
		return numStart
	}

	newStart := numStart

	if (numEnd == 0 || latestNum < numEnd) && latestNum-*gInt64MaxWorkload > numStart {
		newStart = latestNum - *gInt64MaxWorkload
	} else if latestNum >= numEnd && numEnd-numStart > *gInt64MaxWorkload {
		newStart = numEnd - *gInt64MaxWorkload
	}

	if newStart != numStart { // fork sub
		fmt.Printf("%v fork task range: %v ~ %v\n", taskID, numStart, newStart)
		go getBlock(id+1, numStart, newStart)
	}

	return newStart
}
