package task

import (
	"time"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/log"
)

func SyncAssetIssueParticipated() {
	for range time.Tick(1 * time.Hour) {
		start := time.Now()
		log.Info("SyncAssetIssueParticipated start")
		service.SyncAssetIssueParticipated()
		cost := time.Since(start)
		log.Infof("SyncAssetIssueParticipated end, costTime=%v", cost)
	}
}

