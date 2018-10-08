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

// QueryTransactions ...
// @Summary QueryTransactions ...
// @Description Query transactions
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param sort query string false "sort"
// @Param start query string false "start"
// @Param limit query string false "limit"
// @Param number query string false "number"
// @Param hash query string false "hash"
// @Param block query string false "block"
// @Param address query string false "address"
// @Param total query string false "total"
// @Success 200 {string} json "{total":0,"data":[]}"
// @Router /api/transaction [get]
func QueryTransactions(c *gin.Context) {
	req := &entity.Transactions{}
	req.Sort = c.Query("sort")
	req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
	req.Count = c.Query("count")
	req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
	req.Hash = c.Query("hash")
	req.Address = c.Query("address") //按照交易所属人查询，此处包含转入和转出的交易
	req.Number = c.Query("number")   // block
	if c.Query("block") != "" {      //还能用block查
		req.Number = c.Query("block")
	}
	req.Total = mysql.ConvertStringToInt64(c.Query("total"), 0) // 分页查询传入上次总计结果
	log.Debugf("Hello /api/transaction?%#v", req)
	//resp, err := service.QueryTransactions(req)
	resp, err := service.QueryTransactionsBuffer(req)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// QueryTransactionByHash ...
// @Summary QueryTransactionByHash ...
// @Description Query transaction by hash
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param hash query string false "hash"
// @Success 200 {string} json "{}"
// @Router /api/transaction/:hash [get]
func QueryTransactionByHash(c *gin.Context) {
	req := &entity.Transactions{}
	req.Hash = c.Param("hash") //占位符传参
	log.Debugf("Hello /api/transaction/:%#v", req.Hash)
	resp, err := service.QueryTransactionByHashFromBuffer(req)
	if resp == nil {
		resp, err = service.QueryTransaction(req)
	}
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}

// PostTransaction ...
// @Summary PostTransaction ...
// @Description post transaction
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param dry-run query string false "dry-run"
// @Success 200 {string} json "{}"
// @Router /api/transaction [post]
func PostTransaction(c *gin.Context) {
	dryRun := c.Query("dry-run")
	req := &entity.PostTransaction{}
	if c.BindJSON(req) == nil {
		if req == nil {
			log.Errorf("parsing request parameter err!")
			c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
		}
	}
	log.Debugf("Hello /api/transaction:dryRun:[%v]", dryRun)
	resp, err := service.PostTransaction(req, dryRun)
	if err != nil {
		errCode, _ := util.GetErrorCode(err)
		c.JSON(errCode, err)
	}
	c.JSON(http.StatusOK, resp)
}
