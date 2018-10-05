package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/ext/service"
)

func accountRegister(ginRouter *gin.Engine) {
	//创建新地址和private key，不保存数据库
	ginRouter.POST("/api/account", func(c *gin.Context) {
		log.Debugf("Hello /api/account POST")
		resp, err := service.CreateAccount()
		if err != nil {
			errCode, _ := util.GetErrorCode(err)
			c.JSON(errCode, err)
		}
		c.JSON(http.StatusOK, resp)
	})
}
