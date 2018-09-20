package task

import (
	"time"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/log"
)

func SyncVoteWitnessRanking() {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	for {
		start := time.Now()
		log.Info("SyncVoteWitnessRanking start")
		next = next.Add(6 * time.Hour)
		log.Infof("SyncVoteWitnessRanking nextTime:%v, timestamp:%v", next, next.UnixNano()/1e6)
		t := time.NewTimer(next.Sub(now))
		<-t.C
		service.SyncVoteWitnessRanking()
		cost := time.Since(start)
		log.Infof("SyncVoteWitnessRanking end, costTime=%v", cost)
	}
}