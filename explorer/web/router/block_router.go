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

// QueryBlocks ...
// @Summary QueryBlocks ...
// @Description Query blocks
// @Tags Blocks
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param start query string false "start"
// @Param limit query string false "limit"
// @Param order query string false "order"
// @Param number query string false "number"
// @Param producer query string false "producer"
// @Success 200 {string} json "{total":0,"data":[]}"
// @Router /api/block [get]
func QueryBlocks(c *gin.Context) {
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
}

// QueryBlockByID ...
// @Summary QueryBlockByID ...
// @Description Query block by ID
// @Tags Blocks
// @Accept  json
// @Produce  json
// @Param number query string false "number"
// @Success 200 {string} json "{}"
// @Router /api/block/:number [get]
func QueryBlockByID(c *gin.Context) {
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
}
