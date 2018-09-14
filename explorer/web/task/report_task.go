package task

import (
	"time"
	"github.com/wlcy/tron/explorer/web/service"
)

func SyncCacheTodayReportTask() {
	for range time.Tick(3 * time.Minute) {
		service.SyncCacheTodayReport()
	}
}

func SyncCacheHistoryReportTask() {
	for {
		now := time.Now()
		next := now.Add(24 * time.Hour).Add(1 * time.Second)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		service.SyncCacheHistoryReport()
	}
}