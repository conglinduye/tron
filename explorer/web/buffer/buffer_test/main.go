package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/buffer"
)

func main() {
	flag.Parse()
	// mysql.Initialize("127.0.0.1", "3306", "tron", "budev", "tron**1")
	mysql.Initialize("mine", "3306", "tron", "tron", "tron")
	log.ChangeLogLevel(log.Str2Level("INFO"))

	// initRedis([]string{"127.0.0.1:6379"})
	bb := buffer.GetBlockBuffer()
	_ = bb
	cnt := 0
	for cnt < 10 {
		// getBlocks(-1, 0, 100) // get latest 100 blocks
		// getBlocks(-1, 500, 100) // get

		// getTrx()
		// getTrxs()

		// fmt.Printf("\n### %v, %v, %v, %v\n\n", bb.GetMaxBlockID(), bb.GetMaxConfirmedBlockID(), bb.GetSolidityNodeMaxBlockID(), bb.GetFullNodeMaxBlockID())

		time.Sleep(10 * time.Second)
		cnt++
	}
}

func getTrx() {
	// bb := buffer.GetBlockBuffer()

	// ret := bb.GetTransactions(0, 100)
	// fmt.Printf("trx 0~100, len:%v\n", len(ret))

	getMaxBlockIDTrx()
}

func getMaxBlockIDTrx() {
	bb := buffer.GetBlockBuffer()

	maxBlockID := bb.GetMaxBlockID()
	block := bb.GetBlock(maxBlockID)
	trxs := bb.GetTransactionByBlockID(maxBlockID)
	if 0 == len(trxs) {
		fmt.Printf("get max blockID (%v) %v\n\ttrxs empty!\n\n", maxBlockID, utils.ToJSONStr(block))
	} else {
		fmt.Printf("get max blockID (%v) %v\n\ttrxs count:%v\n%v\n\n", maxBlockID, utils.ToJSONStr(block), len(trxs), utils.ToJSONStr(trxs))
	}
}

func getTrxs() {
	bb := buffer.GetBlockBuffer()

	var start, count, trxLen int64
	start = 100
	count = 40
	ts := time.Now()
	trxs := bb.GetTransactions(start, count)
	tsc := time.Since(ts)
	trxLen = int64(len(trxs))
	fmt.Printf("get trxs result, ret count %v, req count %v, cost:%v\n", trxLen, count, tsc)

	if trxLen > 0 {
		var maxBlockID, minBlockID, cur, n int64
		maxBlockID = trxs[0].Block
		minBlockID = trxs[trxLen-1].Block
		cur = maxBlockID
		n++
		for idx, trx := range trxs {
			_ = idx
			// fmt.Printf("%v-->%v\n", idx, trx.Block)
			if cur == trx.Block || cur > trx.Block {
			} else {
				fmt.Printf("trx list block error, Non-continuous number %v--%v\n", cur, trx.Block)
			}
			if cur != trx.Block {
				n++
			}
			cur = trx.Block
		}

		fmt.Printf("maxBlockID:%v, minBlockID:%v, len:%v, block count:%v\n", maxBlockID, minBlockID, trxLen, n)

	}

}

func getBlocks(start, rs, re int64) {
	tsr := time.Now()
	bb := buffer.GetBlockBuffer()

	ret, _ := bb.GetBlocks(start, rs, re)
	retLen := len(ret)
	var c, unc int
	var minCBlockID int64 = 900000000
	var maxCBlockID int64
	var maxUncBlockID int64
	var minUncBlockID int64 = 9000000000
	for _, block := range ret {
		if block.Confirmed {
			c++
			if maxCBlockID < block.Number {
				maxCBlockID = block.Number
			}
			if minCBlockID > block.Number {
				minCBlockID = block.Number
			}
		} else {
			unc++
			if maxUncBlockID < block.Number {
				maxUncBlockID = block.Number
			}
			if minUncBlockID > block.Number {
				minUncBlockID = block.Number
			}
		}
	}
	fmt.Printf("(min, max) confirmed block id:(%v,%v) count:%v;  (min, max) unconfirmed block id:(%v,%v) count:%v\n", minCBlockID, maxCBlockID, maxCBlockID-minCBlockID+1, minUncBlockID, maxUncBlockID, maxUncBlockID-minUncBlockID+1)

	if retLen == 0 || retLen != int(re) {
		fmt.Printf("Get Block failed! start, re, rs = (%v, %v, %v) ret list len:%v\n", start, rs, re, retLen)
	} else {
		fmt.Printf("\nload from buffer %v~ %v (%v), size:%v, ret[0].num:%v, ret[%v].num:%v, cost:%v\n\n", rs, re, re, len(ret), ret[0].Number, retLen, ret[retLen-1].Number, time.Since(tsr))
	}

	fmt.Printf("\n\n")

}
