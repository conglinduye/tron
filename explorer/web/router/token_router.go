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
			tokenReq.Limit = "40"
		}

		tokenResp := &entity.TokenResp{}
		var err error = nil
		/*if tokenReq.Owner == "" && tokenReq.Name == "" && tokenReq.Status == "" {
			tokenResp, err = service.QueryCommonTokensBuffer(tokenReq)
		} else if tokenReq.Status != "" && tokenReq.Status == "ico" {
			tokenResp, err = service.QueryIcoTokensBuffer(tokenReq)
		} else {
			tokenResp, err = service.QueryTokens(tokenReq)
		}*/

		tokenResp, err = service.QueryTokens(tokenReq)

		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
			return
		}
		tokenInfos := tokenResp.Data
		length := len(tokenInfos)
		start := mysql.ConvertStringToInt(tokenReq.Start, 0)
		limit := mysql.ConvertStringToInt(tokenReq.Limit, 0)
		if start > length {
			tokenResp.Data = make([]*entity.TokenInfo, 0)
			c.JSON(http.StatusOK, tokenResp)
			return
		}

		if start + limit < length {
			tokenResp.Data = tokenInfos[start:start+limit]
			c.JSON(http.StatusOK, tokenResp)
			return
		} else {
			tokenResp.Data = tokenInfos[start:length]
			return
		}
		c.JSON(http.StatusOK, tokenResp)
	})

	ginRouter.GET("/api/token/:name", func(c *gin.Context) {
		name := c.Param("name")
		log.Debugf("Hello /api/token/:%#v", name)
		tokenInfo, err := service.QueryToken(name)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
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

	ginRouter.GET("/api/mytoken", func(c *gin.Context) {
		tokenReq := &entity.Token{}
		tokenReq.Owner = c.Query("owner")
		log.Debugf("Hello /api/mytoken?%#v", tokenReq)
		log.Debugf("owner_address=%v", tokenReq.Owner)

		if tokenReq.Start == "" || tokenReq.Limit == "" {
			tokenReq.Start = "0"
			tokenReq.Limit = "40"
		}
		tokenResp, err := service.QueryTokens(tokenReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, tokenResp)
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
