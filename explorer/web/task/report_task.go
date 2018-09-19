package task

import (
	"time"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/log"
)

func SyncCacheTodayReport() {
	for range time.Tick(3 * time.Minute) {
		log.Info("SyncCacheTodayReportTask start")
		start := time.Now()
		service.SyncCacheTodayReport()
		cost := time.Since(start)
		log.Infof("SyncCacheTodayReportTask end, costTime=%v", cost)
	}
}

func SyncPersistYesterdayReport() {
	for {
		start := time.Now()
		log.Info("SyncPersistYesterdayReport start")
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		next = next.Add(24 * time.Hour).Add(1 * time.Second)
		log.Infof("SyncPersistYesterdayReport nextTime:%v, timestamp:%v", next, next.UnixNano() / 1e6)
		t := time.NewTimer(next.Sub(now))
		<-t.C
		service.SyncPersistYesterdayReport()
		cost := time.Since(start)
		log.Infof("SyncPersistYesterdayReport end, costTime=%v", cost)
	}
}