package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/util"
	"net/http"
)

func srAccountRegister(ginRouter *gin.Engine) {
	ginRouter.GET("/api/srAccount", func(c *gin.Context) {
		srAccountReq := &entity.SrAccount{}
		srAccountReq.Start = c.Query("start")
		srAccountReq.Limit = c.Query("limit")
		srAccountReq.Address = c.Query("address")
		log.Debugf("Hello /api/srAccount?%#v", srAccountReq)
		if srAccountReq.Start == "" || srAccountReq.Limit == "" {
			srAccountReq.Start = "0"
			srAccountReq.Limit = "40"
		}
		srAccountResp, err := service.QuerySrAccounts(srAccountReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, srAccountResp)
	})
}