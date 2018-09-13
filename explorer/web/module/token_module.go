package module

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"encoding/json"
	"fmt"
	"sync/atomic"
)

//QueryTokens
func QueryTokensRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.TokenResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTokens error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTokens dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	tokenResp := &entity.TokenResp{}
	tokens := make([]*entity.TokenInfo, 0)

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
		token.Description = dataPtr.GetField("asset_desc")
		token.Url = dataPtr.GetField("utl")
		frozenJson := dataPtr.GetField("frozen_supply")
		var tokenFrozenInfo []entity.TokenFrozenInfo
		if err := json.Unmarshal([]byte(frozenJson), &tokenFrozenInfo); err == nil {
			token.Frozen = tokenFrozenInfo
		}
		token.Abbr = dataPtr.GetField("asset_abbr")

		tokens = append(tokens, token)
	}

	var total = int64(len(tokens))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	tokenResp.Total = total
	tokenResp.Data = tokens

	return tokenResp, nil
}

//QueryToken
func QueryTokenRealize(strSQL, filterSQL string) (*entity.TokenInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryToken error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryToken dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	token := &entity.TokenInfo{}
	for dataPtr.NextT() {
		token.OwnerAddress = dataPtr.GetField("owner_address")
		token.Name = dataPtr.GetField("asset_name")
		token.TotalSupply = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_supply"))
		token.TrxNum = mysql.ConvertDBValueToInt64(dataPtr.GetField("trx_num"))
		token.Num = mysql.ConvertDBValueToInt64(dataPtr.GetField("num"))
		token.Participated = mysql.ConvertDBValueToInt64(dataPtr.GetField("participated"))
		token.EndTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("end_time"))
		token.StartTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("start_time"))
		token.VoteScore = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_score"))
		token.Description = dataPtr.GetField("asset_desc")
		token.Url = dataPtr.GetField("utl")
		frozenJson := dataPtr.GetField("frozen_supply")
		var tokenFrozenInfo []entity.TokenFrozenInfo
		if err := json.Unmarshal([]byte(frozenJson), &tokenFrozenInfo); err == nil {
			token.Frozen = tokenFrozenInfo
		}
		token.Abbr = dataPtr.GetField("asset_abbr")
	}

	return token, nil
}

//QueryTokenBalanceRealize 查询通证余额
func QueryTokenBalanceRealize(strSQL, filterSQL string) (*entity.TokenBalanceInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
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
	where token_name = '%v' `, tokenName)
	log.Debug(strSQL)
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
	log.Debug(strSQL)
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
func QueryTokenExtInfo(addressList []string) ([]*entity.TokenExtInfo, error) {
	filterSQL := mysql.GenSQLPartInStrList("logo.address", addressList, true)
	strSQL := fmt.Sprintf(`
	SELECT logo.address,token_id,token_name,brief, website, white_paper,logo.logo_url,
    	github,country, credit, reddit,twitter,facebook, telegram,steam,
    	medium, webchat,Weibo,review
	FROM wlcy_asset_logo logo
	left join wlcy_asset_info info on logo.address=info.address and info.status=1
	where 1=1 and %v order by info.address `, filterSQL)
	log.Debug(strSQL)
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
	var socialMediaInfos = make([]*entity.SocialMediaInfo, 0)
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
