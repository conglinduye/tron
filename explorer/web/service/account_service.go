package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryAccounts 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryAccounts(req *entity.Accounts) (*entity.AccountsResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string

	strSQL := fmt.Sprintf(`
		   select account_name,acc.address,acc.balance as totalBalance,
		   frozen,create_time,latest_operation_time,votes ,
	       ass.asset_name as token_name,ass.creator_address,ass.balance
	       from tron.tron_account acc
	       left join tron.account_asset_balance ass on ass.address=acc.address
		   where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and acc.address='%v'", req.Address)
	}

	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "balance") > 0 {
			sortTemp = fmt.Sprintf("%v acc.balance", sortTemp)
			if strings.Index(v, "-") == 0 {
				sortTemp = fmt.Sprintf("%v desc", sortTemp)
			}
		}
	}
	if sortTemp != "" {
		if strings.Index(sortTemp, ",") == 0 {
			sortTemp = sortTemp[1:]
		}
		sortTemp = fmt.Sprintf("order by %v", sortTemp)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryAccount 精确查询  	//number=2135998
func QueryAccount(req *entity.Accounts) (*entity.AccountDetail, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
	select account_name,acc.address,acc.balance as totalBalance,frozen,create_time,latest_operation_time,votes ,
        wit.url,wit.is_job,acc.allowance,acc.latest_withdraw_time,
        ass.asset_name as token_name,ass.creator_address,ass.balance
    from tron.tron_account acc
    left join tron.account_asset_balance ass on ass.address=acc.address
    left join tron.witness wit on wit.address=acc.address
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and acc.address='%v'", req.Address)
	}
	return module.QueryAccountRealize(strSQL, filterSQL)
}

//QueryAccountMedia 查询账户媒体信息 	//number=2135998
func QueryAccountMedia(req *entity.Accounts) (*entity.AccountMediaInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
	select address,url
	from tron.wlcy_witness_create_info
	where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", req.Address)
	}
	return module.QueryAccountMediaRealize(strSQL, filterSQL)
}

//UpdateAccountSr 更新超级账户github信息 	//number=2135998
func UpdateAccountSr(req *entity.SuperAccountInfo) (*entity.SuperAccountInfo, error) {
	var id int64
	var err error
	var srAccount = &entity.SuperAccountInfo{}
	if module.CheckSrAccountExist(req.Address) { // 更新
		id, err = module.UpdateSrAccount(req.Address, req.GithubLink)
	} else { // 插入
		id, err = module.InsertSrAccount(req.Address, req.GithubLink)
	}
	if err == nil {
		srAccount.ID = id
		srAccount.Address = req.Address
		srAccount.GithubLink = req.GithubLink
	}

	return srAccount, err
}

//QueryAccountSr 查询超级账户github信息 	//number=2135998
func QueryAccountSr(req *entity.SuperAccountInfo) (*entity.SuperAccountInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select address,url from tron.wlcy_witness_create_info
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", req.Address)
	}
	return module.QueryAccountSrRealize(strSQL, filterSQL)
}
