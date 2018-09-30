package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/main/module/account"
	"github.com/wlcy/tron/explorer/main/module/rawmysql"
)

var gMysqlDSN = flag.String("mysql", "tron:tron@tcp(mysql1:3306)/tron", "mysql db connection string")
var gRedisDSN = flag.String("redis", "redisServer:6379", "redis server connection string")
var gMaxWorker = flag.Int("worker", 30, "max worker of connectting to main net to get account info")
var gMaxDBWorker = flag.Int("dbworker", 2, "max worker of writing data to db")

func main() {
	flag.Parse()
	initDB()
	// initRedis()

	signalHandle()

	log.Str2Level("info")

	addrList := getAssetIssueAccount()

	startWork(addrList)

}

func initDB() {
	rawmysql.DSN = *gMysqlDSN
}

// func initRedis() {
// 	rawredis.DSN = *gRedisDSN
// }

func getAssetIssueAccount() (ret []*account.AddressSyncInfo) {
	dbc := rawmysql.GetMysqlDB()

	strSQL := "select a.owner_address, ifnull(b.latest_operation_time,0) from asset_issue a left join tron_account b on b.address = a.owner_address"
	rows, err := dbc.Query(strSQL)
	if nil != err || nil == rows {
		log.Errorf("get asset_issue owner address failed:%v, sql:%v", err, strSQL)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ownerAddr string
		var latestOpTime int64
		err := rows.Scan(&ownerAddr, &latestOpTime)
		if nil != err || 0 == len(ownerAddr) {
			log.Errorf("scan owner address failed:%v", err)
			continue
		}
		ret = append(ret, &account.AddressSyncInfo{Addr: ownerAddr, LOT: latestOpTime})
	}
	return
}

func startWork(addrList []*account.AddressSyncInfo) {
	accWorker := account.NewAccountWorker(*gMaxWorker, 10240, *gMaxDBWorker)
	accWorker.StartAccountWorker()
	accWorker.StartDBWorker()
	accWorker.AppendTask(addrList)

	for {
		fmt.Println("### wait signal msg......")
		<-sigMsg
		fmt.Println("### got signal msg......")
		if needQuit() {
			accWorker.WaitStop()
			fmt.Println("### quit ......")
			break
		}
		accWorker.AppendTask(addrList)
	}
	fmt.Println("### call WaitStop ......")
	accWorker.WaitStop()
}

var quit = make(chan struct{})
var sigMsg = make(chan struct{})

func signalHandle() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cnt := 3
	go func() {
		for {
			<-sigs
			cnt--
			sigMsg <- struct{}{}
			fmt.Println("receive signal and send sig msg")
			if cnt == 0 {
				if !needQuit() {
					close(quit)
				}
			}
		}
	}()
}

func needQuit() bool {
	select {
	case <-quit:
		return true
	default:
		return false
	}
}
