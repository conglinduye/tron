package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

func TestLoadTrx(*testing.T) {
	initDB("tron:tron@tcp(172.16.21.224:3306)/tron")
	initRedis([]string{"127.0.0.1:6379"})
	// ts := time.Now()
	// blockIDs := genVerifyBlockIDList(0, 1000)
	// trxList := loadTransFromDB(blockIDs)
	// fmt.Printf("load %v trans cost:%v\n", len(trxList), time.Since(ts))

	// for _, trx := range trxList { // 27196, 27291
	// 	trx.ExtractContract()
	// 	// anaylzeTransaction(trx)
	// }

	// load 92498 trans cost:23.430238327s  2000000, 2010000
}

func TestRedis(*testing.T) {
	initRedis([]string{"127.0.0.1:6379"})
	fmt.Println(AddRefreshAddress([]byte("123"), []byte("345"), []byte("456")))

	fmt.Println(_redisCli.Set("123", "4123", time.Duration(0)))
}

func TestGetAccount(*testing.T) {
	initDB("tron:tron@tcp(172.16.21.224:3306)/tron")
	initRedis([]string{"127.0.0.1:6379"})

	// ts := time.Now()
	// blockIDs := genVerifyBlockIDList(0, 1000)
	// trxList := loadTransFromDB(blockIDs)
	// fmt.Printf("load %v trans cost:%v\n", len(trxList), time.Since(ts))

	// for _, trx := range trxList { // 27196, 27291
	// 	trx.ExtractContract()
	// 	// anaylzeTransaction(trx)
	// }

	// list, err := ClearRefreshAddress()
	// fmt.Println(err)
	// for idx, a := range list {
	// 	fmt.Println(idx, "-->", utils.Base58EncodeAddr([]byte(a)))
	// }

	// accList, restAddr, _ := getAccount(list)

	// fmt.Printf("accList size:%v, restAddr size:%v\n", len(accList), len(restAddr))
}

func TestRW(*testing.T) {
	client := grpcclient.GetRandomSolidity()

	utils.VerifyCall(client.GetAccount("123"))
}

func TestRedisB(*testing.T) {
	initRedis([]string{"127.0.0.1:6379"})
	for i := 0; i < 200000; i++ {
		AddRefreshAddress(i)
	}
}

func TestBlockInfo(*testing.T) {
	initDB("tron:tron@tcp(172.16.21.224:3306)/tron")
	getBlockFull(2237300)

}

func getBlockFull(num int64) {
	client := grpcclient.GetRandomWallet()

	// block, err := client.GetBlockByNum(num)
	blocks, err := client.GetBlockByLimitNext(num, num+1)

	if nil != err || nil == blocks {
		fmt.Println(err)
		return
	}

	data, _ := proto.Marshal(blocks[0])
	fmt.Printf("fullnode block:[%v](%T), hash:%v, size:%v\n", num, blocks[0], utils.HexEncode(utils.CalcBlockHash(blocks[0])), len(data))
	showBlockTrx(blocks[0])
	fmt.Printf("\n\n")
	storeBlocks(blocks[0:1])
}

func showBlockTrx(block *core.Block) {
	if nil == block {
		return
	}
	for _, trx := range block.Transactions {
		ctxOwner, _ := utils.GetContractInfoStr(trx.RawData.Contract[0])
		fmt.Printf("trx_hash:%64v\ttype:%-30v\towner_address:%v\ttimestamp:%30v\texpire:%30v\n",
			utils.HexEncode(utils.CalcTransactionHash(trx)), trx.RawData.Contract[0].Type, ctxOwner,
			utils.ConverTimestampStr(trx.RawData.Timestamp), utils.ConverTimestampStr(trx.RawData.Expiration))
	}
}
