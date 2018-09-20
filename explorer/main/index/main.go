package main

import (
	"flag"
	"fmt"
)

var gIndexFile = flag.String("index", "index.idx", "file store the index info")
var gDSN = flag.String("dsn", "tron:tron@tcp(mine:3306)/tron", "msyql dsn")
var gWork = flag.String("work", "update", "work need to do\n\tupdate: update current index\n\tsearch: gen search sql\n\ttest: run test")

func main() {
	flag.Parse()
	initDB(*gDSN)

	// searchIdxIF(getIndex())

	switch *gWork {
	case "update":
		updateIndexIF()
	case "search":
		searchIdxIF(getIndex())
	case "test":
		test()
	case "daemon":
		daemon()
	default:
		flag.Usage()
	}
}

func updateIndexIF() {
	index := getIndex()
	updateIndex(index)
}

func test() {

	index := getIndex()

	printIndex(index)
	fmt.Printf("\n\n\n\n")

	reloadPos := len(index)
	if reloadPos-3 > 0 {
		reloadPos = reloadPos - 3
	}

	updateIndex(index[:reloadPos])
}

func daemon() {
	index := getIndex()
	for {
		index = updateIndex(index)
		storeIdxToDB(index)

		// time.Sleep(30 * time.Second)
		return
	}
}
