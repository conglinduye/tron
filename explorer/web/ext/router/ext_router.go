package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
)

// not used
//Start  启动服务
func Start(address string, objectpool int) {
	ginRouter := gin.Default()
	ginRouter.Use(corsMiddleware())
	// define your register
	//AccountRegister(ginRouter)

	service := http.Server{
		Addr:           address,
		Handler:        ginRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Debugf("Start service, address:[%v],", address)

	service.ListenAndServe()

}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var isAccess = true
		/*origin := c.Request.Header.Get("Origin")
		var filterHost = [...]string{"http://localhost.*", "http://*.hfjy.com"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = false
		for _, v := range filterHost {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}*/
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
