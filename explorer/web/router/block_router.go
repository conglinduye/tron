package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/service"
)

func blockRegister(ginRouter *gin.Engine) {

	//?sort=-number&limit=1&count=true&number=2135998
	ginRouter.GET("/api/block", func(c *gin.Context) {
		blockReq := &entity.Blocks{}
		blockReq.Sort = c.Query("sort")
		blockReq.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		blockReq.Count = c.Query("count")
		blockReq.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		blockReq.Order = c.Query("order")
		blockReq.Number = c.Query("number")
		blockReq.Producer = c.Query("producer")
		log.Debugf("Hello /api/block?%#v", blockReq)
		//blockResp, err := service.QueryBlocks(blockReq)
		blockResp, err := service.QueryBlocksBuffer(blockReq)
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
		//blockResp, err := service.QueryBlock(blockReq)
		blockResp, err := service.QueryBlockBuffer(blockReq)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, blockResp)
	})

}
