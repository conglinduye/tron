package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
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
	blockIDs := genVerifyBlockIDList(20, 1000)
	trxList := loadTransFromDB(blockIDs)
	fmt.Printf("load %v trans cost:%v\n", len(trxList), time.Since(ts))

	for _, trx := range trxList { // 27196, 27291
		trx.ExtractContract()
		anaylzeTransaction(trx)
	}

	list, err := ClearRefreshAddress()
	fmt.Println(err)

	accList, restAddr, _ := getAccount(list)
	storeAccount(accList)

	fmt.Printf("accList size:%v, restAddr size:%v\n", len(accList), len(restAddr))
}

func TestRW(*testing.T) {

	initDB("tron:tron@tcp(172.16.21.224:3306)/tron")
	client := grpcclient.GetRandomSolidity()

	acc, _ := client.GetAccount("TKoU7MkprWw8q142Sd199XU4B5fUaMVBNm")

	accc := new(account)
	accc.SetRaw(acc)

	storeAccount([]*account{accc})

	for {

		fmt.Printf("\n\n--%v--\n", getDBMaxBlockID())
		time.Sleep(3 * time.Second)
	}

}
