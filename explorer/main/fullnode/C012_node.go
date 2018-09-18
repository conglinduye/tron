package main

import (
	"fmt"
	"time"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/wlcy/tron/explorer/core/grpcclient"
)

func startNodeDaemon() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if nodes, ok := getNodes(); ok {
				icnt, ucnt, ecnt, err := storeNodes(nodes)
				fmt.Printf("node daemon work result:(%v, %v, %v, %v)\n", icnt, ucnt, ecnt, err)
				time.Sleep(30 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}

			if needQuit() {
				break
			}
		}
		fmt.Printf("Node Daemon QUIT\n")
	}()
}

func getNodes() ([]*api.Node, bool) {
	client := grpcclient.GetRandomWallet()

	nodeList, err := client.ListNodes()
	if nil != err || len(nodeList) == 0 {
		return nil, false
	}

	return nodeList, true
}

func storeNodes(nodeList []*api.Node) (iCnt int64, uCnt int64, eCnt int64, err error) {
	if len(nodeList) == 0 {
		return
	}

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		fmt.Printf("get db failed:%v\n", err)
		return
	}
	/*
		CREATE TABLE `nodes` (
		  `node_host` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'host',
		  `node_port` int NOT NULL DEFAULT '0' COMMENT 'port',
		  `create_time` timestamp not null default current_timestamp,
		  PRIMARY KEY (`node_host`,`node_port`),
		  KEY `idx_node_host` (`node_host`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	/**/
	sqlI := `insert into nodes ( node_host, node_port ) values  (?, ?)`
	stmt, err := txn.Prepare(sqlI)
	if nil != err {
		fmt.Printf("prepare [%v] failed:%v\n", sqlI, err)
		return
	}
	defer stmt.Close()

	txn.Exec("delete from nodes where 1=1")
	for _, node := range nodeList {
		_, err = stmt.Exec(string(node.Address.Host), node.Address.Port)
		if nil != err {
			// fmt.Printf("insert asset_issue [%T] failed:%v\n", asset, err)
			eCnt++
		} else {
			iCnt++
		}
	}

	err = txn.Commit()
	if err != nil {
		// fmt.Printf("connit witness failed:%v\n", err)
		return
	}
	return
}
