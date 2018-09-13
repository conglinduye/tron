package module

import (
	"fmt"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

// QueryReportBlock
func QueryReportBlock(startTime, endTime int64) (*entity.ReportBlock, error) {
	strSQL := fmt.Sprintf(` 
	select count(1) as totalCount, sum(block_size) as totalSize 
	from blocks
    where create_time >= %v and create_time < %v `, startTime, endTime)
	log.Debug(strSQL)
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

	}
	return reportBlock, nil
}

// QueryTotalReportBlock
func QueryTotalReportBlock(dateTime int64) (*entity.ReportBlock, error) {
	strSQL := fmt.Sprintf(` 
	select count(1) as totalCount, sum(block_size) as totalSize 
	from blocks
    where create_time < %v `, dateTime)
	log.Debug(strSQL)
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
	log.Debug(strSQL)
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
	log.Debug(strSQL)
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




