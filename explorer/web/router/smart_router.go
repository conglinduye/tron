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

func smartRegister(ginRouter *gin.Engine) {

	//查询智能合约
	ginRouter.GET("/api/contracts", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Sort = c.Query("sort")
		req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		req.Count = c.Query("count")
		req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		log.Debugf("Hello /api/contracts?%#v", req)
		resp, err := service.QueryContracts(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	//合约详情信息
	ginRouter.GET("/api/contract/:address", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Address = c.Param("address") //占位符传参
		log.Debugf("Hello /api/transfer/:%#v", req.Address)
		resp, err := service.QueryContractByAddress(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//查询智能合约交易列表
	ginRouter.GET("/api/contracts/transaction", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Sort = c.Query("sort")
		req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		req.Count = c.Query("count")
		req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		req.Address = c.Query("contract")
		req.Type = c.Query("type")
		log.Debugf("Hello /api/contracts/transaction?%#v", req)
		resp, err := service.QueryContractTnx(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	// 查询智能合约code信息
	ginRouter.GET("/api/contracts/code", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Address = c.Query("contract")
		log.Debugf("Hello /api/contracts/code?%#v", req)
		resp, err := service.QueryContractsCode(req)

		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

	//  查询智能合约event信息
	/*ginRouter.GET("/api/contracts/event", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Address = c.Query("contract")
		log.Debugf("Hello /api/contracts/event?%#v", req)
		resp, err := service.QueryContractEvent(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	// 查询智能合约内部交易
	ginRouter.GET("/api/contracts/internalTxs", func(c *gin.Context) {
		req := &entity.Contracts{}
		req.Sort = c.Query("sort")
		req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		req.Count = c.Query("count")
		req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		req.Address = c.Query("contract")
		log.Debugf("Hello /api/contracts/internalTxs?%#v", req)
		resp, err := service.QueryInternalTxs(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			cJSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})*/
	// 校验智能合约
	ginRouter.POST("/api/contracts/verify", func(c *gin.Context) {
		req := &entity.ContractCodeInfo{}
		if c.BindJSON(req) == nil {
			if req == nil {
				log.Errorf("parsing request parameter err!")
				c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
			}
		}
		log.Debugf("Hello /api/contracts/verify?%#v", req)
		resp, err := service.VerifyContractCode(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})

}
