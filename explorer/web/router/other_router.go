package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/service"
)

func otherRegister(ginRouter *gin.Engine) {

	//获得数据同步信息
	ginRouter.GET("/api/system/status", func(c *gin.Context) {
		log.Debugf("Hello /api/system/status")
		resp, err := service.QuerySystemStatus()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	//交易所交易信息
	ginRouter.GET("/api/market/markets", func(c *gin.Context) {
		log.Debugf("Hello /api/market/markets")
		resp, err := service.QueryMarkets()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	//验签
	/*ginRouter.GET("/api/auth", func(c *gin.Context) {
		log.Debugf("Hello /api/auth")
		resp, err := service.QueryAuth()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	//申请测试币
	ginRouter.GET("/api/testnet/request-coins", func(c *gin.Context) {
		log.Debugf("Hello /api/testnet/request-coins")
		resp, err := service.QueryTestRequestCoin()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})*/

}
