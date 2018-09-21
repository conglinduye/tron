package main

import (
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/task"
)

func Asyncload() {
	//init buffer
	buffer.GetBlockBuffer()
	buffer.GetWitnessBuffer()
	buffer.GetMarketBuffer()
	buffer.GetVoteBuffer()
	buffer.GetAccountTokenBuffer()
	buffer.GetTokenBuffer()

	go task.SyncCacheTodayReport()

	go task.SyncPersistYesterdayReport()

	go task.SyncAssetIssueParticipated()

	go task.SyncVoteWitnessRanking()
}
