package router

import (
	"net/http"

	"github.com/wlcy/tron/explorer/lib/websocket"

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
		//resp, err := service.QueryMarkets()
		resp, err := service.QueryMarketsBuffer()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/socket.io/", func(c *gin.Context) {
		log.Debugf("Hello socket.io")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		websocket.WsHandler(c.Writer, c.Request)
		c.JSON(http.StatusOK, "ok")
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
