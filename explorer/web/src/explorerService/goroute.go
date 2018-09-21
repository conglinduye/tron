package main

import (
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/task"
)

func Async() {
	//init buffer
	go buffer.GetBlockBuffer()
	go buffer.GetWitnessBuffer()
	go buffer.GetMarketBuffer()
	go buffer.GetAccountTokenBuffer()
	go buffer.GetTokenBuffer()

	go task.SyncCacheTodayReport()
	go task.SyncPersistYesterdayReport()
	go task.SyncAssetIssueParticipated()
	go task.SyncVoteWitnessRanking()
}
