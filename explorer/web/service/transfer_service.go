package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryTransfers 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryTransfers(req *entity.Transfers) (*entity.TransfersResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from tron.contract_transfer
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
	if req.Limit != "" && req.Start != "" {
		pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	}
	return module.QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryTransfer 精确查询  	//number=2135998
func QueryTransfer(req *entity.Transfers) (*entity.TransferInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,amount,
		asset_name,trx_hash,
		contract_type,confirmed,create_time
		from tron.contract_transfer
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Hash != "" {
		filterSQL = fmt.Sprintf(" and trx_hash='%v'", req.Hash)
	}
	return module.QueryTransferRealize(strSQL, filterSQL)
}
