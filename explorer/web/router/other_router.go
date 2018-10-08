package router

import (
	"net/http"
	"strings"

	"github.com/wlcy/tron/explorer/lib/config"

	"github.com/wlcy/tron/explorer/lib/websocket"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/service"
)

// QuerySystemStatus ...
// @Summary QuerySystemStatus ...
// @Description Query system block data sync status
// @Tags System
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{}"
// @Router /api/system/status [get]
func QuerySystemStatus(c *gin.Context) {
	log.Debugf("Hello /api/system/status")
	resp, err := service.QuerySystemStatus()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryMarkets ...
// @Summary QueryMarkets ...
// @Description Query markets
// @Tags System
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{}"
// @Router /api/market/markets [get]
func QueryMarkets(c *gin.Context) {
	log.Debugf("Hello /api/market/markets")
	//resp, err := service.QueryMarkets()
	resp, err := service.QueryMarketsBuffer()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

//SocketIO ...
func SocketIO(c *gin.Context) {
	log.Debugf("Hello socket.io")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
	c.Set("content-type", "application/json")
	websocket.WsHandler(c.Writer, c.Request)
	c.JSON(http.StatusOK, "ok")
}

//Auth ...
func Auth(c *gin.Context) {
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
}

//RequestCoins ...
func RequestCoins(c *gin.Context) {
	if strings.ToUpper(config.NetType) != "mainnet" {
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
}
