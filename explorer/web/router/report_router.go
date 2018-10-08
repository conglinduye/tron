package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/service"
	"net/http"
	"github.com/wlcy/tron/explorer/web/entity"
)

// @Summary Query statistics
// @Description Query statistics
// @Tags statistics
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"success":true,"data":[{"date":1529884800000,"totalTransaction":3084,"avgBlockTime":3,"avgBlockSize":197,"totalBlockCount":26548,"newAddressSeen":2179,"blockchainSize":5243100,"totalAddress":2208,"newBlockSeen":26547,"newTransactionSeen":3081}...]}"
// @Router /api/stats/overview [get]
func ReportOverview(c *gin.Context) {
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
}

func ReportOverviewInit(c *gin.Context) {
	service.SyncInitReport()
	c.JSON(http.StatusOK, "handle done")
}

func SyncPersistYesterdayReport(c *gin.Context) {
	service.SyncPersistYesterdayReport()
	c.JSON(http.StatusOK, "handle done")
}
