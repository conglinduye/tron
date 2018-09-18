package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var gIntMaxWorker = flag.Int("worker", 10, "maximum worker for fetch blocks")
var gStrMysqlDSN = flag.String("dsn", "tron:tron@tcp(172.16.21.224:3306)/tron", "mysql connection string(DSN)")
var gInt64MaxWorkload = flag.Int64("workload", 10000, "maximum workload for each worker")
var gStartBlokcID = flag.Int64("start_block", 0, "block num start to synchronize")
var gEndBlokcID = flag.Int64("end_block", 0, "block num end to synchronize, default 0 means run as daemon")
var gRedisDSN = flag.String("redisDSN", "127.0.0.1:6379", "redis DSN")
var gMaxErrCntPerNode = flag.Int("max_err_per_node", 10, "max error before we try to other node")
var gMaxAccountWorkload = flag.Int("max_account_workload", 200, "max account a node need handle not fork new worker")
var gIntHandleAccountInterval = flag.Int("account_handle_interval", 30, "account info synchronize handle minmum interval in seconds")

var quit = make(chan struct{}) // quit signal channel
var wg sync.WaitGroup

func signalHandle() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if !needQuit() {
			close(quit)
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

func startDaemon() {
	startAssetDaemon()
	startWintnessDaemon()
	startRedisAccountRefreshPush()
	startAccountDaemon()
	startNodeDaemon()
}

func main() {
	flag.Parse()

	maxErrCnt = *gMaxErrCntPerNode
	getAccountWorkerLimit = *gMaxAccountWorkload

	signalHandle()

	initDB(*gStrMysqlDSN)
	initRedis([]string{*gRedisDSN})
	startDaemon()

	getAllBlocks()
	if !needQuit() {
		close(quit)
	}

	syncAccount() // syn account after getAllBlocks() quit

	fmt.Println("Wait other daemon quit .......")
	wg.Wait()

	fmt.Println("fullnode QUIT")
}

func getAllBlocks() {
	wc1 = newWorkerCounter(*gIntMaxWorker)
	ts := time.Now()
	getBlock(0, *gStartBlokcID, *gEndBlokcID)
	fmt.Printf("get all blocks cost:%v\n", time.Since(ts))
}
