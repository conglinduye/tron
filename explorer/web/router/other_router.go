package router

import (
	"net/http"

	"github.com/wlcy/tron/explorer/lib/config"

	"github.com/wlcy/tron/explorer/lib/websocket"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
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
	ginRouter.GET("/api/auth", func(c *gin.Context) {
		req := &entity.Auth{}
		if c.BindJSON(req) == nil {
			if req == nil {
				log.Errorf("parsing request parameter err!")
				c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
			}
		}
		log.Debugf("Hello /api/auth %#v", req)
		resp, err := service.QueryAuth(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//申请测试币
	ginRouter.GET("/api/testnet/request-coins", func(c *gin.Context) {
		if config.NetType != "mainnet" {
			//获取header
			realIP := c.Request.Header.Get("X-Real-IP")
			log.Debugf("Hello /api/testnet/request-coins get Header[X-Real-IP]:%#v", realIP)
			req := &entity.TestCoin{}
			if c.BindJSON(req) == nil {
				if req == nil {
					log.Errorf("parsing request parameter err!")
					c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
				}
			}
			log.Debugf("Hello /api/testnet/request-coins %#v", req)
			resp, err := service.QueryTestRequestCoin(req, realIP)
			if err != nil {
				errCode, _ := util.GetErrorCode(err)
				c.JSON(errCode, err)
			}
			c.JSON(http.StatusOK, resp)
		}
	})

}
