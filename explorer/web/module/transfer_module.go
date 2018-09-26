package module

import (
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryTransfersRealize 操作数据库
func QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL, filterTempSQL string, needTotal bool) (*entity.TransfersResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTransfersRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTransfersRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	transfersResp := &entity.TransfersResp{}
	transferInfos := make([]*entity.TransferInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var transfer = &entity.TransferInfo{}
		transfer.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		transfer.TransactionHash = dataPtr.GetField("trx_hash")
		transfer.TransferFromAddress = dataPtr.GetField("owner_address")
		createTime := dataPtr.GetField("create_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		transfer.CreateTime = mysql.ConvertDBValueToInt64(createTime)
		transfer.TransferToAddress = dataPtr.GetField("to_address")
		transfer.TokenName = dataPtr.GetField("asset_name")
		transfer.Amount = mysql.ConvertDBValueToInt64(dataPtr.GetField("amount"))
		if transfer.TokenName == "" {
			transfer.TokenName = "TRX"
			//如果是TRX，页面做的单位转换
			//transfer.Amount = transfer.Amount / 1000000
		}
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			transfer.Confirmed = true
		}

		transferInfos = append(transferInfos, transfer)
	}

	var total = int64(len(transferInfos))
	if needTotal {
		//查询该语句所查到的数据集合
		if filterTempSQL != "" { //按地址查询总数  按地址查询太慢，变更查询总计方式
			total, err = querySQLCount(filterTempSQL)
		} else {
			total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL) //
		}

		if err != nil {
			log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
		}
	}
	transfersResp.Total = total
	transfersResp.Data = transferInfos

	return transfersResp, nil

}
func querySQLCount(address string) (int64, error) {
	strFullSQL := fmt.Sprintf(`select count(1) as total
	from tron.contract_transfer
	   where 1=1   and owner_address='%v' 
    union
    select count(1) as total
	  from tron.contract_transfer
	where 1=1   and  to_address='%v'
    `, address, address)
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("querySQLCount error :[%v]\n", err)
		return 0, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("querySQLCount dataPtr is nil ")
		return 0, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var totalNum = int64(0)

	//填充数据
	for dataPtr.NextT() {
		totalNum += mysql.ConvertDBValueToInt64(dataPtr.GetField("total"))
	}
	return totalNum, nil
}

//QueryTransferRealize 操作数据库
func QueryTransferRealize(strSQL, filterSQL string) (*entity.TransferInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTransferRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTransferRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var transfer = &entity.TransferInfo{}
	//填充数据
	for dataPtr.NextT() {
		transfer.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		transfer.TransactionHash = dataPtr.GetField("trx_hash")
		transfer.TransferFromAddress = dataPtr.GetField("owner_address")
		createTime := dataPtr.GetField("create_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		transfer.CreateTime = mysql.ConvertDBValueToInt64(createTime)
		transfer.TransferToAddress = dataPtr.GetField("to_address")
		transfer.TokenName = dataPtr.GetField("asset_name")
		transfer.Amount = mysql.ConvertDBValueToInt64(dataPtr.GetField("amount"))
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			transfer.Confirmed = true
		}
	}

	return transfer, nil

}

//QueryTrxOutByAddress 查询该地址的转出数
func QueryTrxOutByAddress(address string) int64 {
	strSQL := fmt.Sprintf(`select owner_address, count(1) as trxOut from tron.contract_transfer trf where owner_address='%v'`, address)
	log.Sql(strSQL)
	trxOut := int64(0)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil || dataPtr == nil {
		log.Errorf("QueryTrxOutByAddress error :[%v]\n", err)
		return trxOut
	}
	//填充数据
	for dataPtr.NextT() {
		trxOut = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxOut"))
	}
	return trxOut
}

//QueryTrxInByAddress 查询改地址的转入数
func QueryTrxInByAddress(address string) int64 {
	strSQL := fmt.Sprintf(`select to_address, count(1) as trxIn from tron.contract_transfer trf where to_address='%v'`, address)
	log.Sql(strSQL)
	trxIn := int64(0)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil || dataPtr == nil {
		log.Errorf("QueryTrxOutByAddress error :[%v]\n", err)
		return trxIn
	}
	//填充数据
	for dataPtr.NextT() {
		trxIn = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxIn"))
	}
	return trxIn
}
