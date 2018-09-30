package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

func GetApiUser(username string) (*entity.ApiUsers, error) {
	return module.GetApiUser(username)
}