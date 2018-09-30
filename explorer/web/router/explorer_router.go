package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/router/middleware"
)

//Start  启动服务
func Start(address string, objectpool int) {
	ginRouter := gin.Default()
	ginRouter.Use(corsMiddleware())
	// 注册区块链查询路由
	blockRegister(ginRouter)
	// 注册交易查询路由
	transactionRegister(ginRouter)
	// 注册转账查询路由
	transferRegister(ginRouter)
	// 注册账户查询路由
	accountRegister(ginRouter)
	// 注册投票查询路由
	voteRegister(ginRouter)
	// 注册超级代表查询路由
	witnessRegister(ginRouter)
	// 注册通证查询路由
	tokenRegister(ginRouter)
	// 注册统计查询路由
	reportRegister(ginRouter)
	// 注册其他查询路由
	otherRegister(ginRouter)

	//ginRouter.Use(cors.Default())

	ginRouter.POST("/api/login", Login)

	tokenBlacklist := ginRouter.Group("/api/tokenBlacklist")
	// 授权处理
	tokenBlacklist.Use(middleware.AuthMiddleware())
	{
		tokenBlacklist.POST("/add", AddTokenBlackList)
		tokenBlacklist.DELETE("/delete/:id", DeleteTokenBlackList)
		tokenBlacklist.GET("/list", QueryTokenBlackList)
	}

	tokenExt := ginRouter.Group("/api/tokenExt")
	// 授权处理
	tokenExt.Use(middleware.AuthMiddleware())
	{
		tokenExt.POST("/addInfo", AddAssetExtInfo)
		tokenExt.POST("/updateInfo", UpdateAssetExtInfo)
		tokenExt.POST("/addLogo", AddAssetExtLogo)
		tokenExt.POST("/updateLogo", UpdateAssetExtLogo)
	}


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
