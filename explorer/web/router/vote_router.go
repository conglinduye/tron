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

func voteRegister(ginRouter *gin.Engine) {

	ginRouter.GET("/api/vote", func(c *gin.Context) {
		req := &entity.Votes{}
		req.Sort = c.Query("sort")
		req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		req.Count = c.Query("count")
		req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		req.Candidate = c.Query("candidate")
		req.Voter = c.Query("voter")
		log.Debugf("Hello /api/vote?%#v", req)
		resp, err := service.QueryVotes(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/next-cycle", func(c *gin.Context) {
		log.Debugf("Hello /api/vote/next-cycle")
		//resp, err := service.QueryVoteNextCycle()
		resp, err := service.QueryVoteNextCycleBuffer()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/witness", func(c *gin.Context) {
		resp, _ := service.QueryVoteWitnessBuffer()
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/witness/:address", func(c *gin.Context) {
		address := c.Param("address")
		resp, _ := service.QueryVoteWitnessDetail(address)
		c.JSON(http.StatusOK, resp)
	})


	ginRouter.GET("/api/account/:address/votes", func(c *gin.Context) {
		address := c.Param("address")
		log.Debugf("Hello /api/account/:%#v//votes", address)
		resp := service.QueryAddressVoter(address)
		c.JSON(http.StatusOK, resp)
	})

}


// @Summary Query votes
// @Description Query votes
// @Tags Votes
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param start query string false "start"
// @Param limit query string false "limit"
// @Param candidate query string false "candidate"
// @Param voter query string false "voter"
// @Success 200 {string} json "{total":29744,"totalVotes":5465411009,"data":[{"voterAddress":"TZ2MikDjJZ9ghGs8Uc37aTyg2NNXzQtU8s","candidateAddress":"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp","votes":1,"candidateUrl":"","candidateName":"Sesameseed","voterAvailableVotes":1}...]}"
// @Router /api/vote [get]
func QueryVotes(c *gin.Context) {
	req := &entity.Votes{}
	req.Sort = c.Query("sort")
	req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
	req.Count = c.Query("count")
	req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
	req.Candidate = c.Query("candidate")
	req.Voter = c.Query("voter")
	log.Debugf("Hello /api/vote?%#v", req)
	resp, err := service.QueryVotes(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Query votes next cycle
// @Description Query votes next cycle
// @Tags Votes
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"nextCycle":13383000}"
// @Router /api/vote/next-cycle [get]
func QueryVotesNextCycle(c *gin.Context) {
	log.Debugf("Hello /api/vote/next-cycle")
	//resp, err := service.QueryVoteNextCycle()
	resp, err := service.QueryVoteNextCycleBuffer()
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}


// @Summary Query votes current cycle
// @Description Query votes current cycle
// @Tags Votes
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"total_votes":8165982956,"candidates":[{"address":"TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV","name":"","url":"http://TronGr10.com","hasPage":false,"votes":100013931,"change_cycle":0,"change_day":7}...]}"
// @Router /api/vote/current-cycle [get]
func QueryVoteCurrentCycle(c *gin.Context) {
	voteCurrentCycleResp, _:= service.QueryVoteCurrentCycleBuffer()
	c.JSON(http.StatusOK, voteCurrentCycleResp)
}


// @Summary Query votes live
// @Description votes live
// @Tags Votes
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"data":{"TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV":{"address":"TCvwc3FV3ssq2rD82rMmjhT4PVXYTsFcKV","votes":100013931}...}"
// @Router /api/vote/live [get]
func QueryVoteLive(c *gin.Context) {
	voteLiveResp, _:= service.QueryVoteLiveBuffer()
	c.JSON(http.StatusOK, voteLiveResp)
}