package module

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
)

// QueryAllAssetBlacklist
func QueryAllAssetBlacklist() ([]*entity.AssetBlacklist, error) {
	assetBlackLists := make([]*entity.AssetBlacklist, 0)
	strSQL := fmt.Sprintf(`select owner_address, asset_name from wlcy_asset_blacklist`)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryAssetBlacklist error :[%v]", err)
		return assetBlackLists, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAssetBlacklist dataPtr is nil ")
		return assetBlackLists, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		assetBlackList := &entity.AssetBlacklist{}
		assetBlackList.Id = dataPtr.GetField("id")
		assetBlackList.OwnerAddress = dataPtr.GetField("owner_address")
		assetBlackList.TokenName = dataPtr.GetField("asset_name")
		assetBlackList.CreateTime = dataPtr.GetField("create_time")
		assetBlackLists = append(assetBlackLists, assetBlackList)
	}

	return assetBlackLists, nil
}

//QueryAssetBlacklist
func QueryAssetBlacklist(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AssetBlacklistResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Info(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAssetBlacklist error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAssetBlacklist dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	assetBlacklistResp := &entity.AssetBlacklistResp{}
	assetBlackLists := make([]*entity.AssetBlacklist, 0)

	for dataPtr.NextT() {
		assetBlackList := &entity.AssetBlacklist{}
		assetBlackList.Id = dataPtr.GetField("id")
		assetBlackList.OwnerAddress = dataPtr.GetField("owner_address")
		assetBlackList.TokenName = dataPtr.GetField("asset_name")
		assetBlackList.CreateTime = dataPtr.GetField("create_time")
		assetBlackLists = append(assetBlackLists, assetBlackList)

	}

	var total = int64(len(assetBlackLists))
	total, err = mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}

	assetBlacklistResp.Total = total
	assetBlacklistResp.Data = assetBlackLists

	return assetBlacklistResp, nil
}

//InsertAssetBlacklist
func InsertAssetBlacklist(ownerAddress, tokenName string) error {
	strSQL := fmt.Sprintf(`
		insert into wlcy_asset_blacklist 
		(owner_address,asset_name)
		values('%v','%v')`,
		ownerAddress, tokenName)
	insID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("InsertAssetBlacklist fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("InsertAssetBlacklist success, insert id: [%v]", insID)
	return nil
}

//DeleteAssetBlacklist
func DeleteAssetBlacklist(id uint64) error {
	strSQL := fmt.Sprintf(`
		delete from wlcy_asset_blacklist where id=%v`, id)
	insID, _, err := mysql.ExecuteDeleteSQLCommand(strSQL)
	if err != nil {
		log.Errorf("DeleteAssetBlacklist fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("DeleteAssetBlacklist success, insert id: [%v]", insID)
	return nil
}

// InsertAssetExtInfo
func InsertAssetExtInfo(info *entity.AssetExtInfo) error {
	strSQL := fmt.Sprintf(`
		insert into wlcy_asset_info 
		(address, token_name, token_id, brief, website, white_paper, github, country, credit, reddit,
         twitter, facebook, telegram, steam, medium, webchat, weibo, review, status)
		values('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v',
		'%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v)`,
		info.Address, info.TokenName, info.TokenId, info.Brief, info.Website, info.WhitePaper,
		info.Github, info.Country, info.Credit, info.Reddit, info.Twitter, info.Facebook,
		info.Telegram, info.Steam, info.Medium, info.Webchat, info.Weibo, info.Review, info.Status)
	insertID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("InsertAssetExtInfo fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("InsertAssetExtInfo success, insert id: [%v]", insertID)
	return nil
}

// UpdateAssetExtInfo
func UpdateAssetExtInfo(info *entity.AssetExtInfo) error {
	strSQL := fmt.Sprintf(`
		update wlcy_asset_info set token_name = '%v', white_paper = '%v', github = '%v', country = '%v', credit = %v,
		reddit = '%v', twitter = '%v', facebook = '%v', telegram = '%v', steam = '%v', medium = '%v', webchat = '%v',
		weibo = '%v', review = %v, status = %v, update_time = '%v' where address = '%v' `,
		info.TokenName, info.WhitePaper, info.Github, info.Country, info.Credit,
		info.Reddit, info.Twitter, info.Facebook, info.Telegram, info.Steam, info.Medium, info.Webchat,
		info.Weibo, info.Review, info.Status, info.UpdateTime, info.Address)
	updateID, _, err := mysql.ExecuteSQLCommand(strSQL, false)
	if err != nil {
		log.Errorf("UpdateAssetExtInfo fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("UpdateAssetExtInfo success, update id: [%v]", updateID)
	return nil
}


// InsertAssetExtLogo
func InsertAssetExtLogo(logo *entity.AssetExtLogo) error {
	strSQL := fmt.Sprintf(`
		insert into wlcy_asset_logo 
		(address, logo_url) values('%v', '%v')`,
		logo.Address, logo.LogoUrl)
	insertID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("InsertAssetExtLogo fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("InsertAssetExtLogo success, insert id: [%v]", insertID)
	return nil
}

// UpdateAssetExtInfo
func UpdateAssetExtLogo(logo *entity.AssetExtLogo) error {
	strSQL := fmt.Sprintf(`
		update wlcy_asset_logo set address = '%v', logo_url = '%v', update_time = '%v' where address = '%v' `,
		logo.Address, logo.LogoUrl, logo.UpdateTime, logo.Address)
	updateID, _, err := mysql.ExecuteSQLCommand(strSQL, false)
	if err != nil {
		log.Errorf("UpdateAssetExtLogo fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("UpdateAssetExtLogo success, update id: [%v]", updateID)
	return nil
}








