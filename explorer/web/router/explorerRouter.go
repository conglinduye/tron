package router

import (
	"net/http"
	"time"

	"github.com/wlcy/tron/explorer/web/service"

	"github.com/wlcy/tron/explorer/web/entity"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
)

//Start  启动服务
func Start(address string, objectpool int) {
	router := gin.Default()
	//?sort=-number&limit=1&count=true&number=2135998
	router.GET("/api/block", func(c *gin.Context) {
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
	service := http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Debugf("Start service, address:[%v],", address)

	service.ListenAndServe()

}
