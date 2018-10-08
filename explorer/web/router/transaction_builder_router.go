package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/service"
)

func registerTransactionBuilderRouter(router *gin.Engine) {

	router.GET("/api/transaction-builder/contract/transfer", service.TBTransfer)
	router.GET("/api/transaction-builder/contract/transferasset", service.TBTransferAsset)
	router.GET("/api/transaction-builder/contract/accountcreate", service.TBAccountCreate)
	router.GET("/api/transaction-builder/contract/accountupdate", service.TBAccountUpdate)
	router.GET("/api/transaction-builder/contract/withdrawbalance", service.TBWithdrawBalance)

	router.POST("/api/transaction-builder/contract/transfer", service.TBTransfer)
	router.POST("/api/transaction-builder/contract/transferasset", service.TBTransferAsset)
	router.POST("/api/transaction-builder/contract/accountcreate", service.TBAccountCreate)
	router.POST("/api/transaction-builder/contract/accountupdate", service.TBAccountUpdate)
	router.POST("/api/transaction-builder/contract/withdrawbalance", service.TBWithdrawBalance)
}
