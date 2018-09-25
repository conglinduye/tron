package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/service"
	"net/http"
	"github.com/wlcy/tron/explorer/web/entity"
)

func reportRegister(ginRouter *gin.Engine) {
	ginRouter.GET("/api/stats/overview", func(c *gin.Context) {
		reportResp := &entity.ReportResp{}
		data := make([]*entity.ReportOverview, 0)
		reportResp, err := service.QueryReport()
		if err != nil {
			reportResp.Data = data
			reportResp.Success = false
			c.JSON(http.StatusOK, reportResp)
			return
		}

		if reportResp == nil {
			reportResp.Data = data
			reportResp.Success = false
			c.JSON(http.StatusOK, reportResp)
			return
		}

		c.JSON(http.StatusOK, reportResp)
	})


	ginRouter.GET("/api/stats/overview/init", func(c *gin.Context) {
		service.SyncInitReport()
		c.JSON(http.StatusOK, "handle done")
	})
}
