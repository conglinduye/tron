package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"time"
)

func QueryReport() (*entity.ReportResp, error) {
	reportResp := &entity.ReportResp{}
	reportOverviews := make([]*entity.ReportOverview, 0)
	for i := 14; i >= 1; i-- {
		reportOverview := &entity.ReportOverview{}
		t := time.Now()
		t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
		for j := 1; j <= i; j++ {
			t1 = t1.Add(-24 * time.Hour)
		}
		t1 = t1.Add(-8 * time.Hour)
		t2 := t1.Add(24 * time.Hour)
		startTime := t1.UnixNano() / 1e6
		endTime := t2.UnixNano() / 1e6
		//fmt.Printf("t1:%v, t2:%v\n", t1, t2)
		//fmt.Printf("startTime:%v, endTime:%v\n", startTime, endTime)

		syncReportBetweenTime(startTime, endTime, reportOverview)
		syncReportByTime(endTime, reportOverview)
		reportOverview.Date = startTime
		reportOverviews = append(reportOverviews, reportOverview)
	}
	reportResp.Success = true
	reportResp.Data = reportOverviews

	return reportResp, nil
}

// syncReportBetweenTime
func syncReportBetweenTime(startTime, endTime int64, overview *entity.ReportOverview) {
	reportBlock, _ := module.QueryReportBlock(startTime, endTime)
	if reportBlock.TotalCount == 0 {
		overview.AvgBlockSize = 0
	} else {
		overview.AvgBlockSize = reportBlock.TotalSize / reportBlock.TotalCount
	}
	overview.AvgBlockTime = 3

	totalAccount, _ := module.QueryReportAccount(startTime, endTime)
	overview.NewAddressSeen =totalAccount

}

// syncReportByTime
func syncReportByTime(dateTime int64, overview *entity.ReportOverview) {
	reportBlock, _ := module.QueryTotalReportBlock(dateTime)
	totalTransaction, _ := module.QueryTotalReportTransaction(dateTime)
	totalAccount, _ := module.QueryTotalReportAccount(dateTime)

	overview.TotalBlockCount = reportBlock.TotalCount
	overview.BlockchainSize = reportBlock.TotalSize
	overview.TotalTransaction = totalTransaction
	overview.TotalAddress = totalAccount
}
