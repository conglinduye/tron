package module

import (
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryTransfersRealize 操作数据库
func QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL, filterTempSQL string) (*entity.TransfersResp, error) {
	strFullSQL := strSQL + " " + filterSQL + "" + filterTempSQL + " " + sortSQL + " " + pageSQL
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

	//查询该语句所查到的数据集合
	var total = int64(len(transferInfos))
	total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	transfersResp.Total = total
	transfersResp.Data = transferInfos

	return transfersResp, nil

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
