package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryTransfersBuffer ...从缓存中获取数据
func QueryTransfersBuffer(req *entity.Transfers) (*entity.TransfersResp, error) {
	transfers := &entity.TransfersResp{}
	if req.Number != "" { //按blockID查询
		transfers.Data = buffer.GetBlockBuffer().GetTransferByBlockID(mysql.ConvertStringToInt64(req.Number, 0))
		transfers.Total = int64(len(transfers.Data))
	} else if req.Hash != "" { //按照交易hash查询
		transact := buffer.GetBlockBuffer().GetTransferByHash(req.Hash)
		if transact == nil {
			transact, _ = QueryTransfer(req)
		}
		transacts := make([]*entity.TransferInfo, 0)
		transacts = append(transacts, transact)
		transfers.Data = transacts
		transfers.Total = int64(len(transfers.Data))
	} else if req.Address != "" { //按照交易所属人查询，包含转出的交易，和转入的交易
		return QueryTransfers(req)
	} else { //分页查询
		transfers.Data = buffer.GetBlockBuffer().GetTransfers(req.Start, req.Limit)
		transfers.Total = buffer.GetBlockBuffer().GetTotalTransfers()
	}
	return transfers, nil

}

//QueryTransferByHashFromBuffer 从缓存中精确查询  	//number=2135998   TODO: cache
func QueryTransferByHashFromBuffer(req *entity.Transfers) (*entity.TransferInfo, error) {
	return buffer.GetBlockBuffer().GetTransferByHash(req.Hash), nil
}

//QueryTransfers 条件查询  	//?sort=-number&limit=1&count=true&number=2135998  TODO: cache
func QueryTransfers(req *entity.Transfers) (*entity.TransfersResp, error) {
	var filterSQL, sortSQL, pageSQL, filterTempSQL string
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
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and (owner_address='%v' or to_address='%v')", req.Address, req.Address)
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

	/*if filterSQL == "" {
		hourBefore, _ := time.ParseDuration("-24h")
		filterTempSQL = fmt.Sprintf("and create_time>%v", time.Now().Add(hourBefore).UnixNano())
	}*/

	return module.QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL, filterTempSQL)
}

//QueryTransfer 精确查询  	//number=2135998   TODO: cache
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
