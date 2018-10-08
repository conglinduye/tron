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

// QueryTransfers ...
// @Summary QueryTransfers ...
// @Description Query transfer
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param start query string false "start"
// @Param limit query string false "limit"
// @Param number query string false "number"
// @Param hash query string false "hash"
// @Param block query string false "block"
// @Param address query string false "address"
// @Success 200 {string} json "{total":0,"data":[]}"
// @Router /api/transfer [get]
func QueryTransfers(c *gin.Context) {
	req := &entity.Transfers{}
	req.Sort = c.Query("sort")
	req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
	req.Count = c.Query("count")
	req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
	req.Hash = c.Query("hash")
	req.Address = c.Query("address")
	req.Number = c.Query("number")
	if c.Query("block") != "" { //也能用block过滤
		req.Number = c.Query("block")
	}
	req.Total = mysql.ConvertStringToInt64(c.Query("total"), 0) // 分页查询传入上次总计结果
	log.Debugf("Hello /api/transfer?%#v", req)
	resp, err := service.QueryTransfersBuffer(req)
	//resp, err := service.QueryTransfers(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryTransferByHash ...
// @Summary QueryTransferByHash ...
// @Description Query transfer by hash
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param hash query string false "hash"
// @Success 200 {string} json "{}"
// @Router /api/transfer/:hash [get]
func QueryTransferByHash(c *gin.Context) {
	req := &entity.Transfers{}
	req.Hash = c.Param("hash") //占位符传参
	log.Debugf("Hello /api/transfer/:%#v", req.Hash)
	resp, err := service.QueryTransferByHashFromBuffer(req)
	//resp, err := service.QueryTransfer(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}
