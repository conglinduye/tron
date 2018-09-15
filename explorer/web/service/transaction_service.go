package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryTransactions 条件查询  	//?sort=-number&limit=1&count=true&number=2135998 TODO: cache
func QueryTransactions(req *entity.Transactions) (*entity.TransactionsResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,
			trx_hash,contract_data,result_data,fee
			contract_type,confirmed,create_time,expire_time
			from tron.transactions
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Hash != "" {
		filterSQL = fmt.Sprintf(" and trx_hash='%v'", req.Hash)
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortTemp = fmt.Sprintf("%v ,", sortTemp)
			}
			sortTemp = fmt.Sprintf("%v create_time", sortTemp)
			if strings.Index(v, "-") == 0 {
				sortTemp = fmt.Sprintf("%v desc", sortTemp)
			}
			mutiFilter = true
		}

		if strings.Index(v, "number") > 0 {
			if mutiFilter {
				sortTemp = fmt.Sprintf("%v ,", sortTemp)
			}
			sortTemp = fmt.Sprintf("%v block_id", sortTemp)
			if strings.Index(v, "-") == 0 {
				sortTemp = fmt.Sprintf("%v desc", sortTemp)
			}
			mutiFilter = true
		}
	}
	if sortTemp != "" {
		if strings.Index(sortTemp, ",") == 0 {
			sortTemp = sortTemp[1:]
		}
		sortTemp = fmt.Sprintf("order by %v", sortTemp)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryTransactionsRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryTransaction 精确查询  	//number=2135998   TODO: cache
func QueryTransaction(req *entity.Transactions) (*entity.TransactionInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,
		trx_hash,contract_data,result_data,fee
		contract_type,confirmed,create_time,expire_time
		from tron.transactions
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Hash != "" {
		filterSQL = fmt.Sprintf(" and trx_hash='%v'", req.Hash)
	}
	return module.QueryTransactionRealize(strSQL, filterSQL)
}

//PostTransaction 创建交易
/*func PostTransaction(req *entity.PostTransaction) (*entity.TransactionInfo, error) {
	if req.Transaction == "" {
		log.Errorf("no transaction received")
		return nil, util.NewErrorMsg(util.Error_common_request_json_no_data)
	}

	return module.QueryTransactionRealize(strSQL, filterSQL)
}
*/
