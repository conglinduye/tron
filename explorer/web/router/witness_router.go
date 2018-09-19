package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/service"
)

func witnessRegister(ginRouter *gin.Engine) {

	ginRouter.GET("/api/witness", func(c *gin.Context) {
		log.Debugf("Hello /api/witness")
		resp, err := service.QueryWitness()
		//resp, err := service.QueryWitnessBuffer()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/witness/maintenance-statistic", func(c *gin.Context) {
		log.Debugf("Hello /api/witness/maintenance-statistic")
		//resp, err := service.QueryWitnessStatistic()
		resp, err := service.QueryWitnessStatisticBuffer()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

}
