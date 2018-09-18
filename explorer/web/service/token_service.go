package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
	"strings"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"encoding/base64"
	"bytes"
	"time"
	"os"
	"image"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"image/gif"
	"io"
	"errors"
	"github.com/wlcy/tron/explorer/web/buffer"
)

// QueryCommonTokensBuffer
func QueryCommonTokensBuffer(req *entity.Token) (*entity.TokenResp, error) {
	tokenBuffer := buffer.GetTokenBuffer()
	commonTokenResp, _ := tokenBuffer.GetCommonTokenResp()

	return commonTokenResp, nil

}

// QueryIcoTokensBuffer
func QueryIcoTokensBuffer(req *entity.Token) (*entity.TokenResp, error) {
	tokenBuffer := buffer.GetTokenBuffer()
	icoTokenResp, _ := tokenBuffer.GetIcoTokenResp()

	return icoTokenResp, nil

}

//QueryTokens
func QueryTokens(req *entity.Token) (*entity.TokenResp, error) {
	var filterSQL, sortSQL string

	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 and asset_name not in('XP', 'WWGoneWGA', 'ZTX', 'Fortnite', 'ZZZ', 'VBucks', 'CheapAirGoCoin') `)

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
	if req.Status == "ico" {
		t := time.Now()
		dateTime := t.UnixNano() / 1e6
		filterSQL = fmt.Sprintf(" and start_time<=%v and end_time>=%v", dateTime, dateTime)
	}
	sortSQL = "order by participated desc"

	tokenResp, err := module.QueryTokensRealize(strSQL, filterSQL, sortSQL, "")
	if err != nil {
		log.Errorf("queryTokens list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if len(tokenResp.Data) == 0 {
		return tokenResp, nil
	}

	// calculateTokens
	calculateTokens(tokenResp)

	// filterIcoAssetExpire
	tokenResp.Data = filterIcoAssetExpire(req, tokenResp)

	// queryCreateTime
	tokens := tokenResp.Data
	for index := range tokens {
		token := tokens[index]
		createTime := queryAssetCreateTime(token.Name)
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
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
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

	return tokenListResp, nil
}

//QueryTokens
func QueryToken(name string) (*entity.TokenInfo, error) {
	var filterSQL string

	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, participated, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 and asset_name not in('XP', 'WWGoneWGA', 'ZTX', 'Fortnite', 'ZZZ', 'VBucks', 'CheapAirGoCoin') `)

	filterSQL = fmt.Sprintf(" and asset_name='%v'", name)

	token, err := module.QueryTokenRealize(strSQL, filterSQL)
	if err != nil {
		log.Errorf("queryToken list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	if token.Name == "" {
		return  nil, nil
	}

	// calculateToken
	calculateToken(token)


	// queryCreateTime
	createTime := queryAssetCreateTime(token.Name)
	token.DateCreated = createTime

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

//UploadTokenLogo 保存图片
func UploadTokenLogo(defaultPath, imgURL, imageData, address string) (string, error) {
	imgData, err := base64.StdEncoding.DecodeString(imageData)
	if nil != err {
		return "", err
	}
	buffer := bytes.NewBuffer(imgData)
	tempFileName := fmt.Sprintf("tokenLogo_%v.jpeg", time.Now().Format("20060102150405.000000"))

	dist, err := os.Create(defaultPath + "/" + tempFileName)
	if err != nil {
		log.Error(err)
	}
	defer dist.Close()

	err = scale(buffer, dist, 0, 0, 0)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("save file %v ", dist)
	dst := fmt.Sprintf("%v/%v", imgURL, tempFileName)
	err = InsertOrUpdateLogo(address, dst)

	return dst, err
}


/*
* 缩略图生成
* 入参:
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度为0则精度保持不变
* 矩形坐标系起点是左上
* 返回:error
 */
func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		log.Error(err)
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Resize(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	log.Debugf("fm:%v", fm)
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
		/*case "bmp":
		return bmp.Encode(out, canvas)*/
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}


//InsertOrUpdateLogo 更新或插入logo
func InsertOrUpdateLogo(address, url string) error {
	if address == "" || url == "" {
		log.Error("address or url is nil")
		return util.NewErrorMsg(util.Error_common_internal_error)
	}

	addressInfo, err := module.IsAddressExist(address)
	if err != nil {
		log.Errorf("check address:[%v] isExist err:[%v]", address, err)
		return err
	}
	log.Errorf("addressInfo:[%#v] ", addressInfo)
	if addressInfo {
		err = module.UpdateLogoInfo(address, url)
	} else {
		err = module.InsertLogoInfo(address, url)
	}
	return err
}

func SyncAssetIssueParticipated() {
	assetIssues, _ := module.QueryAllAssetIssue()
	if len(assetIssues) == 0 {
		log.Info("SyncAssetIssueParticipated len(assetIssues) == 0")
		return
	}
	for index := range assetIssues {
		assetIssue := assetIssues[index]
		participateAsset, _ := module.QueryParticipateAsset(assetIssue.AssetName)
		if participateAsset.AssetName != "" && participateAsset.TotalAmount > assetIssue.Participated {
			module.UpdateAssetIssue(assetIssue.AssetName, participateAsset.TotalAmount)
		}
	}

}

// QueryAssetBalances
func QueryAssetBalances(req *entity.Token) (*entity.AssetBalanceResp, error){
	assetBalanceResp, err := module.QueryAssetBalances(req)
	if err != nil {
		log.Errorf("QueryAssetBalances list is nil or err:[%v]", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	return assetBalanceResp, nil
}

// QueryAssetCreateTime
func queryAssetCreateTime(tokenName string) int64 {
	createTime, err := module.QueryAssetCreateTime(tokenName)
	if err != nil {
		log.Errorf("QueryAssetCreateTime list is nil or err:[%v]", err)
		t := time.Now()
		createTime = t.UnixNano() / 1e6
	}
	return createTime
}

// filterIcoAssetExpire
func filterIcoAssetExpire(req *entity.Token, tokenResp *entity.TokenResp)  []*entity.TokenInfo {
	tokens := make([]*entity.TokenInfo, 0)
	data := tokenResp.Data
	for index := range data {
		token := data[index]
		if req.Status !="" && req.Status == "ico" && token.IssuedPercentage == 100 {
			// do nothing
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens
}
