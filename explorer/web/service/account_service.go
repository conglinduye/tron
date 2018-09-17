package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/buffer"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryAccounts 条件查询  	//?sort=-number&limit=1&count=true&number=2135998  TODO  cache
func QueryAccounts(req *entity.Accounts) (*entity.AccountsResp, error) {
	var filterSQL, sortSQL, pageSQL string
	/*strSQL := fmt.Sprintf(`
		   select account_name,acc.address,acc.balance as totalBalance,
		   frozen,create_time,latest_operation_time,votes ,
	       ass.asset_name as token_name,ass.creator_address,ass.balance
	       from tron.tron_account acc
	       left join tron.account_asset_balance ass on ass.address=acc.address
		   where 1=1 `)
	*/
	strSQL := fmt.Sprintf(`
		   select account_name,address,balance as totalBalance,
		   frozen,create_time,latest_operation_time,votes
	       from tron.tron_account acc
		   where 1=1 `)
	//按传入条件拼接sql，很容易错误，需要注意
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and acc.address='%v'", req.Address)
	}

	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "balance") > 0 {
			sortSQL = fmt.Sprintf("%v acc.balance", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
		}
	}
	if sortSQL != "" {
		if strings.Index(sortSQL, ",") == 0 {
			sortSQL = sortSQL[1:]
		}
		sortSQL = fmt.Sprintf("order by %v", sortSQL)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountsRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountsRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	accountsResp := &entity.AccountsResp{}
	accountInfos := make([]*entity.AccountInfo, 0)
	//accountListMap := make(map[string]*entity.AccountInfo, 0) //保存每个账户信息，用于去重
	var oldBalance = make([]*entity.BalanceInfoDB, 0) //解析冻结信息
	var totalFrozen = int64(0)                        //power信息

	//填充数据
	//ass.asset_name as token_name,ass.creator_address,ass.balance
	for dataPtr.NextT() {
		var account = &entity.AccountInfo{}
		frozen := dataPtr.GetField("frozen")
		if frozen != "" {
			if err := json.Unmarshal([]byte(frozen), &oldBalance); err != nil {
				log.Errorf("Unmarshal data failed:[%v]-[%v]", err, frozen)
				return nil, util.NewErrorMsg(util.Error_common_request_json_convert_error)
			}
		}
		for _, blanceFrozen := range oldBalance {
			totalFrozen += blanceFrozen.Amount
		}
		account.Power = totalFrozen
		account.Address = dataPtr.GetField("address")
		account.CreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		account.UpdateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_operation_time"))
		account.Name = dataPtr.GetField("account_name")
		account.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBalance"))
		/*tokenInfo, err := querytokenBalanceInfo(account.Address)
		if err != nil {
			log.Errorf("get token balance info err:[%v] by adderss:[%v]", err, account.Address)
			return nil, err
		}*/
		//从缓存中获取数据
		account.TokenBalances = buffer.GetAccountTokenBuffer().GetAccountTokenBuffer(account.Address)
		accountInfos = append(accountInfos, account)
	}

	//查询该语句所查到的数据集合
	var total = int64(len(accountInfos))
	total, err = mysql.QueryTableDataCount("tron.tron_account")
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	accountsResp.Total = total
	accountsResp.Data = accountInfos
	return accountsResp, nil

	//return module.QueryAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryAccount 精确查询  	//number=2135998   TODO  cache
func QueryAccount(req *entity.Accounts) (*entity.AccountDetail, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
	select account_name,acc.address,acc.balance as totalBalance,frozen,create_time,latest_operation_time,votes ,
		wit.url,wit.is_job,acc.allowance,acc.latest_withdraw_time,
		acc.net_usage,acc.free_net_limit,acc.net_used,acc.net_limit,acc.asset_net_used,acc.asset_net_limit,
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
