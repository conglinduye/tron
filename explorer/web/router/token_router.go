package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/util"
	"net/http"
	"strings"
	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"sync/atomic"
)

func tokenRegister(ginRouter *gin.Engine) {
	ginRouter.GET("/api/token", func(c *gin.Context) {
		tokenReq := &entity.Token{}
		tokenReq.Start = c.Query("start")
		tokenReq.Limit = c.Query("limit")
		tokenReq.Owner = c.Query("owner")
		tokenReq.Name = c.Query("name")
		tokenReq.Status = c.Query("status")
		log.Debugf("Hello /api/token?%#v", tokenReq)
		if tokenReq.Start == "" || tokenReq.Limit == "" {
			tokenReq.Start = "0"
			tokenReq.Limit = "20"
		}
		tokenResp := &entity.TokenResp{}
		tokenList := make([]*entity.TokenInfo, 0)

		if tokenReq.Owner == "" && tokenReq.Name == "" && tokenReq.Status == "" {
			log.Info("service.QueryCommonTokenListBuffer")
			tokenList, _ = service.QueryCommonTokenListBuffer()
		} else if tokenReq.Status != "" && tokenReq.Status == "ico" {
			log.Info("service.QueryIcoTokenListBuffer")
			tokenList, _ = service.QueryIcoTokenListBuffer()
		} else if tokenReq.Owner != "" && tokenReq.Name != "" {
			log.Info("service.QueryTokenDetailListBuffer")
			tokenList,  _ = service.QueryTokenDetailListBuffer()
			tokenList, total := hanldeTokenDetail(tokenReq.Owner, tokenReq.Name, tokenList)
			tokenResp.Total = total
			tokenResp.Data = tokenList
			c.JSON(http.StatusOK, tokenResp)
			return
		} else if tokenReq.Name != "" && strings.HasPrefix(tokenReq.Name, "%") && strings.HasSuffix(tokenReq.Name, "%") {
			log.Info("service.QueryCommonTokenListBuffer NameFuzzyQuery")
			tokenList, _ = service.QueryCommonTokenListBuffer()
			tokenList, total := hanldeTokenList4FuzzyQuery(tokenReq.Name, tokenList)
			tokenResp.Total = total
			tokenResp.Data = tokenList
			c.JSON(http.StatusOK, tokenResp)
			return
		} else {
			log.Info("service.QueryCommonTokenListBuffer OtherQuery")
			tokenList, _ = service.QueryCommonTokenListBuffer()
			tokenList, total := hanldeTokenList4QueryCondition(tokenReq, tokenList)
			tokenResp.Total = total
			tokenResp.Data = tokenList
			c.JSON(http.StatusOK, tokenResp)
			return
		}

		// copyTokenList
		tokenList = copyTokenList(tokenList)

		// paging handle
		length := len(tokenList)
		tokenResp.Total = int64(length)
		start := mysql.ConvertStringToInt(tokenReq.Start, 0)
		limit := mysql.ConvertStringToInt(tokenReq.Limit, 0)
		if start > length {
			tokenResp.Data = make([]*entity.TokenInfo, 0)
		} else {
			if start + limit < length {
				tokenResp.Data = tokenList[start:start+limit]
			} else {
				tokenResp.Data = tokenList[start:length]
			}
		}

		// handleTokenListIndex
		handleTokenListIndex(tokenReq, tokenResp.Data)

		c.JSON(http.StatusOK, tokenResp)
	})

	ginRouter.GET("/api/token/:name", func(c *gin.Context) {
		tokenReq := &entity.Token{}
		name := c.Param("name")
		tokenReq.Name = name
		log.Debugf("Hello /api/token/:%#v", name)
		log.Info("service.QueryTokenDetailListBuffer")
		tokenList,  _ := service.QueryTokenDetailListBuffer()
		tokenList, _ = hanldeTokenList4QueryCondition(tokenReq, tokenList)
		tokenInfo := &entity.TokenInfo{}
		if len(tokenList) > 0 {
			tokenInfo = tokenList[0]
		}
		c.JSON(http.StatusOK, tokenInfo)
	})


	ginRouter.GET("/api/token/:name/address", func(c *gin.Context) {
		tokenReq := &entity.Token{}
		tokenReq.Name = c.Param("name")
		tokenReq.Start = c.Query("start")
		tokenReq.Limit = c.Query("limit")

		if tokenReq.Start == "" || tokenReq.Limit == "" {
			tokenReq.Start = "0"
			tokenReq.Limit = "50"
		}

		assetBalanceResp, err := service.QueryAssetBalances(tokenReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, assetBalanceResp)
	})


	ginRouter.POST("/api/uploadLogo", func(c *gin.Context) {
		res := &entity.UploadLogoRes{}
		var uploadLogoReq entity.UploadLogoReq
		if err := c.Bind(&uploadLogoReq); err != nil {
			res.Success = false
			c.JSON(http.StatusBadRequest, res)
			return
		}

		if uploadLogoReq.ImageData == "" || uploadLogoReq.Address == "" {
			res.Success = false
			c.JSON(http.StatusBadRequest, res)
			return
		}
		//传入data格式：data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAACgCAYA...
		if len(strings.Split(uploadLogoReq.ImageData, ",")) > 1 {
			uploadLogoReq.ImageData = strings.Split(uploadLogoReq.ImageData, ",")[1]
		}

		dst, err := service.UploadTokenLogo(config.DefaultPath, config.ImgURL, uploadLogoReq.ImageData, uploadLogoReq.Address)

		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			res.Success = false
			c.JSON(errCode, res)
			return
		}

		res.Success =true
		res.Data = dst
		c.JSON(http.StatusOK, res)
	})

	ginRouter.GET("/api/download/tokenInfo", func(c *gin.Context) {
		res := &entity.TokenDownloadInfoRes{}
		tokenFile := config.TokenTemplateFile
		if tokenFile == "" {
			tokenFile = "http://coin.top/tokenTemplate/TronscanTokenInformationSubmissionTemplate.xlsx"
		}
		res.Success = true
		res.Data = tokenFile
		c.JSON(http.StatusOK, res)
	})

	ginRouter.GET("/api/sync/participated", func(c *gin.Context) {
		service.SyncAssetIssueParticipated()
		c.JSON(http.StatusOK, "handle done")
	})
}


// hanldeTokenDetail
func hanldeTokenDetail(address string, name string, tokenList []*entity.TokenInfo) ([]*entity.TokenInfo, int64) {
	newTokenInfoList := make([]*entity.TokenInfo, 0)
	tokenInfo := &entity.TokenInfo{}
	for _, token := range tokenList {
		if token.OwnerAddress == address && token.Name == name {
			tokenInfo = token
			tokenInfo.Index = 1
			newTokenInfoList = append(newTokenInfoList, tokenInfo)
			break
		}
	}
	total := len(tokenList)
	return newTokenInfoList, int64(total)
}

// hanldeTokenList4QueryCondition
func hanldeTokenList4QueryCondition(tokenReq *entity.Token, tokenList []*entity.TokenInfo) ([]*entity.TokenInfo, int64) {
	newTokenInfoList := make([]*entity.TokenInfo, 0)
	for _, tokenInfo := range tokenList {
		if tokenReq.Owner != "" && tokenInfo.OwnerAddress == tokenReq.Owner {
			temp := new(entity.TokenInfo)
			*temp = *tokenInfo
			newTokenInfoList = append(newTokenInfoList, temp)
		}
		if tokenReq.Name != "" && tokenInfo.Name == tokenReq.Name {
			newTokenInfoList = append(newTokenInfoList, tokenInfo)
		}
	}
	total := len(tokenList)
	return newTokenInfoList, int64(total)
}

// hanldeTokenList4FuzzyQuery
func hanldeTokenList4FuzzyQuery(name string, tokenList []*entity.TokenInfo) ([]*entity.TokenInfo, int64) {
	rs := []rune(name)
	name = string(rs[1:len(name)-1])
	newTokenInfoList := make([]*entity.TokenInfo, 0)
	for _, tokenInfo := range tokenList {
		if strings.Contains(tokenInfo.Name, name) {
			temp := new(entity.TokenInfo)
			*temp = *tokenInfo
			newTokenInfoList = append(newTokenInfoList, temp)
		}

	}
	total := len(tokenList)
	return newTokenInfoList, int64(total)
}


// handleTokenListIndex
func handleTokenListIndex(req *entity.Token, tokenList []*entity.TokenInfo) {
	var index = mysql.ConvertStringToInt32(req.Start, 0)

	for _, token := range tokenList {
		atomic.AddInt32(&index, 1)
		token.Index = index
	}
}

// copyTokenList
func copyTokenList(tokenList []*entity.TokenInfo) []*entity.TokenInfo {
	newTokenList := make([]*entity.TokenInfo, 0, len(tokenList))
	for _, tokenInfo := range tokenList {
		newTokenInfo := new(entity.TokenInfo)
		*newTokenInfo = *tokenInfo
		newTokenList = append(newTokenList, newTokenInfo)
	}
	return newTokenList
}
