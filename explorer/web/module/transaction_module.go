package module

import (
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryTransactionsRealize 操作数据库
func QueryTransactionsRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.TransactionsResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTransactionsRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTransactionsRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	transactionResp := &entity.TransactionsResp{}
	transactionInfos := make([]*entity.TransactionInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var transaction = &entity.TransactionInfo{}
		transaction.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		transaction.Hash = dataPtr.GetField("trx_hash")
		transaction.ToAddress = dataPtr.GetField("to_address")
		transaction.OwnerAddress = dataPtr.GetField("owner_address")
		createTime := dataPtr.GetField("create_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		transaction.CreateTime = mysql.ConvertDBValueToInt64(createTime)
		transaction.ContractType = mysql.ConvertDBValueToInt64(dataPtr.GetField("contract_type"))
		if dataPtr.GetField("contract_data") != "" {
			_, transaction.ContractData = utils.GetContractInfoStr3(int32(transaction.ContractType), utils.HexDecode(dataPtr.GetField("contract_data")))
		}
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			transaction.Confirmed = true
		}

		transactionInfos = append(transactionInfos, transaction)
	}

	//查询该语句所查到的数据集合
	var total = int64(len(transactionInfos))
	total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	transactionResp.Total = total
	transactionResp.Data = transactionInfos

	return transactionResp, nil

}

//QueryTransactionRealize 操作数据库
func QueryTransactionRealize(strSQL, filterSQL string) (*entity.TransactionInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTransactionRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTransactionRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var transaction = &entity.TransactionInfo{}
	//填充数据
	for dataPtr.NextT() {
		transaction.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		transaction.Hash = dataPtr.GetField("trx_hash")
		transaction.ToAddress = dataPtr.GetField("to_address")
		transaction.OwnerAddress = dataPtr.GetField("owner_address")
		createTime := dataPtr.GetField("create_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		transaction.CreateTime = mysql.ConvertDBValueToInt64(createTime)
		transaction.ContractType = mysql.ConvertDBValueToInt64(dataPtr.GetField("contract_type"))
		_, transaction.ContractData = utils.GetContractInfoStr3(int32(transaction.ContractType), utils.HexDecode(dataPtr.GetField("contract_data")))
		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			transaction.Confirmed = true
		}
	}

	return transaction, nil

}
