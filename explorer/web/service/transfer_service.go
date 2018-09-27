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
	if req.Number != "" { //按blockID查询, 分页 总量等于 block transcation
		transfers.Data = buffer.GetBlockBuffer().GetTransferByBlockID(mysql.ConvertStringToInt64(req.Number, 0))
		transfers.Total = int64(len(transfers.Data))
	} else if req.Hash != "" { //按照交易hash查询，不分页
		transact := buffer.GetBlockBuffer().GetTransferByHash(req.Hash)
		if transact == nil {
			transact, _ = QueryTransfer(req)
		}
		transacts := make([]*entity.TransferInfo, 0)
		transacts = append(transacts, transact)
		transfers.Data = transacts
		transfers.Total = int64(len(transfers.Data))
	} else if req.Address != "" { //按照交易所属人查询，包含转出的交易，和转入的交易， 分页 总量等于用户的transactions
		//return QueryTransfers(req)
		return QueryTransfersByAddress(req)
	} else { //分页查询, 分页 总量== totalTransaction
		transfers.Data = buffer.GetBlockBuffer().GetTransfers(req.Start, req.Limit, req.Total)
		transfers.Total = buffer.GetBlockBuffer().GetTotalTransfers()
	}
	return transfers, nil

}

//QueryTransferByHashFromBuffer 从缓存中精确查询  	//number=2135998
func QueryTransferByHashFromBuffer(req *entity.Transfers) (*entity.TransferInfo, error) {
	return buffer.GetBlockBuffer().GetTransferByHash(req.Hash), nil
}

/*
//QueryTransfersByAddress 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryTransfersByAddress(req *entity.Transfers) (*entity.TransfersResp, error) {
	var resp = &entity.TransfersResp{}
	pageSQL := fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from contract_transfer
			where 1=1 and owner_address='%v'`, req.Address)

	transOutResp, err := module.QueryTransfersByAddressRealize(strSQL, pageSQL, true)
	if err != nil {
		log.Errorf("QueryTransfersByAddressRealize query out transfer for address[%v] error[%v] ", req.Address, err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	strSQL = fmt.Sprintf(`
			select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from contract_transfer
			where 1=1 and to_address='%v'`, req.Address)

	transInResp, err := module.QueryTransfersByAddressRealize(strSQL, pageSQL, true)
	if err != nil {
		log.Errorf("QueryTransfersByAddressRealize query in transfer for address[%v] error[%v] ", req.Address, err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	resp.Total = transOutResp.Total + transInResp.Total
	transferInfos := append(transOutResp.Data, transInResp.Data...)
	sort.SliceStable(transferInfos, func(i, j int) bool { return transferInfos[i].CreateTime > transferInfos[j].CreateTime })
	if int64(len(transferInfos)) > req.Limit {
		resp.Data = transferInfos[:req.Limit]
	} else {
		resp.Data = transferInfos
	}

	return resp, nil
}
*/
//QueryTransfers 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryTransfers(req *entity.Transfers) (*entity.TransfersResp, error) {
	var filterSQL, sortSQL, pageSQL, filterTempSQL string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from contract_transfer
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
		filterTempSQL = req.Address
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

	return module.QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL, filterTempSQL, true)
}

//QueryTransfersByAddress  根据地址查询其下所有相关的交易列表
func QueryTransfersByAddress(req *entity.Transfers) (*entity.TransfersResp, error) {
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false
	strSQL := fmt.Sprintf(`
	select oo.block_id,oo.owner_address,oo.to_address,oo.amount,
	oo.asset_name,oo.trx_hash,
	oo.contract_type,oo.confirmed,oo.create_time
	from
	(select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from contract_transfer
			where 1=1 and owner_address='%v'
	union 
	select block_id,owner_address,to_address,amount,
			asset_name,trx_hash,
			contract_type,confirmed,create_time
			from contract_transfer
			where 1=1 and to_address='%v'
	) oo where 1=1`, req.Address, req.Address)

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

	return module.QueryTransfersRealize(strSQL, filterSQL, sortSQL, pageSQL, "", true)
}

//QueryTransfer 精确查询  	//number=2135998   TODO: cache
func QueryTransfer(req *entity.Transfers) (*entity.TransferInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,amount,
		asset_name,trx_hash,
		contract_type,confirmed,create_time
		from contract_transfer
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
