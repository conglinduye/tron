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

func transactionRegister(ginRouter *gin.Engine) {

	//?sort=-number&limit=1&count=true&number=2135998
	ginRouter.GET("/api/transaction", func(c *gin.Context) {
		req := &entity.Transactions{}
		req.Sort = c.Query("sort")
		req.Limit = mysql.ConvertStringToInt64(c.Query("limit"), 40)
		req.Count = c.Query("count")
		req.Start = mysql.ConvertStringToInt64(c.Query("start"), 0)
		req.Hash = c.Query("hash")
		req.Number = c.Query("number")
		if c.Query("block") != "" { //还能用block查
			req.Number = c.Query("block")
		}
		log.Debugf("Hello /api/transaction?%#v", req)
		resp, err := service.QueryTransactions(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	//:number=2135998
	ginRouter.GET("/api/transaction/:hash", func(c *gin.Context) {
		req := &entity.Transactions{}
		req.Hash = c.Param("hash") //占位符传参
		log.Debugf("Hello /api/transaction/:%#v", req.Hash)
		resp, err := service.QueryTransaction(req)
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
	/*
		ginRouter.POST("/api/transaction", func(c *gin.Context) {
			req := &entity.PostTransaction{}
			if c.BindJSON(req) == nil {
				if req == nil {
					log.Errorf("parsing request parameter err!")
					c.JSON(http.StatusInternalServerError, http.ErrBodyNotAllowed)
				}
			}
			log.Debugf("Hello /api/transaction")
			resp, err := service.PostTransaction(req)
			if err != nil {
				errCode, _ := util.GetErrorCode(err)
				c.JSON(errCode, err)
			}
			c.JSON(http.StatusOK, resp)
		})
	*/
}
