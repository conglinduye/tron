package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
	"strings"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"sync/atomic"
)

//QueryTokens
func QueryTokens(req *entity.Token) (*entity.TokenResp, error) {
	var filterSQL, sortSQL, pageSQL string

	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 `)

	if req.Owner != "" {
		filterSQL = fmt.Sprintf(" and owner_address='%v'", req.Owner)
	}
	if req.Name != "" {
		if strings.HasPrefix(req.Name, "%") && strings.HasSuffix(req.Name, "%") {
			filterSQL = fmt.Sprintf(" and asset_name like '%v'", req.Name)
		} else {
			filterSQL = fmt.Sprintf(" and asset_name='%v'", req.Name)
		}
	}

	sortSQL = "order by participated desc"

	if req.Limit != "" && req.Start != "" {
		pageSQL = fmt.Sprintf(" limit %v, %v", req.Start, req.Limit)
	}

	tokenResp, err := module.QueryTokensRealize(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		log.Errorf("queryTokens list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	// calculateTokens
	calculateTokens(tokenResp)

	tokenAddressList := make([]string, 0)
	for _, tokenInfo := range tokenResp.Data {
		if tokenInfo.OwnerAddress != "" {
			tokenAddressList = append(tokenAddressList, tokenInfo.OwnerAddress)
		}
	}

	tokenExtList, err := module.QueryTokenExtInfo(tokenAddressList)
	if err != nil {
		log.Errorf("queryTokenExtInfo list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	var tokenListResp = &entity.TokenResp{}
	tokenList := make([]*entity.TokenInfo, 0)
	var index = mysql.ConvertStringToInt32(req.Start, 0)
	tokenExtEmptyInfoList := module.InitTokenExtInfos()

	for _, tokenInfo := range tokenResp.Data {
		atomic.AddInt32(&index, 1)
		tokenInfo.Index = index

		for _, tokenExtInfo := range tokenExtList {

			if tokenInfo.OwnerAddress == tokenExtInfo.OwnerAddress {
				//tokenInfo.TokenExtInfo = tokenExtInfo
				tokenInfo.Country = tokenExtInfo.Country
				tokenInfo.GitHub = tokenExtInfo.GitHub
				tokenInfo.ImgURL = tokenExtInfo.ImgURL
				tokenInfo.Reputation = tokenExtInfo.Reputation
				tokenInfo.TokenID = tokenExtInfo.TokenID
				tokenInfo.WebSite = tokenExtInfo.WebSite
				tokenInfo.WhitePaper = tokenExtInfo.WhitePaper
				tokenInfo.SocialMedia = tokenExtInfo.SocialMedia
				break
			} else {
				tokenInfo.ImgURL = tokenExtEmptyInfoList[0].ImgURL
				tokenInfo.Country = tokenExtEmptyInfoList[0].Country
				tokenInfo.GitHub = tokenExtEmptyInfoList[0].GitHub
				tokenInfo.Reputation = tokenExtEmptyInfoList[0].Reputation
				tokenInfo.TokenID = tokenExtEmptyInfoList[0].TokenID
				tokenInfo.WebSite = tokenExtEmptyInfoList[0].WebSite
				tokenInfo.WhitePaper = tokenExtEmptyInfoList[0].WhitePaper
				tokenInfo.SocialMedia = tokenExtEmptyInfoList[0].SocialMedia
			}
		}
		tokenList = append(tokenList, tokenInfo)

		if len(tokenList) > 0 {
			tokenListResp.Data = tokenList
			tokenListResp.Total = tokenResp.Total
		}
	}

	return tokenListResp, nil
}

//QueryTokens
func QueryToken(name string) (*entity.TokenInfo, error) {
	var filterSQL string

	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 `)

	filterSQL = fmt.Sprintf(" and asset_name='%v'", name)

	token, err := module.QueryTokenRealize(strSQL, filterSQL)
	if err != nil {
		log.Errorf("queryToken list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	// calculateToken
	calculateToken(token)

	// QueryTotalTokenTransfers
	totalTokenTransfers, _ := module.QueryTotalTokenTransfers(name)
	token.TotalTransactions = totalTokenTransfers
	// QueryTotalTokenHolders
	totalTokenHolders, _ := module.QueryTotalTokenHolders(name)
	token.NrOfTokenHolders = totalTokenHolders

	tokenAddressList := make([]string, 0)
	tokenAddressList = append(tokenAddressList, token.OwnerAddress)
	tokenExtList, err := module.QueryTokenExtInfo(tokenAddressList)
	if err != nil {
		log.Errorf("queryTokenExtInfo list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	for _, tokenExtInfo := range tokenExtList {
		token.Country = tokenExtInfo.Country
		token.GitHub = tokenExtInfo.GitHub
		token.ImgURL = tokenExtInfo.ImgURL
		token.Reputation = tokenExtInfo.Reputation
		token.TokenID = tokenExtInfo.TokenID
		token.WebSite = tokenExtInfo.WebSite
		token.WhitePaper = tokenExtInfo.WhitePaper
		token.SocialMedia = tokenExtInfo.SocialMedia
		break
	}

	return token, nil
}

// QueryTokenBalance
func QueryTokenBalance(address, tokenName string) (*entity.TokenBalanceInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select address, asset_name, creator_address, balance
		from account_asset_balance
		where 1=1 `)
	filterSQL = fmt.Sprintf(" and address='%v' and asset_name='%v'", address, tokenName)

	return module.QueryTokenBalanceRealize(strSQL, filterSQL)
}

// calculateTokens
func calculateTokens(tokenResp *entity.TokenResp) {
	tokens := tokenResp.Data
	for index := range tokens {
		token := tokens[index]
		calculateToken(token)
	}
}

// calculateToken
func calculateToken(token *entity.TokenInfo) {
	frozen := token.Frozen

	var frozenSupply int64 = 0
	for index := range frozen {
		frozenSupply = frozenSupply + frozen[index].Amount
	}
	totalSupply := token.TotalSupply
	availableSupply := totalSupply - frozenSupply

	var availableTokens int64 = 0
	tokenBalanceInfo, err := QueryTokenBalance(token.OwnerAddress, token.Name)
	if err == nil {
		availableTokens = tokenBalanceInfo.Balance
	}

	issuedTokens := availableSupply - availableTokens

	issuedPercentage := float64(issuedTokens) / float64(totalSupply) * 100
	remainingTokens := totalSupply - frozenSupply - issuedTokens
	percentage := float64(remainingTokens) / float64(totalSupply) * 100
	frozenSupplyPercentage := float64(frozenSupply) / float64(totalSupply) * 100

	if token.Num != 0 {
		price := token.TrxNum / token.Num
		token.Price = price
	} else {
		token.Price = 0
	}
	token.Issued = issuedTokens
	token.IssuedPercentage = issuedPercentage
	token.Available = availableTokens
	token.AvailableSupply = availableSupply
	token.Remaining = remainingTokens
	token.RemainingPercentage = percentage
	token.Percentage = percentage
	token.FrozenTotal = frozenSupply
	token.FrozenPercentage = frozenSupplyPercentage

}
