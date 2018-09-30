package module

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
	"time"
)

//QueryTokenList
func QueryTokenList(strSQL, filterSQL, sortSQL, pageSQL string) ([]*entity.TokenInfo, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTokensList error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTokensList dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	tokenList := make([]*entity.TokenInfo, 0)

	for dataPtr.NextT() {
		token := &entity.TokenInfo{}
		token.OwnerAddress = dataPtr.GetField("owner_address")
		token.Name = dataPtr.GetField("asset_name")
		token.TotalSupply = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_supply"))
		token.TrxNum = mysql.ConvertDBValueToInt64(dataPtr.GetField("trx_num"))
		token.Num = mysql.ConvertDBValueToInt64(dataPtr.GetField("num"))
		token.Participated = mysql.ConvertDBValueToInt64(dataPtr.GetField("participated"))
		token.EndTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("end_time"))
		token.StartTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("start_time"))
		token.VoteScore = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_score"))
		token.Description = string(utils.HexDecode(dataPtr.GetField("asset_desc")))
		token.Url = dataPtr.GetField("url")
		frozenJson := dataPtr.GetField("frozen_supply")
		if frozenJson == "" {
			frozenJson = "[]"
		}
		tokenFrozenInfoList := make([]*entity.TokenFrozenInfo, 0)
		tokenFrozenSupplyList := make([]*entity.TokenFrozenSupply, 0)
		err := json.Unmarshal([]byte(frozenJson), &tokenFrozenSupplyList)
		if err != nil {
			log.Errorf("QueryTokenList json.Unmarshal error :[%v]\n", err)
		} else {
			for _, tokenFrozenSupply := range tokenFrozenSupplyList {
				tokenFrozenInfo := &entity.TokenFrozenInfo{}
				tokenFrozenInfo.Amount = tokenFrozenSupply.FrozenBalance
				t := time.Now()
				datetime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),t.Minute(),0,0, time.UTC)
				days := (tokenFrozenSupply.ExpireTime - datetime.UnixNano() / 1e6) / (1000 * 60 * 60 * 24)
				tokenFrozenInfo.Days = days
				tokenFrozenInfoList = append(tokenFrozenInfoList, tokenFrozenInfo)
			}
			token.Frozen = tokenFrozenInfoList
		}
		token.Abbr = dataPtr.GetField("asset_abbr")

		if dataPtr.IsFieldExist("totalTokenTransfers") {
			token.TotalTransactions = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTokenTransfers"))
		}

		if dataPtr.IsFieldExist("totalTokenHolders") {
			token.NrOfTokenHolders = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTokenHolders"))
		}

		tokenList = append(tokenList, token)
	}

	return tokenList, nil
}

//QueryTokenBalanceRealize 查询通证余额
func QueryTokenBalanceRealize(strSQL, filterSQL string) (*entity.TokenBalanceInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTokenBalanceRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTokenBalanceRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var tokenBalanceInfo = &entity.TokenBalanceInfo{}

	for dataPtr.NextT() {
		tokenBalanceInfo.Address = dataPtr.GetField("address")
		tokenBalanceInfo.AssetName = dataPtr.GetField("asset_name")
		tokenBalanceInfo.CreatorAddress = dataPtr.GetField("creator_address")
		tokenBalanceInfo.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))

	}
	return tokenBalanceInfo, nil
}

//QueryTotalTokenTransfers
func QueryTotalTokenTransfers(tokenName string) (int64, error) {
	var totalTokenTransfers = int64(0)
	strSQL := fmt.Sprintf(`
    select count(1) as totalTokenTransfers
	from contract_asset_transfer
	where  asset_name = '%v' `, tokenName)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalTokenTransfers error :[%v]\n", err)
		return totalTokenTransfers, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalTokenTransfers dataPtr is nil ")
		return totalTokenTransfers, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalTokenTransfers = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTokenTransfers"))
	}
	return totalTokenTransfers, nil

}

//QueryTotalTokenHolders
func QueryTotalTokenHolders(tokenName string) (int64, error) {
	var totalTokenHolders = int64(0)
	strSQL := fmt.Sprintf(` 
	select count(1) as totalTokenHolders
	from account_asset_balance 
	where asset_name = '%v' `, tokenName)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryTotalTokenHolders error :[%v]\n", err)
		return totalTokenHolders, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTotalTokenHolders dataPtr is nil ")
		return totalTokenHolders, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		totalTokenHolders = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalTokenHolders"))
	}
	return totalTokenHolders, nil
}

// QueryTokenExtInfo
func QueryTokenExtInfo() ([]*entity.TokenExtInfo, error) {
	strSQL := fmt.Sprintf(`
	SELECT logo.address,token_id,token_name,brief, website, white_paper,logo.logo_url,
    	github,country, credit, reddit,twitter,facebook, telegram,steam,
    	medium, webchat,Weibo,review
	FROM wlcy_asset_logo logo
	left join wlcy_asset_info info on logo.address=info.address and info.status=1
	where 1=1 order by info.address `)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("queryTokenExtInfo error :[%v]\n", err)
		return nil, err
	}
	if dataPtr == nil {
		log.Error("dataPtr is nil")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var index int32

	var tokenExtInfos = make([]*entity.TokenExtInfo, 0)
	if dataPtr.ResNum() == 0 {
		tokenExtInfos = InitTokenExtInfos()
	}

	for dataPtr.NextT() {
		atomic.AddInt32(&index, 1)
		tokenExtInfo := &entity.TokenExtInfo{}
		tokenExtInfo.OwnerAddress = dataPtr.GetField("address")
		tokenExtInfo.Country = mysql.SetDefaultVal(dataPtr.GetField("country"), "no_message")
		tokenExtInfo.GitHub = mysql.SetDefaultVal(dataPtr.GetField("github"), "no_message")
		tokenExtInfo.ImgURL = mysql.SetDefaultVal(dataPtr.GetField("logo_url"), "")
		tokenExtInfo.Index = fmt.Sprintf("%v", index)
		//转换信用评级  0-Pending，1-Ok ，2-Neutral ，3-insufficient_message，4-Fake
		var reputation = "insufficient_message"
		credit := dataPtr.GetField("credit")
		if credit == "0" {
			reputation = "Pending"
		} else if credit == "1" {
			reputation = "Ok"
		} else if credit == "2" {
			reputation = "Neutral"
		} else if credit == "4" {
			reputation = "Fake"
		}
		tokenExtInfo.Reputation = reputation
		tokenExtInfo.TokenID = dataPtr.GetField("token_id")
		tokenExtInfo.WebSite = mysql.SetDefaultVal(dataPtr.GetField("website"), "no_message")
		tokenExtInfo.WhitePaper = mysql.SetDefaultVal(dataPtr.GetField("white_paper"), "no_message")
		var socialMediaInfos = make([]*entity.SocialMediaInfo, 0)
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Reddit", URL: dataPtr.GetField("reddit")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Twitter", URL: dataPtr.GetField("twitter")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Facebook", URL: dataPtr.GetField("facebook")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Telegram", URL: dataPtr.GetField("telegram")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Steem", URL: dataPtr.GetField("steam")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Medium", URL: dataPtr.GetField("medium")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Wechat", URL: dataPtr.GetField("webchat")})
		socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Weibo", URL: dataPtr.GetField("Weibo")})

		tokenExtInfo.SocialMedia = socialMediaInfos
		tokenExtInfos = append(tokenExtInfos, tokenExtInfo)
	}
	return tokenExtInfos, nil
}

// InitTokenExtInfos
func InitTokenExtInfos() []*entity.TokenExtInfo {
	var tokenExtInfos = make([]*entity.TokenExtInfo, 0)
	var socialMediaInfos = make([]*entity.SocialMediaInfo, 0)
	tokenExtInfo := &entity.TokenExtInfo{}
	tokenExtInfo.Country = "no_message"
	tokenExtInfo.GitHub = "no_message"
	tokenExtInfo.Reputation = "insufficient_message"
	tokenExtInfo.WebSite = "no_message"
	tokenExtInfo.WhitePaper = "no_message"
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Reddit", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Twitter", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Facebook", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Telegram", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Steem", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Medium", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Wechat", URL: ""})
	socialMediaInfos = append(socialMediaInfos, &entity.SocialMediaInfo{Name: "Weibo", URL: ""})

	tokenExtInfo.SocialMedia = socialMediaInfos

	tokenExtInfos = append(tokenExtInfos, tokenExtInfo)

	return tokenExtInfos
}

//InsertLogoInfo
func InsertLogoInfo(address, url string) error {
	strSQL := fmt.Sprintf(`
		insert into wlcy_asset_logo 
		(address,logo_url)
		values('%v','%v')`,
		address, url)
	insID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("insert logo url fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("insert logo url success, insert id: [%v]", insID)
	return nil
}

//UpdateLogoInfo
func UpdateLogoInfo(address, url string) error {
	strSQL := fmt.Sprintf(`
	update wlcy_asset_logo
	set logo_url='%v' where address='%v'`,
		url, address)
	_, _, err := mysql.ExecuteSQLCommand(strSQL, false)
	if err != nil {
		log.Errorf("update logoInfo result fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("update logoInfo result success  sql:%s", strSQL)
	return nil
}

//IsAddressExist
func IsAddressExist(address string) (bool, error) {
	isExist := false
	strSQL := fmt.Sprintf(`
	select address,logo_url from wlcy_asset_logo
	where 1=1 and address='%v' limit 1 `, address)
	//查询结果
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("isAddressNotExist error :[%v]\n", err)
		return isExist, err
	}
	if dataPtr == nil {
		log.Error("dataPtr is nil")
		return isExist, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr.ResNum() > 0 {
		isExist = true
	}
	return isExist, nil
}

// QueryAllAssetIssue
func QueryAllAssetIssue() ([]*entity.AssetIssue, error) {
	strSQL := fmt.Sprintf(`
		select owner_address, asset_name, participated
		from asset_issue `)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryAllAssetIssue error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAllAssetIssue dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	assetIssues := make([]*entity.AssetIssue, 0)

	for dataPtr.NextT() {
		assetIssue := &entity.AssetIssue{}
		assetIssue.OwnerAddress = dataPtr.GetField("owner_address")
		assetIssue.AssetName = dataPtr.GetField("asset_name")
		assetIssue.Participated = mysql.ConvertDBValueToInt64(dataPtr.GetField("participated"))
		assetIssues = append(assetIssues, assetIssue)
	}

	return assetIssues, nil
}

// QueryParticipateAsset
func QueryParticipateAsset(toAddress, assetName string) (*entity.ParticipateAsset, error) {
	strSQL := fmt.Sprintf(`
		select asset_name, sum(amount) as totalAmount
		from contract_participate_asset
		where to_address = '%v' and asset_name = '%v'`, toAddress, assetName)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryParticipateAsset error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryParticipateAsset dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	participateAsset := &entity.ParticipateAsset{}

	for dataPtr.NextT() {
		participateAsset.AssetName = dataPtr.GetField("asset_name")
		participateAsset.TotalAmount = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalAmount"))
	}

	return participateAsset, nil
}

// UpdateAssetIssue
func UpdateAssetIssue(address string, assetName string, participated int64) error {
	strSQL := fmt.Sprintf(`
	update asset_issue set participated=%v where owner_address='%v' and asset_name='%v'`,
		participated, address, assetName)
	log.Sql(strSQL)
	_, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("UpdateAssetIssue result fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("UpdateAssetIssue result success  sql:%s", strSQL)
	return nil
}

//QueryAssetBalances
func QueryAssetBalances(req *entity.Token) (*entity.AssetBalanceResp, error) {
	strSQL := fmt.Sprintf(` 
	select address, asset_name, balance
	from account_asset_balance 
	where  asset_name = '%v' order by balance desc `, req.Name)
	pageSQL := fmt.Sprintf(` limit %v, %v `, req.Start, req.Limit)
	fullSQL := strSQL + pageSQL
	log.Sql(fullSQL)
	dataPtr, err := mysql.QueryTableData(fullSQL)
	if err != nil {
		log.Errorf("QueryAssetBalances error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAssetBalances dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	assetBalanceResp := &entity.AssetBalanceResp{}
	assetBalances := make([]*entity.AssetBalance, 0)
	for dataPtr.NextT() {
		assetBalance := &entity.AssetBalance{}
		assetBalance.Address = dataPtr.GetField("address")
		assetBalance.Name = dataPtr.GetField("asset_name")
		assetBalance.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))

		assetBalances = append(assetBalances, assetBalance)

	}

	var total = int64(len(assetBalances))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}

	assetBalanceResp.Total = total
	assetBalanceResp.Data = assetBalances
	return assetBalanceResp, nil
}

// QueryAssetCreateTime
func QueryAssetCreateTime(ownerAddress, tokenName string) (int64, error) {
	var assetCreateTime = int64(0)
	strSQL := fmt.Sprintf(` 
	select c.create_time as createTime
	from asset_issue a, contract_asset_issue b, blocks c
    where a.asset_name = b.asset_name and b.block_id = c.block_id 
	and a.owner_address = '%v' and a.asset_name = '%v'`, ownerAddress, tokenName)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryAssetCreateTime error :[%v]\n", err)
		return assetCreateTime, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAssetCreateTime dataPtr is nil ")
		return assetCreateTime, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for dataPtr.NextT() {
		assetCreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("createTime"))
	}
	return assetCreateTime, nil
}

// QueryAssetTransfer
func QueryAssetTransfer(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AssetTransferResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAssetTransfer error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAssetTransfer dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	assetTransferResp := &entity.AssetTransferResp{}
	assetTransferList := make([]*entity.AssetTransfer, 0)

	for dataPtr.NextT() {
		assetTransfer := &entity.AssetTransfer{}
		assetTransfer.BlockId = mysql.ConvertDBValueToInt64(dataPtr.GetField("block_id"))
		assetTransfer.TransactionHash = dataPtr.GetField("trx_hash")
		assetTransfer.Timestamp = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		assetTransfer.TransferFromAddress = dataPtr.GetField("owner_address")
		assetTransfer.TransferToAddress = dataPtr.GetField("to_address")
		assetTransfer.Amount = mysql.ConvertDBValueToInt64(dataPtr.GetField("amount"))
		assetTransfer.TokenName = dataPtr.GetField("asset_name")

		confirmed := dataPtr.GetField("confirmed")
		if confirmed == "1" {
			assetTransfer.Confirmed = true
		} else {
			assetTransfer.Confirmed = false
		}

		assetTransferList = append(assetTransferList, assetTransfer)
	}

	total, err := mysql.QuerySQLViewCount(strSQL + " " + filterSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}

	assetTransferResp.Total = total
	assetTransferResp.Data = assetTransferList

	return assetTransferResp, nil
}