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

func accountRegister(ginRouter *gin.Engine) {

	//?sort=-balance&limit=1&count=true
	ginRouter.GET("/api/account", func(c *gin.Context) {
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
	})
	//:number=2135998
	ginRouter.GET("/api/account/:address", func(c *gin.Context) {
		req := &entity.Accounts{}
		req.Address = c.Param("address") //占位符传参
		log.Debugf("Hello /api/account/:%#v", req.Address)
		resp, err := service.QueryAccount(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//查询某地址的媒体信息
	ginRouter.GET("/api/account/:address/media", func(c *gin.Context) {
		req := &entity.Accounts{}
		req.Address = c.Param("address") //占位符传参
		log.Debugf("Hello /api/account/:%#v//media", req.Address)
		resp, err := service.QueryAccountMedia(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//修改超级代表github信息
	ginRouter.POST("/api/account/:address/sr", func(c *gin.Context) {
		req := &entity.SuperAccountInfo{}
		if c.BindJSON(req) == nil {
			if req == nil {
				log.Errorf("parsing request parameter err!")
				c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
			}
		}
		log.Debugf("Hello /api/account/:%#v//sr", req)
		resp, err := service.UpdateAccountSr(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//查询超级代表github信息
	ginRouter.GET("/api/account/:address/sr", func(c *gin.Context) {
		req := &entity.SuperAccountInfo{}
		req.Address = c.Param("address") //占位符传参
		log.Debugf("Hello /api/account/:%#v//sr", req.Address)
		resp, err := service.QueryAccountSr(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

}
