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
	commonTokenResp *entity.TokenResp
	icoTokenResp    *entity.TokenResp
}

func GetTokenBuffer() *tokenBuffer {
	return getTokenBuffer()
}

func getTokenBuffer() *tokenBuffer {
	onceTokenBuffer.Do(func() {
		_tokenBuffer = &tokenBuffer{}
		_tokenBuffer.loadCommonQueryTokens()
		_tokenBuffer.loadIcoQueryTokens()

		go tokenInfoBufferLoader()
	})

	return _tokenBuffer

}
func tokenInfoBufferLoader() {
	for {
		_tokenBuffer.loadCommonQueryTokens()
		_tokenBuffer.loadIcoQueryTokens()
		time.Sleep(20 * time.Second)
	}
}

func (w *tokenBuffer) GetCommonTokenResp() (commonTokenResp *entity.TokenResp, ok bool) {
	w.RLock()
	if w.commonTokenResp == nil {
		log.Debugf("GetCommonTokenResp from buffer nil, data reload")
		w.loadCommonQueryTokens()
		log.Debugf("GetCommonTokenResp from buffer, buffer data updated ")
	}
	commonTokenResp = w.commonTokenResp
	w.RUnlock()
	return
}

// loadCommonQueryTokens
func (w *tokenBuffer) loadCommonQueryTokens() {
	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 and asset_name not in('XP', 'WWGoneWGA', 'ZTX', 'Fortnite', 'ZZZ', 'VBucks', 'CheapAirGoCoin') 
			order by participated desc `)

	commonTokenResp, err := module.QueryTokensRealize(strSQL, "", "", "")
	if err != nil {
		log.Errorf("loadCommonQueryTokens list is nil or err:[%v]", err)
	}
	if len(commonTokenResp.Data) == 0 {
		return
	}

	subHandle(commonTokenResp)

	w.Lock()
	w.commonTokenResp = commonTokenResp
	w.Unlock()

}

// GetIcoTokenResp
func (w *tokenBuffer) GetIcoTokenResp() (icoTokenResp *entity.TokenResp, ok bool) {
	w.RLock()
	if w.icoTokenResp == nil {
		log.Debugf("GetIcoTokenResp from buffer nil, data reload")
		w.loadIcoQueryTokens()
		log.Debugf("GetIcoTokenResp from buffer, buffer data updated ")
	}
	icoTokenResp = w.icoTokenResp
	w.RUnlock()
	return
}

// loadIcoQueryTokens
func (w *tokenBuffer) loadIcoQueryTokens() {
	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 and asset_name not in('XP', 'WWGoneWGA', 'ZTX', 'Fortnite', 'ZZZ', 'VBucks', 'CheapAirGoCoin') `)

	t := time.Now()
	dateTime := t.UnixNano() / 1e6
	filterSQL := fmt.Sprintf(" and start_time<=%v and end_time>=%v", dateTime, dateTime)

	sortSQL := "order by participated desc"

	icoTokenResp, err := module.QueryTokensRealize(strSQL, filterSQL, sortSQL, "")
	if err != nil {
		log.Errorf("loadIcoQueryTokens list is nil or err:[%v]", err)
	}
	if len(icoTokenResp.Data) == 0 {
		return
	}

	subHandle(icoTokenResp)

	icoTokenResp.Data = filterIcoAssetExpire(icoTokenResp)

	w.Lock()
	w.icoTokenResp = icoTokenResp
	w.Unlock()

}

// subHandle
func subHandle(tokenResp *entity.TokenResp) {
	// calculateTokens
	calculateTokens(tokenResp)

	// queryCreateTime
	tokens := tokenResp.Data
	for index := range tokens {
		token := tokens[index]
		createTime := queryAssetCreateTime(token.OwnerAddress, token.Name)
		tokens[index].DateCreated = createTime
	}

	tokenAddressList := make([]string, 0)
	for _, tokenInfo := range tokenResp.Data {
		if tokenInfo.OwnerAddress != "" {
			tokenAddressList = append(tokenAddressList, tokenInfo.OwnerAddress)
		}
	}

	tokenExtList, err := module.QueryTokenExtInfo(tokenAddressList)
	if err != nil {
		log.Errorf("queryTokenExtInfo list is nil or err:[%v]", err)
	}

	var tokenListResp = &entity.TokenResp{}
	tokenList := make([]*entity.TokenInfo, 0)
	//var index = mysql.ConvertStringToInt32(req.Start, 0)
	tokenExtEmptyInfoList := module.InitTokenExtInfos()

	for _, tokenInfo := range tokenResp.Data {
		//atomic.AddInt32(&index, 1)
		//tokenInfo.Index = index

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

// filterIcoAssetExpire
func filterIcoAssetExpire(tokenResp *entity.TokenResp) []*entity.TokenInfo {
	tokens := make([]*entity.TokenInfo, 0)
	data := tokenResp.Data
	for index := range data {
		token := data[index]
		if token.IssuedPercentage == 100 {
			// do nothing
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens
}
