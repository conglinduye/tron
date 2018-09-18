package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/web/buffer"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryTransactionsBuffer ...
func QueryTransactionsBuffer(req *entity.Transactions) (*entity.TransactionsResp, error) {
	transactions := &entity.TransactionsResp{}
	if req.Number != "" { //按blockID查询
		transactions.Data = buffer.GetBlockBuffer().GetTransactionByBlockID(mysql.ConvertStringToInt64(req.Number, 0))
		transactions.Total = int64(len(transactions.Data))
	} else if req.Hash != "" { //按照交易hash查询
		transact := buffer.GetBlockBuffer().GetTransactionByHash(req.Hash)
		if transact == nil {
			transact, _ = QueryTransaction(req)
		}
		transacts := make([]*entity.TransactionInfo, 0)
		transacts = append(transacts, transact)
		transactions.Data = transacts
		transactions.Total = int64(len(transactions.Data))
	} else if req.Address != "" { //按照交易所属人查询，包含转出的交易，和转入的交易
		transactions, _ = QueryTransactions(req)
	} else { //分页查询
		transactions.Data = buffer.GetBlockBuffer().GetTransactions(req.Start, req.Limit)
		transactions.Total = buffer.GetBlockBuffer().GetTotalTransactions()
	}

	return transactions, nil
}

//QueryTransactions 条件查询  	//?sort=-number&limit=1&count=true&number=2135998 TODO: cache
func QueryTransactions(req *entity.Transactions) (*entity.TransactionsResp, error) {
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,
			trx_hash,contract_data,result_data,fee,
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
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and (owner_address='%v' or to_address='%v'", req.Address, req.Address)
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v create_time", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}

		if strings.Index(v, "number") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v block_id", sortSQL)
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

	return module.QueryTransactionsRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryTransactionByHashFromBuffer 精确查询
func QueryTransactionByHashFromBuffer(req *entity.Transactions) (*entity.TransactionInfo, error) {
	return buffer.GetBlockBuffer().GetTransactionByHash(req.Hash), nil
}

//QueryTransaction 精确查询  	//number=2135998   TODO: cache
func QueryTransaction(req *entity.Transactions) (*entity.TransactionInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,
		trx_hash,contract_data,result_data,fee,
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
