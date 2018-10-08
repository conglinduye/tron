package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/ext/service"
)

// QueryBalanceByAddress ...
// @Summary QueryBalanceByAddress ...
// @Description Query account by address  balance
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param address query string false "address"
// @Success 200 {string} json "{}"
// @Router /api/account/:address/balance [get]
func QueryBalanceByAddress(c *gin.Context) {
	address := c.Param("address") //占位符传参
	log.Debugf("Hello /api/account/:%#v/balance", address)
	resp, err := service.QueryAccountBalance(address)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// GenAccount ...
// @Summary GenAccount ...
// @Description gen account address
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{}"
// @Router /api/account [post]
func GenAccount(c *gin.Context) {
	log.Debugf("Hello /api/account POST")
	resp, err := service.CreateAccount()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}
