package module

import (
	"encoding/json"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryContractsRealize 操作数据库
func QueryContractsRealize(strSQL, filterSQL, sortSQL, pageSQL string, needTotal bool) (*entity.ContractsResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	contractsResp := &entity.ContractsResp{}
	status := &entity.State{}
	if err != nil || dataPtr == nil {
		log.Errorf("QueryContractsRealize error :[%v]\n", err)
		status.Code = util.Error_common_internal_error
		status.Message = util.GetErrorMsgSleek(util.Error_common_internal_error)
		contractsResp.Status = status
		return contractsResp, util.NewErrorMsg(util.Error_common_internal_error)
	}

	contractListInfos := make([]*entity.ContractListInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var contractListInfo = &entity.ContractListInfo{}
		contractListInfo.Address = dataPtr.GetField("address")
		contractListInfo.Name = dataPtr.GetField("contract_name")
		contractListInfo.Compiler = dataPtr.GetField("compiler_version")
		createTime := dataPtr.GetField("verify_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		contractListInfo.DateVerified = mysql.ConvertDBValueToInt64(createTime)
		contractListInfo.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		contractListInfo.TrxCount = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxNum"))
		isOptimized := dataPtr.GetField("is_optimized")
		if isOptimized == "1" {
			contractListInfo.IsSetting = true
		}

		contractListInfos = append(contractListInfos, contractListInfo)
	}

	var total = int64(len(contractListInfos))
	if needTotal {
		//查询该语句所查到的数据集合
		total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL) //
		if err != nil {
			log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
		}
	}
	contractsResp.Total = total
	status.Code = util.OK
	status.Message = util.Success
	contractsResp.Data = contractListInfos

	return contractsResp, nil

}

//QueryContractsByAddressRealize 操作数据库
func QueryContractsByAddressRealize(strSQL, filterSQL, sortSQL, pageSQL string, needTotal bool) (*entity.ContractBaseResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	contractsResp := &entity.ContractBaseResp{}
	status := &entity.State{}
	if err != nil || dataPtr == nil {
		log.Errorf("QueryContractsByAddressRealize error :[%v]\n", err)
		status.Code = util.Error_common_internal_error
		status.Message = util.GetErrorMsgSleek(util.Error_common_internal_error)
		contractsResp.Status = status
		return contractsResp, util.NewErrorMsg(util.Error_common_internal_error)
	}

	contractInfos := make([]*entity.ContractBaseInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var contract = &entity.ContractBaseInfo{}
		contract.Address = dataPtr.GetField("address")
		contract.TokenTracke = dataPtr.GetField("tokenContract")
		contract.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		contract.TrxCount = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxNum"))
		creator := &entity.Creator{}
		creator.Address = dataPtr.GetField("owner_address")
		creator.TokenBalance = 0 //TODO
		creator.TxHash = dataPtr.GetField("trx_hash")
		contract.Creator = creator
		contractInfos = append(contractInfos, contract)
	}
	status.Code = util.OK
	status.Message = util.Success
	contractsResp.Data = contractInfos

	return contractsResp, nil

}

//QueryContractTnxRealize 。。。
func QueryContractTnxRealize(strSQL, filterSQL, sortSQL, pageSQL string, needTotal bool) (*entity.ContractTransactionResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	contractsResp := &entity.ContractTransactionResp{}
	status := &entity.State{}
	if err != nil || dataPtr == nil {
		log.Errorf("QueryContractTnxRealize error :[%v]\n", err)
		status.Code = util.Error_common_internal_error
		status.Message = util.GetErrorMsgSleek(util.Error_common_internal_error)
		contractsResp.Status = status
		return contractsResp, util.NewErrorMsg(util.Error_common_internal_error)
	}

	contractInfos := make([]*entity.ContractTransactionInfo, 0)

	//填充数据
	for dataPtr.NextT() {
		var contract = &entity.ContractTransactionInfo{}
		contract.TxHash = dataPtr.GetField("trx_hash")
		contract.Block = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		contract.OwnAddress = dataPtr.GetField("owner_address")
		contract.ToAddress = dataPtr.GetField("contract_address")
		contract.Value = mysql.ConvertDBValueToInt64(dataPtr.GetField("call_value"))
		contract.Token = dataPtr.GetField("call_data") //TODO
		createTime := dataPtr.GetField("create_time")
		if len(createTime) > 13 {
			createTime = createTime[:13]
		}
		contract.Timestamp = mysql.ConvertDBValueToInt64(createTime)
		contractInfos = append(contractInfos, contract)
	}
	var total = int64(len(contractInfos))
	if needTotal {
		//查询该语句所查到的数据集合
		total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL) //
		if err != nil {
			log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
		}
	}
	contractsResp.Total = total
	status.Code = util.OK
	status.Message = util.Success
	contractsResp.Data = contractInfos

	return contractsResp, nil
}

//QueryContractsCodeRealize ...
func QueryContractsCodeRealize(strSQL, filterSQL, sortSQL, pageSQL string, needTotal bool) (*entity.ContractCodeResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	contractsResp := &entity.ContractCodeResp{}
	status := &entity.State{}
	if err != nil || dataPtr == nil {
		log.Errorf("QueryContractsCodeRealize error :[%v]\n", err)
		status.Code = util.Error_common_internal_error
		status.Message = util.GetErrorMsgSleek(util.Error_common_internal_error)
		contractsResp.Status = status
		return contractsResp, util.NewErrorMsg(util.Error_common_internal_error)
	}
	//填充数据
	var contractCodeInfo = &entity.ContractCodeInfo{}
	for dataPtr.NextT() {
		contractCodeInfo.Address = dataPtr.GetField("address")
		contractCodeInfo.Name = dataPtr.GetField("contract_name")
		contractCodeInfo.Compiler = dataPtr.GetField("compiler_version")
		contractCodeInfo.Source = dataPtr.GetField("source_code")
		contractCodeInfo.ByteCode = dataPtr.GetField("byte_code")
		contractCodeInfo.ABI = dataPtr.GetField("abi")
		contractCodeInfo.AbiEncoded = dataPtr.GetField("abi_encoded")
		lib := dataPtr.GetField("contract_library")
		librarys := make([]*entity.Library, 0)
		if lib != "" {
			if err := json.Unmarshal([]byte(lib), librarys); err != nil {
				log.Errorf("Unmarshal data failed:[%v]-[%v]", err, lib)
				status.Code = util.Error_common_request_json_convert_error
				status.Message = util.GetErrorMsgSleek(util.Error_common_request_json_convert_error)
				contractsResp.Status = status
				return contractsResp, util.NewErrorMsg(util.Error_common_request_json_convert_error)
			}
		}
		contractCodeInfo.Librarys = librarys
		isOptimized := dataPtr.GetField("is_optimized")
		if isOptimized == "1" {
			contractCodeInfo.IsSetting = true
		}
	}

	status.Code = util.OK
	status.Message = util.Success
	contractsResp.Data = contractCodeInfo

	return contractsResp, nil

}

//QueryContractsCodeVerify ...
func QueryContractsCodeVerify(strSQL string, req *entity.ContractCodeInfo) (*entity.State, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	status := &entity.State{}
	if err != nil || dataPtr == nil {
		log.Errorf("QueryContractsCodeVerify error :[%v]\n", err)
		status.Code = util.Error_common_internal_error
		status.Message = util.GetErrorMsgSleek(util.Error_common_internal_error)
		return status, util.NewErrorMsg(util.Error_common_internal_error)
	}
	isPass := false
	//填充数据
	for dataPtr.NextT() {
		abiMain := dataPtr.GetField("abi")
		byteCodeMain := dataPtr.GetField("byte_code")
		if abiMain == req.ABI && req.ByteCode == byteCodeMain {
			isPass = true
		}
	}
	status.Code = 2001
	status.Message = "verify failed!"
	if isPass {
		status.Code = util.OK
		status.Message = util.Success
		err = req.InsertContractCode()
		if err != nil {
			status.Code = 3001
			status.Message = "verify success! insert db failed"
		}

	}
	return status, err
}
