package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"fmt"
	"time"
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


func InsertAssetExtInfo(info *entity.AssetExtInfo) error {
	return module.InsertAssetExtInfo(info)
}

func UpdateAssetExtInfo(info *entity.AssetExtInfo) error {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	updateTime := tm.Format("2006-01-02 15:04:05")
	info.UpdateTime = updateTime
	return module.UpdateAssetExtInfo(info)
}

func InsertAssetExtLogo(logo *entity.AssetExtLogo) error {
	return module.InsertAssetExtLogo(logo)
}

func UpdateAssetExtLogo(logo *entity.AssetExtLogo) error {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	updateTime := tm.Format("2006-01-02 15:04:05")
	logo.UpdateTime = updateTime
	return module.UpdateAssetExtLogo(logo)
}

