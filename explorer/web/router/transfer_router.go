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

func transferRegister(ginRouter *gin.Engine) {

	//?sort=-number&limit=1&count=true&number=2135998
	ginRouter.GET("/api/transfer", func(c *gin.Context) {
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
	})
	//:number=2135998
	ginRouter.GET("/api/transfer/:hash", func(c *gin.Context) {
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
	})

}
