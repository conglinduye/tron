package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/handler"
	"github.com/wlcy/tron/explorer/web/errno"
	"github.com/wlcy/tron/explorer/web/service"
	"github.com/wlcy/tron/explorer/web/router/middleware"
)

func Login(c *gin.Context) {
	var u entity.ApiUsers
	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	d, err := service.GetApiUser(u.Username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if err := middleware.Compare(d.Password, u.Password); err != nil {
		handler.SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	token, err := middleware.Sign(c, middleware.Context{ID: d.Id, Username: d.Username}, "")
	if err != nil {
		handler.SendResponse(c, errno.ErrToken, nil)
		return
	}

	handler.SendResponse(c, nil, entity.ApiToken{Token: token})
}
