package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"fmt"
)

func QueryAssetBlacklist(req *entity.AssetBlacklistReq) (*entity.AssetBlacklistResp, error) {
	var filterSQL, sortSQL, pageSQL string
	strSQL := fmt.Sprintf(`select id, owner_address, asset_name, create_time from wlcy_asset_blacklist where 1=1 `)

	if req.OwnerAddress != "" {
		filterSQL = fmt.Sprintf(" and owner_address='%v'", req.OwnerAddress)
	}
	if req.TokenName != "" {
		filterSQL = fmt.Sprintf(" and asset_name='%v'", req.TokenName)
	}

	sortSQL = fmt.Sprintf("order by id desc")
	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryAssetBlacklist(strSQL, filterSQL, sortSQL, pageSQL)
}

func InsertAssetBlacklist(ownerAddress, tokenName string) error {
	return module.InsertAssetBlacklist(ownerAddress, tokenName)
}

func DeleteAssetBlacklist(id uint64) error {
	return module.DeleteAssetBlacklist(id)
}

