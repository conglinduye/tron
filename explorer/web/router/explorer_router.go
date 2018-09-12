package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/lib/log"
)

//Start  启动服务
func Start(address string, objectpool int) {
	ginRouter := gin.Default()
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
