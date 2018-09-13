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
		tokenResp, err := service.QueryTokens(tokenReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
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

	ginRouter.GET("/api/mytoken", func(c *gin.Context) {
		tokenReq := &entity.Token{}
		tokenReq.Owner = c.Query("owner")
		log.Debugf("Hello /api/mytoken?%#v", tokenReq)
		log.Debugf("owner_address=%v", tokenReq.Owner)
		if tokenReq.Owner == "" {
			c.JSON(http.StatusBadRequest, nil)
		}

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
		var uploadLogoReq entity.UploadLogoReq
		if err := c.Bind(&uploadLogoReq); err != nil {
			c.JSON(http.StatusBadRequest, nil)
		}

		if uploadLogoReq.ImageData == "" || uploadLogoReq.Address == "" {
			c.JSON(http.StatusBadRequest, nil)
		}
		//传入data格式：data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAACgCAYA...
		if len(strings.Split(uploadLogoReq.ImageData, ",")) > 1 {
			uploadLogoReq.ImageData = strings.Split(uploadLogoReq.ImageData, ",")[1]
		}

		dst, err := service.UploadTokenLogo(config.DefaultPath, config.ImgURL, uploadLogoReq.ImageData, uploadLogoReq.Address)

		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		
		c.JSON(http.StatusOK, dst)
	})

	ginRouter.GET("/api/download/tokenInfo", func(c *gin.Context) {
		tokenFile := config.TokenTemplateFile
		if tokenFile == "" {
			tokenFile = "http://coin.top/tokenTemplate/TronscanTokenInformationSubmissionTemplate.xlsx"
		}

		c.JSON(http.StatusOK, tokenFile)
	})
}
