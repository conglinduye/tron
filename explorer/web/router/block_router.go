package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/service"
)

func blockRegister(ginRouter *gin.Engine) {

	//?sort=-number&limit=1&count=true&number=2135998
	ginRouter.GET("/api/block", func(c *gin.Context) {
		blockReq := &entity.Blocks{}
		blockReq.Sort = c.Query("sort")
		blockReq.Limit = c.Query("limit")
		blockReq.Count = c.Query("count")
		blockReq.Start = c.Query("start")
		blockReq.Order = c.Query("order")
		blockReq.Number = c.Query("number")
		log.Debugf("Hello /api/block?%#v", blockReq)
		blockResp, err := service.QueryBlocks(blockReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, blockResp)
	})
	//:number=2135998
	ginRouter.GET("/api/block/:number", func(c *gin.Context) {
		blockReq := &entity.Blocks{}
		blockReq.Number = c.Param("number") //占位符传参
		log.Debugf("Hello /api/block/:%#v", blockReq.Number)
		blockResp, err := service.QueryBlock(blockReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, blockResp)
	})

}
