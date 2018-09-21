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

	//?sort=-number&limit=1&count=true&number=2135998
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
		req := &entity.VoteWitnessReq{}
		req.Start = c.Query("start")
		req.Limit = c.Query("limit")
		req.Address = c.Query("address")

		resp, err := service.QueryVoteWitness(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/witness/:address", func(c *gin.Context) {
		address := c.Param("address")
		resp, _ := service.QueryVoteWitnessDetail(address)
		c.JSON(http.StatusOK, resp)
	})

}
