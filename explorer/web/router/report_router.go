package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/lib/util"
	"net/http"
)

func reportRegister(ginRouter *gin.Engine) {
	ginRouter.GET("/api/stats/overview", func(c *gin.Context) {
		reportResp, err := service.QueryReport()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, reportResp)
	})


	ginRouter.GET("/api/stats/overview/init", func(c *gin.Context) {
		service.SyncInitReport()

		c.JSON(http.StatusOK, "handle done")
	})

}
