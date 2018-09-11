package router

import (
	"net/http"

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
		req.Limit = c.Query("limit")
		req.Count = c.Query("count")
		req.Start = c.Query("start")
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

	ginRouter.GET("/api/vote/live", func(c *gin.Context) {
		log.Debugf("Hello /api/vote/live")
		resp, err := service.QueryVoteLive()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/current-cycle", func(c *gin.Context) {
		log.Debugf("Hello /api/vote/current-cycle")
		resp, err := service.QueryVoteCurrentCycle()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	ginRouter.GET("/api/vote/next-cycle", func(c *gin.Context) {
		log.Debugf("Hello /api/vote/next-cycle")
		resp, err := service.QueryVoteNextCycle()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

}
