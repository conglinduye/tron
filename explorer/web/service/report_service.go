package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"time"
	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/log"
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v4"
	"strconv"
)

const HistoryOverviewKey = "org.tron.explorer.report.history.overview"

const TodayOverviewKey = "org.tron.explorer.report.today.overview"

const YesterdayOverviewKey = "org.tron.explorer.report.yesterday.overview"

func QueryReport() (*entity.ReportResp, error) {
	reportResp := &entity.ReportResp{}

	historyOverviewValue, err := config.RedisCli.Get(HistoryOverviewKey).Result()
	if err == redis.Nil {
		SyncCacheHistoryReport()
		historyOverviewValue, _ = config.RedisCli.Get(HistoryOverviewKey).Result()
	} else if err != nil {
		log.Errorf("QueryReport historyOverviewValue redis get value error :[%v]\n", err)
		return nil, err
	}

	if historyOverviewValue == "" || historyOverviewValue == "[]" {
		SyncCacheHistoryReport()
		historyOverviewValue, _ = config.RedisCli.Get(HistoryOverviewKey).Result()
	}
	totalReportOverviews := make([]*entity.ReportOverview, 0)
	err = json.Unmarshal([]byte(historyOverviewValue), &totalReportOverviews)
	if err != nil {
		log.Errorf("QueryReport json.Unmarshal totalReportOverviews error :[%v]\n", err)
		return nil, err
	}

	if len(totalReportOverviews) < 1 {
		log.Error("QueryReport len(totalReportOverviews) < 1\n")
		return nil, err
	}

	last := totalReportOverviews[len(totalReportOverviews) - 1]
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC).Add(-24 * time.Hour)
	dateTime := t1.UnixNano() / 1e6
	if dateTime != last.Date {
		SyncPersistYesterdayReport()
	}

	todayOverviewValue, err := config.RedisCli.Get(TodayOverviewKey).Result()
	if err == redis.Nil {
		SyncCacheTodayReport()
		todayOverviewValue, _ = config.RedisCli.Get(TodayOverviewKey).Result()
	} else if err != nil {
		log.Errorf("QueryReport todayOverviewValue redis get value error :[%v]\n", err)
		return nil, err
	}
	if todayOverviewValue == "" {
		SyncCacheTodayReport()
		todayOverviewValue, _ = config.RedisCli.Get(TodayOverviewKey).Result()
	}
	todayReportOverview := &entity.ReportOverview{}
	err = json.Unmarshal([]byte(todayOverviewValue), &todayReportOverview)
	if err != nil {
		log.Errorf("QueryReport json.Unmarshal todayReportOverview error :[%v]\n", err)
		return nil, err
	}

	totalReportOverviews = append(totalReportOverviews, todayReportOverview)

	reportResp.Success = true
	reportResp.Data = totalReportOverviews
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
	overview.NewBlockSeen = reportBlock.TotalCount
	overview.NewTransactionSeen = reportBlock.TotalTransaction

	totalAccount, _ := module.QueryReportAccount(startTime, endTime)
	overview.NewAddressSeen =totalAccount

}

// syncReportByTime
func syncReportByTime(dateTime int64, overview *entity.ReportOverview) {
	reportBlock, _ := module.QueryTotalReportBlock(dateTime)
	totalAccount, _ := module.QueryTotalReportAccount(dateTime)

	overview.TotalBlockCount = reportBlock.TotalCount
	overview.TotalTransaction = reportBlock.TotalTransaction
	overview.BlockchainSize = reportBlock.TotalSize
	overview.TotalAddress = totalAccount
}

func SyncInitReport() {
	count, _ := module.QueryTotalStatistics()
	if count <= 0 {
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, time.UTC)

		t,_ := time.Parse("20060102150405", "20180625000000")
		t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
		t2 := t1.Add(24 * time.Hour)
		nowTime := now.UnixNano() / 1e6
		startTime := t1.UnixNano() / 1e6
		endTime := t2.UnixNano() / 1e6
		fmt.Printf("nowTime:%v, startTime:%v, endTime:%v\n",nowTime, startTime, endTime)
		for ; startTime < nowTime ; {
			reportOverview := &entity.ReportOverview{}
			syncReportBetweenTime(startTime, endTime, reportOverview)
			syncReportByTime(endTime, reportOverview)
			reportOverview.Date = startTime
			module.InsertStatistics(reportOverview)
			t1 = t1.Add(24 * time.Hour)
			t2 = t1.Add(24 * time.Hour)
			startTime = t1.UnixNano() / 1e6
			endTime = t2.UnixNano() / 1e6
			fmt.Printf("startTime:%v, endTime:%v\n", startTime, endTime)
		}
	}

	SyncCacheHistoryReport()

	SyncCacheTodayReport()
}

func SyncPersistYesterdayReport() {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
	t3 := t1.Add(-24 * time.Hour)
	dateTime := t3.UnixNano() / 1e6

	key := YesterdayOverviewKey + "." + strconv.FormatInt(dateTime,10)
	if !lockYesterdayOverview(key) {
		log.Infof("SyncPersistYesterdayReport lock exit, key:%v", key)
		return
	}
	err := config.RedisCli.Set(key, "SyncPersistYesterdayReport", 60 * time.Second).Err()
	if err != nil {
		log.Errorf("SyncPersistYesterdayReport lock redis set err:[%v]", err)
	}

	strSQL := fmt.Sprintf(`
			select date, avg_block_time, avg_block_size, new_block_seen, new_transaction_seen, 
			new_address_seen, total_block_count, total_transaction, total_address, blockchain_size
			from wlcy_statistics order by date desc limit 1`)
	reportOverviews, _:= module.QueryStatistics(strSQL)
	if len(reportOverviews) == 0 {
		log.Error("SyncPersistYesterdayReport QueryStatistics len(reportOverviews) == 0")
		return
	}
	if reportOverviews[0].Date < dateTime {
		t1 = t1.Add(-24 * time.Hour)
		t2 := t1.Add(24 * time.Hour)
		startTime := t1.UnixNano() / 1e6
		endTime := t2.UnixNano() / 1e6
		fmt.Printf("SyncPersistYesterdayReport startTime:%v, endTime:%v\n", startTime, endTime)
		reportOverview := &entity.ReportOverview{}
		syncReportBetweenTime(startTime, endTime, reportOverview)
		syncReportByTime(endTime, reportOverview)
		reportOverview.Date = startTime
		err := module.InsertStatistics(reportOverview)
		if err != nil {
			log.Errorf("InsertStatistics error :[%v]\n", err)
			return
		}
		historyOverviewValue, err := config.RedisCli.Get(HistoryOverviewKey).Result()
		if err == redis.Nil {
			SyncCacheHistoryReport()
			historyOverviewValue, _ = config.RedisCli.Get(HistoryOverviewKey).Result()
		} else if err != nil {
			log.Errorf("SyncPersistYesterdayReport historyOverviewValue redis get value error :[%v]\n", err)
			return
		}
		if historyOverviewValue == "" {
			SyncCacheHistoryReport()
			historyOverviewValue, _ = config.RedisCli.Get(HistoryOverviewKey).Result()
		}
		reportOverviews := make([]*entity.ReportOverview, 0)
		json.Unmarshal([]byte(historyOverviewValue), &reportOverviews)
		reportOverviews = append(reportOverviews, reportOverview)
		value, _ := json.Marshal(reportOverviews)
		err = config.RedisCli.Set(HistoryOverviewKey, string(value), 0).Err()
		if err != nil {
			log.Errorf("SyncPersistYesterdayReport redis set err:[%v]", err)
		}
	}
	log.Info("SyncPersistYesterdayReport handle done")

}

func SyncCacheHistoryReport() {
	log.Info("SyncCacheHistoryReport handle start")
	strSQL := fmt.Sprintf(`
			select date, avg_block_time, avg_block_size, new_block_seen, new_transaction_seen, 
			new_address_seen, total_block_count, total_transaction, total_address, blockchain_size
			from wlcy_statistics order by date asc `)
	reportOverviews, err := module.QueryStatistics(strSQL)
	if err != nil {
		log.Errorf("SyncCacheHistoryReport set err:[%v]", err)
	}
	value, err := json.Marshal(reportOverviews)
	if err != nil {
		log.Errorf("SyncCacheHistoryReport redis set err:[%v]", err)
	}
	err = config.RedisCli.Set(HistoryOverviewKey, string(value), 0).Err()
	if err != nil {
		log.Errorf("SyncCacheHistoryReport set err:[%v]", err)
	}
	log.Info("SyncCacheHistoryReport handle done")
}

func SyncCacheTodayReport() {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0,0,0,0, time.UTC)
	t2 := t
	startTime := t1.UnixNano() / 1e6
	endTime := t2.UnixNano() / 1e6
	fmt.Printf("SyncCacheTodayReport, startTime:%v, endTime:%v \n", startTime, endTime)
	reportOverview := &entity.ReportOverview{}
	syncReportBetweenTime(startTime, endTime, reportOverview)
	syncReportByTime(endTime, reportOverview)
	reportOverview.Date = startTime

	value, err := json.Marshal(reportOverview)
	if err != nil {
		log.Errorf("SyncCacheTodayReport json.Marshal reportOverview err:[%v]", err)
	}

	err = config.RedisCli.Set(TodayOverviewKey, string(value), 0).Err()
	if err != nil {
		log.Errorf("SyncCacheTodayReport set err:[%v]", err)
	}

	log.Info("SyncCacheTodayReport handle done")
}

func lockYesterdayOverview(key string) bool {
	value, err := config.RedisCli.Get(key).Result()
	if err == redis.Nil {
		return true
	} else if err != nil {
		log.Errorf("lock redis get value error :[%v]\n", err)
		return false
	}
	if value == "" {
		return true
	} else {
		return false
	}
}
