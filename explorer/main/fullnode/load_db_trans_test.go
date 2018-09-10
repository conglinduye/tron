package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
)

func TestLoadTrx(*testing.T) {
	initDB("tron:tron@tcp(172.16.21.224:3306)/tron")
	initRedis([]string{"127.0.0.1:6379"})
	ts := time.Now()
	blockIDs := genVerifyBlockIDList(0, 1000)
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("load %v trans cost:%v\n", len(trxList), time.Since(ts))

	for _, trx := range trxList { // 27196, 27291
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}

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

	ts := time.Now()
	blockIDs := genVerifyBlockIDList(0, 1000)
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("load %v trans cost:%v\n", len(trxList), time.Since(ts))

	for _, trx := range trxList { // 27196, 27291
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}

	list, err := ClearRefreshAddress()
	fmt.Println(err)
	for idx, a := range list {
		fmt.Println(idx, "-->", utils.Base58EncodeAddr([]byte(a)))
	}

	accList, restAddr, _ := getAccount(list)

	fmt.Printf("accList size:%v, restAddr size:%v\n", len(accList), len(restAddr))
}

func TestRW(*testing.T) {
	client := grpcclient.GetRandomSolidity()

	utils.VerifyCall(client.GetAccount("123"))
}
