package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"time"
	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/log"
	"encoding/json"
	"fmt"
)

const BaseOverviewKey = "org.tron.explorer.report.base.overview"

const TotalOverviewKey = "org.tron.explorer.report.total.overview"

func QueryReport() (*entity.ReportResp, error) {
	reportResp := &entity.ReportResp{}
	reportOverviewsList := make([]*entity.ReportOverview, 0)

	totalOverviewValue, err := config.RedisCli.Get(TotalOverviewKey).Result()
	if err != nil  {
		log.Errorf("syncReportJob get err:[%v]", err)
	}

	if totalOverviewValue == "" {
		baseOverviewValue, _ := config.RedisCli.Get(BaseOverviewKey).Result()
		if baseOverviewValue == "" {
			SyncReportJob()
			baseOverviewValue, _ = config.RedisCli.Get(BaseOverviewKey).Result()
		}
		if err := json.Unmarshal([]byte(baseOverviewValue), &reportOverviewsList); err == nil {
			reportOverview := &entity.ReportOverview{}
			t := time.Now()
			t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
			t1 = t1.Add(-8 * time.Hour)
			t2 := t
			t3 := t1.Add(-24 * time.Hour)
			startTime := t1.UnixNano() / 1e6
			endTime := t2.UnixNano() / 1e6
			dateTime := t3.UnixNano() / 1e6
			fmt.Printf("startTime:%v, endTime:%v, dateTime:%v\n", startTime, endTime, dateTime)
			last := reportOverviewsList[len(reportOverviewsList) - 1]
			if dateTime == last.Date {
				syncReportBetweenTime(startTime, endTime, reportOverview)
				syncReportByTime(endTime, reportOverview)
				reportOverview.Date = startTime

				reportOverviewsList = reportOverviewsList[:len(reportOverviewsList)-1]
				reportOverviewsList = append(reportOverviewsList, reportOverview)
				value, err := json.Marshal(reportOverviewsList)
				if err != nil {
					log.Errorf("syncReportJob json marshal err:[%v]", err)
				}
				err = config.RedisCli.Set(TotalOverviewKey, string(value), 1 * time.Minute).Err()
				if err != nil {
					log.Errorf("syncReportJob set err:[%v]", err)
				}
			} else {
				SyncReportJob()
			}
		}

	} else {
		if err := json.Unmarshal([]byte(totalOverviewValue), reportOverviewsList); err != nil {
			log.Errorf("syncReportJob json unmarshal err:[%v]", err)
		}
	}

	reportResp.Success = true
	reportResp.Data = reportOverviewsList
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

// SyncReportJob
func SyncReportJob() {
	baseOverviewValue, err := config.RedisCli.Get(BaseOverviewKey).Result()
	if err != nil {
		log.Errorf("syncReportJob get err:[%v]", err)
	}
	if baseOverviewValue == "" {
		reportOverviewsList := make([]*entity.ReportOverview, 0)
		for i := 13; i >= 1; i-- {
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

			syncReportBetweenTime(startTime, endTime, reportOverview)
			syncReportByTime(endTime, reportOverview)
			reportOverview.Date = startTime
			reportOverviewsList = append(reportOverviewsList, reportOverview)
		}
		fmt.Printf("reportOverviewsList size: %v\n", len(reportOverviewsList))
		value, err := json.Marshal(reportOverviewsList)
		if err != nil {
			log.Errorf("syncReportJob json marshal err:[%v]", err)
		}
		err = config.RedisCli.Set(BaseOverviewKey, string(value), 0).Err()
		if err != nil {
			log.Errorf("syncReportJob set err:[%v]", err)
		}

	} else {
		reportOverviewsList := make([]*entity.ReportOverview, 0)
		if err := json.Unmarshal([]byte(baseOverviewValue), &reportOverviewsList); err == nil {
			reportOverview := &entity.ReportOverview{}
			t := time.Now()
			t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
			t1 = t1.Add(-24 * time.Hour)
			t1 = t1.Add(-8 * time.Hour)
			t2 := t1.Add(24 * time.Hour)
			startTime := t1.UnixNano() / 1e6
			endTime := t2.UnixNano() / 1e6
			fmt.Printf("startTime:%v, endTime:%v\n", startTime, endTime)
			last := reportOverviewsList[len(reportOverviewsList) - 1]
			if startTime != last.Date {
				syncReportBetweenTime(startTime, endTime, reportOverview)
				syncReportByTime(endTime, reportOverview)
				reportOverview.Date = startTime
				reportOverviewsList = append(reportOverviewsList, reportOverview)
				reportOverviewsList = reportOverviewsList[1:]
				value, err := json.Marshal(reportOverviewsList)
				if err != nil {
					log.Errorf("syncReportJob json marshal err:[%v]", err)
				}
				err = config.RedisCli.Set(BaseOverviewKey, string(value), 0).Err()
				if err != nil {
					log.Errorf("syncReportJob set err:[%v]", err)
				}
			}
		}

	}

}
