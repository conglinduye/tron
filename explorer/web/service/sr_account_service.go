package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
)

//QuerySrAccounts
func QuerySrAccounts(req *entity.SrAccount) (*entity.SrAccountResp, error) {
	var filterSQL, sortSQL, pageSQL string

	strSQL := fmt.Sprintf(`
			select address, github_link, create_time, modified_time
			from wlcy_sr_account
			where 1=1 `)

	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", req.Address)
	}
	if req.Limit != "" && req.Start != "" {
		pageSQL = fmt.Sprintf(" limit %v, %v", req.Start, req.Limit)
	}

	return module.QuerySrAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL)
}
