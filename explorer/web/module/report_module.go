package module

import (
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

// QueryReportBlock
func QueryReportBlock(startTime, endTime int64) (*entity.ReportBlock, error) {
	strSQL := fmt.Sprintf(` 
	select count(1) as totalCount, sum(block_size) as totalSize, sum(transaction_num) as totalTransaction
	from blocks
    where create_time >= %v and create_time < %v `, startTime, endTime)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryReportBlock error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryReportBlock dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var reportBlock = &entity.ReportBlock{}

	for dataPtr.NextT() {
		reportBlock.TotalCount = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalCount"))
		reportBlock.TotalSize = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalSize"))
		reportBlock.TotalTransaction = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTransaction"))

	}
	return reportBlock, nil
}

// QueryTotalReportBlock
func QueryTotalReportBlock(dateTime int64) (*entity.ReportBlock, error) {
	strSQL := fmt.Sprintf(` 
	select count(1) as totalCount, sum(block_size) as totalSize, sum(transaction_num) as totalTransaction
	from blocks
    where create_time < %v `, dateTime)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalReportBlock error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalReportBlock dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var reportBlock = &entity.ReportBlock{}

	for dataPtr.NextT() {
		reportBlock.TotalCount = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalCount"))
		reportBlock.TotalSize = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalSize"))
		reportBlock.TotalTransaction = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTransaction"))

	}
	return reportBlock, nil
}

//QueryTotalReportTransaction
func QueryTotalReportTransaction(dateTime int64) (int64, error) {
	var totalTransaction = int64(0)
	strSQL := fmt.Sprintf(`
    select count(1) as totalTransaction
	from transactions 
	where create_time < %v `, dateTime)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalReportTransactions error :[%v]\n", err)
		return totalTransaction, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalReportTransactions dataPtr is nil ")
		return totalTransaction, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalTransaction = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTransaction"))
	}

	return totalTransaction, nil
}

// QueryReportAccount
func QueryReportAccount(startTime, endTime int64) (int64, error) {
	var totalAccount = int64(0)
	strSQL := fmt.Sprintf(` 
	select count(1) as totalAccount
	from tron_account
    where create_time >= %v and create_time < %v `, startTime, endTime)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryReportAccounts error :[%v]\n", err)
		return totalAccount, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryReportAccounts dataPtr is nil ")
		return totalAccount, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalAccount = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalAccount"))
	}

	return totalAccount, nil
}

//QueryTotalReportAccount
func QueryTotalReportAccount(dateTime int64) (int64, error) {
	var totalAccount = int64(0)
	strSQL := fmt.Sprintf(`
    select count(1) as totalAccount
	from tron_account
	where create_time < %v `, dateTime)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalReportAccounts error :[%v]\n", err)
		return totalAccount, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalReportAccounts dataPtr is nil ")
		return totalAccount, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalAccount = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalAccount"))
	}
	return totalAccount, nil

}

// QueryTotalStatistics
func QueryTotalStatistics() (int64, error) {
	var totalStatistics = int64(0)
	strSQL := fmt.Sprintf(`
    select count(1) as totalStatistics
	from wlcy_statistics `)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalStatistics error :[%v]\n", err)
		return totalStatistics, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalStatistics dataPtr is nil ")
		return totalStatistics, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalStatistics = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalStatistics"))
	}
	return totalStatistics, nil

}

// InsertStatistics
func InsertStatistics(overview *entity.ReportOverview) error {
	strSQL := fmt.Sprintf(`
		insert into wlcy_statistics 
		(date, avg_block_time, avg_block_size, new_block_seen, new_transaction_seen, new_address_seen, 
         total_block_count, total_transaction, total_address, blockchain_size)
		values(%v, %v, %v, %v, %v, %v, %v, %v, %v, %v)`,
		overview.Date, overview.AvgBlockTime, overview.AvgBlockSize, overview.NewBlockSeen, overview.NewTransactionSeen,
		overview.NewAddressSeen, overview.TotalBlockCount, overview.TotalTransaction, overview.TotalAddress, overview.BlockchainSize)
	insID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("insert logo url fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("insert logo url success, insert id: [%v]", insID)
	return nil
}

//UpdateStatistics
func UpdateStatistics(overview *entity.ReportOverview) error {
	strSQL := fmt.Sprintf(`
	update wlcy_statistics
	set avg_block_time=%v, avg_block_size=%v, new_block_seen=%v, new_transaction_seen=%v, new_address_seen=%v, total_block_count=%v, 
	total_transaction=%v, total_address=%v, blockchain_size=%v where date=%v`,
		overview.AvgBlockTime, overview.AvgBlockSize, overview.NewBlockSeen, overview.NewTransactionSeen,
		overview.NewAddressSeen, overview.TotalBlockCount, overview.TotalTransaction, overview.TotalAddress,
		overview.BlockchainSize, overview.Date)
	_, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("update statistics result fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("update statistics result success  sql:%s", strSQL)
	return nil
}

//QueryStatistics
func QueryStatistics(strSQL string) ([]*entity.ReportOverview, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryStatistics error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryStatistics dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	reportOverviews := make([]*entity.ReportOverview, 0)

	for dataPtr.NextT() {
		reportOverview := &entity.ReportOverview{}
		reportOverview.Date = mysql.ConvertDBValueToInt64(dataPtr.GetField("date"))
		reportOverview.AvgBlockTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("avg_block_time"))
		reportOverview.AvgBlockSize = mysql.ConvertDBValueToInt64(dataPtr.GetField("avg_block_size"))
		reportOverview.NewBlockSeen = mysql.ConvertDBValueToInt64(dataPtr.GetField("new_block_seen"))
		reportOverview.NewTransactionSeen = mysql.ConvertDBValueToInt64(dataPtr.GetField("new_transaction_seen"))
		reportOverview.NewAddressSeen = mysql.ConvertDBValueToInt64(dataPtr.GetField("new_address_seen"))
		reportOverview.TotalBlockCount = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_block_count"))
		reportOverview.TotalTransaction = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_transaction"))
		reportOverview.TotalAddress = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_address"))
		reportOverview.BlockchainSize = mysql.ConvertDBValueToInt64(dataPtr.GetField("blockchain_size"))

		reportOverviews = append(reportOverviews, reportOverview)
	}

	return reportOverviews, nil
}
