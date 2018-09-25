package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var gIndexFile = flag.String("index", "index.idx", "file store the index info")
var gDSN = flag.String("dsn", "tron:tron@tcp(mine:3306)/tron", "msyql dsn")
var gWork = flag.String("work", "update", "work need to do, could be:\n\tupdate: update current index\n\tsearch: gen search sql\n\ttest: run test\n\tdaemon: run as daemon until receive signal")
var gTable = flag.String("table", "transactions", "table need to update index, table need has trx_hash & block_id column")
var gStep = flag.Int64("step", 1000, "step of records of each index")

func main() {
	flag.Parse()
	step = *gStep
	initDB(*gDSN)

	// searchIdxIF(getIndex())

	switch *gWork {
	case "update":
		updateIndexIF()
	case "search":
		searchIdxIF()
	case "test":
		test()
	case "daemon":
		daemon()
	default:
		flag.Usage()
	}
}

func updateIndexIF() {
	tableName := *gTable
	index := getIndex(tableName)
	updateIndex(index, tableName)
}

func test() {

	tableName := *gTable
	index := getIndex(tableName)

	printIndex(index, tableName)
	fmt.Printf("\n\n\n\n")

	reloadPos := len(index)
	if reloadPos-3 > 0 {
		reloadPos = reloadPos - 3
	}

	updateIndex(index[:reloadPos], tableName)
}

func daemon() {
	signalHandle()
	tableName := *gTable
	var index []*TrxIndex
	var ok bool
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(30 * time.Second)

	updateLoop:
		for {
			if index, ok = updateIndex(index, tableName); ok {
				storeIdxToDB(index, tableName)
			}

			select {
			case <-ticker.C:
				continue
			case <-quit:
				break updateLoop
			}
		}
		ticker.Stop()
	}()

	<-quit

	fmt.Printf("Daemon Quit\n")
	wg.Wait()
}

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
