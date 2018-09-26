package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

var _tokenBuffer *tokenBuffer
var onceTokenBuffer sync.Once

type tokenBuffer struct {
	sync.RWMutex
	commonTokenList     []*entity.TokenInfo
	icoTokenList        []*entity.TokenInfo
	tokenDetailList     []*entity.TokenInfo
}

func GetTokenBuffer() *tokenBuffer {
	return getTokenBuffer()
}

func getTokenBuffer() *tokenBuffer {
	onceTokenBuffer.Do(func() {
		_tokenBuffer = &tokenBuffer{}
		go tokenListBufferLoader()
		go tokenDetailListBufferLoader()
	})

	return _tokenBuffer
}

func tokenListBufferLoader() {
	for {
		go _tokenBuffer.loadQueryCommonTokenList()
		go _tokenBuffer.loadQueryIcoTokenList()
		time.Sleep(120 * time.Second)
	}
}

func tokenDetailListBufferLoader() {
	for {
		go _tokenBuffer.loadQueryTokensDetailList()
		time.Sleep(180 * time.Second)
	}
}

func (w *tokenBuffer) GetCommonTokenList() (commonTokenList []*entity.TokenInfo) {
	w.RLock()
	commonTokenList = w.commonTokenList
	w.RUnlock()
	return
}

// loadCommonQueryTokenList
func (w *tokenBuffer) loadQueryCommonTokenList() {
	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 order by participated desc `)

	commonTokenList, err := module.QueryTokenList(strSQL, "", "", "")
	if err != nil {
		log.Errorf("loadCommonQueryTokensLists list is nil or err:[%v]", err)
		return
	}
	if len(commonTokenList) == 0 {
		return
	}

	commonTokenList = subHandle(commonTokenList)

	w.Lock()
	w.commonTokenList = commonTokenList
	w.Unlock()

}

// GetIcoTokenList
func (w *tokenBuffer) GetIcoTokenList() (icoTokenList []*entity.TokenInfo) {
	w.RLock()
	icoTokenList = w.icoTokenList
	w.RUnlock()
	return
}

// loadIcoQueryTokens
func (w *tokenBuffer) loadQueryIcoTokenList() {
	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 `)

	t := time.Now()
	dateTime := t.UnixNano() / 1e6
	filterSQL := fmt.Sprintf(" and start_time<=%v and end_time>=%v", dateTime, dateTime)

	sortSQL := "order by participated desc"

	icoTokenList, err := module.QueryTokenList(strSQL, filterSQL, sortSQL, "")
	if err != nil {
		log.Errorf("loadIcoQueryTokens list is nil or err:[%v]", err)
		return
	}
	if len(icoTokenList) == 0 {
		return
	}

	icoTokenList = subHandle(icoTokenList)

	icoTokenList = filterIcoTokenExpire(icoTokenList)

	w.Lock()
	w.icoTokenList = icoTokenList
	w.Unlock()

}

// GetTokensDetailList
func (w *tokenBuffer) GetTokensDetailList() (tokenDetailList []*entity.TokenInfo) {
	w.RLock()
	tokenDetailList = w.tokenDetailList
	w.RUnlock()
	return
}

// loadQueryTokensDetailList
func (w *tokenBuffer) loadQueryTokensDetailList() {
	strSQL := fmt.Sprintf(`
		select a.owner_address, a.asset_name, a.asset_abbr, a.total_supply, a.frozen_supply,
			a.trx_num, a.num, a.participated, a.start_time, a.end_time, a.order_num, a.vote_score, a.asset_desc, a.url,
			b.totalTokenTransfers, c.totalTokenHolders
		from asset_issue a
			left join(
				select asset_name, count(1) as totalTokenTransfers
				from contract_asset_transfer
				group by asset_name
			) b on b.asset_name = a.asset_name
			left join(
				select asset_name, count(1) as totalTokenHolders
				from account_asset_balance 
				group by asset_name
			) c on c.asset_name = a.asset_name
		where 1=1 order by a.participated desc  `)

	tokenDetailList, err := module.QueryTokenList(strSQL, "", "", "")
	if err != nil {
		log.Errorf("loadQueryTokensDetailList list is nil or err:[%v]", err)
		return
	}
	if len(tokenDetailList) == 0 {
		return
	}
	subHandle(tokenDetailList)

	w.Lock()
	w.tokenDetailList = tokenDetailList
	w.Unlock()
}

// subHandle
func subHandle(tokenList []*entity.TokenInfo) []*entity.TokenInfo {
	// filterAssetBlacklist
	tokenList = filterAssetBlacklist(tokenList)

	// calculateTokens
	calculateTokens(tokenList)

	// queryCreateTime
	for _, token := range tokenList {
		createTime := queryAssetCreateTime(token.OwnerAddress, token.Name)
		token.DateCreated = createTime
	}

	tokenExtList, err := module.QueryTokenExtInfo()
	if err != nil {
		log.Errorf("queryTokenExtInfo list is nil or err:[%v]", err)
	}

	newTokenList := make([]*entity.TokenInfo, 0)
	tokenExtEmptyInfoList := module.InitTokenExtInfos()

	for _, tokenInfo := range tokenList {
		flag := true
		for _, tokenExtInfo := range tokenExtList {
			if tokenInfo.OwnerAddress == tokenExtInfo.OwnerAddress {
				tokenInfo.Country = tokenExtInfo.Country
				tokenInfo.GitHub = tokenExtInfo.GitHub
				tokenInfo.ImgURL = tokenExtInfo.ImgURL
				tokenInfo.Reputation = tokenExtInfo.Reputation
				tokenInfo.TokenID = tokenExtInfo.TokenID
				tokenInfo.WebSite = tokenExtInfo.WebSite
				tokenInfo.WhitePaper = tokenExtInfo.WhitePaper
				tokenInfo.SocialMedia = tokenExtInfo.SocialMedia
				flag = false
				break
			}
		}
		if flag == true {
			tokenInfo.Country = tokenExtEmptyInfoList[0].Country
			tokenInfo.GitHub = tokenExtEmptyInfoList[0].GitHub
			tokenInfo.ImgURL = tokenExtEmptyInfoList[0].ImgURL
			tokenInfo.Reputation = tokenExtEmptyInfoList[0].Reputation
			tokenInfo.TokenID = tokenExtEmptyInfoList[0].TokenID
			tokenInfo.WebSite = tokenExtEmptyInfoList[0].WebSite
			tokenInfo.WhitePaper = tokenExtEmptyInfoList[0].WhitePaper
			tokenInfo.SocialMedia = tokenExtEmptyInfoList[0].SocialMedia
		}
		newTokenList = append(newTokenList, tokenInfo)
	}

	return newTokenList

}

// calculateTokens
func calculateTokens(tokenList []*entity.TokenInfo) {
	for _, token := range tokenList {
		calculateToken(token)
	}
}

// calculateToken
func calculateToken(token *entity.TokenInfo) {
	frozenList := token.Frozen

	var frozenSupply int64 = 0
	if frozenList != nil {
		for _, frozen := range frozenList {
			frozenSupply = frozenSupply + frozen.Amount
		}
	} else {
		token.Frozen = make([]entity.TokenFrozenInfo, 0)
	}

	totalSupply := token.TotalSupply
	availableSupply := totalSupply - frozenSupply

	var availableTokens int64 = 0
	tokenBalanceInfo, err := queryTokenBalance(token.OwnerAddress, token.Name)
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

// queryTokenBalance
func queryTokenBalance(address, tokenName string) (*entity.TokenBalanceInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select address, asset_name, creator_address, balance
		from account_asset_balance
		where 1=1 `)
	filterSQL = fmt.Sprintf(" and address='%v' and asset_name='%v'", address, tokenName)

	return module.QueryTokenBalanceRealize(strSQL, filterSQL)
}

// queryAssetCreateTime
func queryAssetCreateTime(ownerAddress, tokenName string) int64 {
	createTime, err := module.QueryAssetCreateTime(ownerAddress, tokenName)
	if err != nil {
		log.Errorf("QueryAssetCreateTime list is nil or err:[%v]", err)
		t := time.Now()
		createTime = t.UnixNano() / 1e6
	}
	return createTime
}

// filterIcoTokenExpire
func filterIcoTokenExpire(tokenList []*entity.TokenInfo) []*entity.TokenInfo {
	newTokenList := make([]*entity.TokenInfo, 0)
	for _, token := range tokenList {
		if token.IssuedPercentage == 100 {
			// do nothing
		} else {
			newTokenList = append(newTokenList, token)
		}
	}
	return newTokenList
}

// filterAssetBlacklist
func filterAssetBlacklist(tokenList []*entity.TokenInfo) []*entity.TokenInfo{
	newTokenInfoList := make([]*entity.TokenInfo, 0)
	assetBlacklists, err := module.QueryAssetBlacklist()
	if err != nil {
		log.Errorf("QueryAssetBlacklist err:[%v]", err)
	}

	if len(assetBlacklists) == 0 {
		return tokenList
	}
	for _, token := range tokenList {
		flag := false
		for _, assetBlacklist := range assetBlacklists {
			if assetBlacklist.OwnerAddress == token.OwnerAddress && assetBlacklist.AssetName == token.Name {
				flag = true
				break
			}
		}
		if flag == false {
			newTokenInfoList =append(newTokenInfoList, token)
		}
	}
	return newTokenInfoList
}