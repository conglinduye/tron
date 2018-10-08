package router

import (
	"net/http"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/service"
)

// QueryAccounts ...
// @Summary QueryAccounts ...
// @Description Query account
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param start query string false "start"
// @Param limit query string false "limit"
// @Param address query string false "address"
// @Success 200 {string} json "{total":0,"data":[]}"
// @Router /api/account [get]
func QueryAccounts(c *gin.Context) {
	req := &entity.Accounts{}
	req.Sort = c.Query("sort")
	req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
	req.Count = c.Query("count")
	req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
	req.Address = c.Query("address")
	log.Debugf("Hello /api/account?%#v", req)
	resp, err := service.QueryAccounts(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryAccountByAddress ...
// @Summary QueryAccountByAddress ...
// @Description Query account by address
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param address query string false "address"
// @Success 200 {string} json "{}"
// @Router /api/account/:address [get]
func QueryAccountByAddress(c *gin.Context) {
	req := &entity.Accounts{}
	req.Address = c.Param("address") //占位符传参
	log.Debugf("Hello /api/account/:%#v", req.Address)
	resp, err := service.QueryAccount(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryAccountAddressMedia ...
// @Summary QueryAccountAddressMedia ...
// @Description Query account by address  media
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param address query string false "address"
// @Success 200 {string} json "{}"
// @Router /api/account/:address/media [get]
func QueryAccountAddressMedia(c *gin.Context) {
	req := &entity.Accounts{}
	req.Address = c.Param("address") //占位符传参
	log.Debugf("Hello /api/account/:%#v//media", req.Address)
	if req.Address == "" {
		errCode, _ := util.GetErrorCode(util.NewErrorMsg(util.Error_common_not_suport_parameter))
		c.JSON(errCode, util.NewErrorMsg(util.Error_common_not_suport_parameter))
	}
	resp, err := service.QueryAccountMedia(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// PostAccountAddressSr ...
// @Summary PostAccountAddressSr ...
// @Description Post super account github
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{}"
// @Router /api/account/:address/media [post]
func PostAccountAddressSr(c *gin.Context) {
	//获取header
	token := c.Request.Header.Get("X-Key")
	req := &entity.SuperAccountInfo{}
	if c.BindJSON(req) == nil {
		if req == nil {
			log.Errorf("parsing request parameter err!")
			c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
		}
	}
	log.Debugf("Hello /api/account/:%#v//sr, header token:[%v]", req, token)
	resp, err := service.UpdateAccountSr(req, token)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryAccountAddressSr ...
// @Summary QueryAccountAddressSr ...
// @Description Query account by address  sr
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param address query string false "address"
// @Success 200 {string} json "{}"
// @Router /api/account/:address/sr [get]
func QueryAccountAddressSr(c *gin.Context) {
	req := &entity.SuperAccountInfo{}
	req.Address = c.Param("address") //占位符传参
	log.Debugf("Hello /api/account/:%#v//sr", req.Address)
	if req.Address == "" {
		errCode, _ := util.GetErrorCode(util.NewErrorMsg(util.Error_common_not_suport_parameter))
		c.JSON(errCode, util.NewErrorMsg(util.Error_common_not_suport_parameter))
	}
	resp, err := service.QueryAccountSr(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryAccountAddressStats ...
// @Summary QueryAccountAddressStats ...
// @Description Query account by address  stats
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param address query string false "address"
// @Success 200 {string} json "{}"
// @Router /api/account/:address/stats [get]
func QueryAccountAddressStats(c *gin.Context) {
	address := c.Param("address") //占位符传参
	log.Debugf("Hello /api/account/:%#v//stats", address)
	if address == "" {
		errCode, _ := util.GetErrorCode(util.NewErrorMsg(util.Error_common_not_suport_parameter))
		c.JSON(errCode, util.NewErrorMsg(util.Error_common_not_suport_parameter))
	}
	resp, err := service.QueryAccountStats(address)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)

}
