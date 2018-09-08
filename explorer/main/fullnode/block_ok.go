package main

import (
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

var bulkFetchLimit = int64(100)

func getBlock(id int, b, e int64) {
	startWorker()

	ts := time.Now()

	servAddr := fmt.Sprintf("%s:50051", utils.GetRandFullNode())
	taskID := fmt.Sprintf("[%04v|%v~%v|%v]", id, b, e, servAddr)

	client := grpcclient.NewWallet(servAddr)
	client.Connect()
	dbc := grpcclient.NewDatabase(servAddr)
	dbc.Connect()

	le := getLatestNum(dbc)
	if le == 0 {
		stopWorker()
		getBlock(id, b, e)
		return
	}
	fmt.Printf("%v latestNum is [%v]\n", taskID, le)
	b = checkForkTask(id, "", le, b, e)

	bb := b
	cnt := int64(0)

	for {

		if e > 0 && b >= e {
			break
		}

		if id == 0 && b >= le {
			time.Sleep(3 * time.Second)

			le = getLatestNum(dbc)
			runTaskCnt := workingTaskCnt()
			fmt.Printf("Current working task:[%v]--max task:[%v]\n", runTaskCnt, *gIntMaxWorker)
			if e > 0 && 1 == runTaskCnt {
				fmt.Printf("Sync all data cost:%v\n", time.Since(ts))
				break
			}
		}

		newE := b + bulkFetchLimit

		if e > 0 && newE > e {
			newE = e
		} else if e == 0 && newE > le {
			newE = le
		}

		blocks, err := client.GetBlockByLimitNext(b, newE)
		_ = err
		cnt += int64(len(blocks))

		storeBlocks(blocks)

		c := int64(len(blocks))
		cnt += c
		b += c
	}
	fmt.Printf("%v Finish work, total cost:%v, total block:%v(%v), begin:%v, end:%v\n", taskID, time.Since(ts), cnt, b-bb, bb, b)

	stopWorker()
}

func getLatestNum(dbc *grpcclient.Database) int64 {
	prop, err := dbc.GetDynamicProperties()
	if nil == err && nil != prop {
		return prop.LastSolidityBlockNum
	}
	return 0
}

func checkForkTask(id int, taskID string, latestE, b, e int64) (newB int64) {
	newB = b
	if e == 0 {
		if id != 0 { // e == 0 only for task id == 0
			return
		}

		if latestE-b > *gInt64MaxWorkload {
			newB = latestE - *gInt64MaxWorkload
			forkBlockTask(id+1, b, newB)
		}
	} else {
		if e-b > *gInt64MaxWorkload {
			newB = e - *gInt64MaxWorkload
			forkBlockTask(id+1, b, newB)
		}
	}
	return
}

func forkBlockTask(id int, b, e int64) {
	go getBlock(id, b, e)
}
