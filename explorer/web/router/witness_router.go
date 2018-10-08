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
		//resp, err := service.QueryWitness()
		resp, err := service.QueryWitnessBuffer()
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


// @Summary Query witness
// @Description Query witness
// @Tags Witness
// @Accept  json
// @Produce  json
// @Success 200 {string} json "[{"address":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp","name":"Sesameseed","url":"https://www.sesameseed.org","producer":true,"latestBlockNumber":3021119,"latestSlotNumber":512995842,"missedTotal":232,"producedTotal":94778,"producedTrx":0,"votes":694559265,"producePercentage":99.75521745552766,"votesPercentage":8.503402454074505}...]"
// @Router /api/witness [get]
func QueryWitness(c *gin.Context) {
	log.Debugf("Hello /api/witness")
	//resp, err := service.QueryWitness()
	resp, err := service.QueryWitnessBuffer()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Query witness statistic
// @Description Query witness statistic
// @Tags Witness
// @Accept  json
// @Produce  json
// @Success 200 {string} json "[{"address":"TBXB5tobBPCFkC8ihFDBWjaiwW2iSpzSfr","name":"TRONVIETNAM","url":"https://www.tronvietnam.com/","blockProduced":115,"total":3106,"percentage":0.037}...]"
// @Router /api/witness/maintenance-statistic [get]
func QueryWitnessStatistic(c *gin.Context) {
	log.Debugf("Hello /api/witness/maintenance-statistic")
	//resp, err := service.QueryWitnessStatistic()
	resp, err := service.QueryWitnessStatisticBuffer()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}
