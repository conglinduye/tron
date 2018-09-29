package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/util"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryContracts 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryContracts(req *entity.Contracts) (*entity.ContractsResp, error) {
	filterSQL, sortSQL, pageSQL := parsingSQL(req)

	strSQL := fmt.Sprintf(`
	select acc.address,sm.contract_name,sm.compiler_version,
	sm.is_optimized,-- source_code,byte_code,abi,abi_encoded,contract_library,
	sm.verify_time,acc.balance,trxCount.trxNum
	from tron_test_net.tron_account acc
	left join tron_test_net.wlcy_smart_contract sm on acc.address=sm.address 
	left join (
		select contract_address,count(1) as trxNum from tron_test_net.contract_trigger_smart group by contract_address
	) trxCount on trxCount.contract_address=acc.address
	where 1=1 and acc.account_type=2 `)

	return module.QueryContractsRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//QueryContractByAddress 根据地址查询
func QueryContractByAddress(req *entity.Contracts) (*entity.ContractBaseResp, error) {
	filterSQL, sortSQL, pageSQL := parsingSQL(req)

	strSQL := fmt.Sprintf(`
			select sm.address,contract_name,compiler_version,
			is_optimized,verify_time,ta.balance,trxCount.trxNum,
			cr.owner_address,cr.trx_hash,cr.name
			from tron_account ta 
			left join wlcy_smart_contract sm on ta.address=sm.address 
			left join contract_create_smart cr on cr.contract_address=sm.address
			left join (
				select contract_address,count(1) as trxNum from contract_trigger_smart group by contract_address
			) trxCount on trxCount.contract_address=sm.address
			where 1=1 and ta.account_type=2 `)

	return module.QueryContractsByAddressRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//QueryContractTnx 查询交易
func QueryContractTnx(req *entity.Contracts) (*entity.ContractTransactionResp, error) {
	filterSQL, sortSQL, pageSQL := parsingSQL(req)

	strSQL := fmt.Sprintf(`
			select sm.address,sm.trx_hash,sm.block_id,
			sm.create_time,sm.owner_address,sm.contract_address,sm.call_value,
			sm.call_data
			from contract_trigger_smart sm
			where 1=1 `)

	return module.QueryContractTnxRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//QueryContractsCode 查询合约代码信息
func QueryContractsCode(req *entity.Contracts) (*entity.ContractCodeResp, error) {
	filterSQL, sortSQL, pageSQL := parsingSQL(req)

	strSQL := fmt.Sprintf(`
			select sm.address,sm.contract_name,sm.compiler_version,
			sm.is_optimized,sm.verify_time,sm.source_code,sm.byte_code,sm.abi,
			sm.abi_encoded,sm.contract_library
			from wlcy_smart_contract sm
			where 1=1 `)

	return module.QueryContractsCodeRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//VerifyContractCode ...
func VerifyContractCode(req *entity.ContractCodeInfo) (*entity.State, error) {
	status := &entity.State{}

	if req.Address == "" || req.ABI == "" || req.ByteCode == "" {
		status.Code = util.Error_common_parameter_invalid
		status.Message = util.GetErrorMsgSleek(util.Error_common_parameter_invalid)
		return status, util.NewErrorMsg(util.Error_common_parameter_invalid)
	}
	strSQL := fmt.Sprintf(`
			select sm.contract_address,sm.abi,sm.byte_code
			from contract_create_smart sm
			where 1=1 and contract_address='%v'`, req.Address)

	return module.QueryContractsCodeVerify(strSQL, req)

}

func parsingSQL(req *entity.Contracts) (string, string, string) {
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false
	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" %v and sm.address='%v'", filterSQL, req.Address)
	}
	if req.Type == "internal" {
		filterSQL = fmt.Sprintf(" %v and sm.contract_type='%v'", filterSQL, 31) //TODO
	}
	if req.Type == "token" {
		filterSQL = fmt.Sprintf(" %v and sm.contract_type='%v'", filterSQL, 30) //TODO
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v verify_time", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}
	}
	if sortSQL != "" {
		if strings.Index(sortSQL, ",") == 0 {
			sortSQL = sortSQL[1:]
		}
		sortSQL = fmt.Sprintf("order by %v", sortSQL)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	return filterSQL, sortSQL, pageSQL
}
