package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/errno"
	"github.com/wlcy/tron/explorer/web/handler"
	"encoding/json"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/service"
	"strconv"
)

func QueryTokenBlackList(c *gin.Context) {
	req := &entity.AssetBlacklistReq{}
	req.Start = c.Query("start")
	req.Limit = c.Query("limit")
	req.OwnerAddress = c.Query("ownerAddress")
	req.TokenName = c.Query("tokenName")
	log.Infof("QueryTokenBlackList req:%#v", req)
	if req.Start == "" || req.Limit == "" {
		req.Start = "0"
		req.Limit = "20"
	}

	assetBlacklistResp, err := service.QueryAssetBlacklist(req)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}
	handler.SendResponse(c, nil, assetBlacklistResp)
}

func AddTokenBlackList(c *gin.Context) {
	var r entity.AssetBlacklist
	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	if r.OwnerAddress == "" || r.TokenName == "" {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	str, _ := json.Marshal(r)
	log.Infof("Create msg: %s\n", str)

	if err := service.InsertAssetBlacklist(r.OwnerAddress, r.TokenName); err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}

func DeleteTokenBlackList(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	log.Infof("Delete TokenBlackList Id: %d\n", userId)
	if err := service.DeleteAssetBlacklist(uint64(userId)); err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}
