package module

import (
	"encoding/json"
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryAccountsRealize 操作数据库
func QueryAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AccountsResp, error) {
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
	accountTokenMap := make(map[string]map[string]int64, 0) //保存每个账户的token信息

	//填充数据
	for dataPtr.NextT() {
		var account = &entity.AccountInfo{}
		account.Power = dataPtr.GetField("votes") //TODO
		account.Address = dataPtr.GetField("address")
		account.CreateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("create_time"))
		account.UpdateTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_operation_time"))
		account.Name = dataPtr.GetField("account_name")
		account.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBalance"))
		tokenName := dataPtr.GetField("token_name")
		balance := mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		if account.Address != "" {
			if tokenInfo, ok := accountTokenMap[account.Address]; ok {
				tokenInfo[tokenName] = balance
				accountTokenMap[account.Address] = tokenInfo
			} else {
				tokenMap := make(map[string]int64, 0)
				tokenMap[tokenName] = balance
				accountTokenMap[account.Address] = tokenMap
			}
		}

		accountInfos = append(accountInfos, account)
	}

	//拼接tokeninfo列表
	for _, accountInfo := range accountInfos {
		if tokenInfo, ok := accountTokenMap[accountInfo.Address]; ok {
			accountInfo.TokenBalances = tokenInfo
		}
	}

	//查询该语句所查到的数据集合
	var total = int64(len(accountInfos))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	accountsResp.Total = total
	accountsResp.Data = accountInfos

	return accountsResp, nil

}

//QueryAccountRealize 操作数据库
func QueryAccountRealize(strSQL, filterSQL string) (*entity.AccountDetail, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var account = &entity.AccountDetail{}
	var oldBalance = make([]*entity.BalanceInfoDB, 0)
	var apiBalance = make([]*entity.BalanceInfo, 0)
	var frozenInfo = &entity.Frozen{}
	var represent = &entity.Represent{}
	var totalFrozen = int64(0)

	accountTokenMap := make(map[string][]*entity.Balance, 0) //保存每个账户的token信息

	// account.Bandwidth TODO
	//填充数据
	for dataPtr.NextT() {
		var balance = &entity.Balance{}
		if account.Address == "" {
			account.Address = dataPtr.GetField("address")
			account.Name = dataPtr.GetField("account_name")
			account.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBalance"))
			represent.Allowance = mysql.ConvertDBValueToInt64(dataPtr.GetField("allowance"))
			isJob := dataPtr.GetField("is_job")
			if isJob == "1" {
				represent.Enabled = true
			}
			represent.LastWithDrawTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_withdraw_time"))
			represent.URL = dataPtr.GetField("url")
			account.Representative = represent
			//[{"frozen_balance":4306000000,"expire_time":1534794417000}]
			frozen := dataPtr.GetField("frozen")
			if err := json.Unmarshal([]byte(frozen), oldBalance); err != nil {
				log.Errorf("Unmarshal data failed:[%v]", err)
			}
			for _, blanceFrozen := range oldBalance {
				apiFrozen := &entity.BalanceInfo{}
				apiFrozen.Amount = blanceFrozen.Amount
				apiFrozen.Expires = blanceFrozen.Expires
				apiBalance = append(apiBalance, apiFrozen)
				totalFrozen += blanceFrozen.Amount
			}
			frozenInfo.Total = totalFrozen
			frozenInfo.Balances = apiBalance
		}

		balance.Name = dataPtr.GetField("token_name")
		balance.Balance = mysql.ConvertDBValueToFloat64(dataPtr.GetField("balance"))

		if account.Address != "" {
			if tokenInfo, ok := accountTokenMap[account.Address]; ok {
				tokenInfo = append(tokenInfo, balance)
				accountTokenMap[account.Address] = tokenInfo
			} else {
				tokenArr := make([]*entity.Balance, 0)
				tokenArr = append(tokenArr, balance)
				accountTokenMap[account.Address] = tokenArr
			}
		}
	}

	//拼接tokeninfo列表
	if tokenInfo, ok := accountTokenMap[account.Address]; ok {
		account.TokenBalances = tokenInfo
		account.Balances = tokenInfo
	}

	return account, nil

}

//QueryAccountMediaRealize 操作数据库
func QueryAccountMediaRealize(strSQL, filterSQL string) (*entity.AccountMediaInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountMediaRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountMediaRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var mediaInfo = &entity.AccountMediaInfo{}
	//填充数据
	for dataPtr.NextT() {
		mediaInfo.Success = false
		mediaInfo.Image = dataPtr.GetField("url") //TODO
		mediaInfo.Reason = "Could not retrieve file"
	}

	return mediaInfo, nil

}

//CheckSrAccountExist 校验账户是否有github地址
func CheckSrAccountExist(address string) bool {
	var filterSQL string
	exist := false
	strSQL := fmt.Sprintf(`
	SELECT address,github_link 
	FROM tron.wlcy_sr_account
	where 1=1   `)

	//按传入条件拼接sql，很容易错误，需要注意
	if address != "" {
		filterSQL = fmt.Sprintf(" and address='%v'", address)
	}
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("CheckSrAccountExist error :[%v]\n", err)
		return exist
	}
	if dataPtr == nil {
		log.Errorf("CheckSrAccountExist dataPtr is nil ")
		return exist
	}
	//填充数据
	if dataPtr.ResNum() > 0 {
		exist = true
	}
	return exist

}

//InsertSrAccount 插入github地址
func InsertSrAccount(address, github string) (int64, error) {
	strSQL := fmt.Sprintf(`insert into tron.wlcy_witness_create_info
			(address,url) value( '%v','%v')`,
		address, github)

	log.Debugf(strSQL)
	instID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("InsertSrAccount result fail:[%v]  sql:%s", err, strSQL)
		return instID, err
	}
	return instID, err
}

//UpdateSrAccount 更新github地址
func UpdateSrAccount(address, github string) (int64, error) {
	strSQL := fmt.Sprintf(`update from tron.wlcy_witness_create_info set url='%v' where address='%v'`,
		github, address)

	log.Debugf(strSQL)
	instID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("UpdateSrAccount result fail:[%v]  sql:%s", err, strSQL)
		return instID, err
	}
	return instID, err
}

//QueryAccountSrRealize 按账户查询github信息
func QueryAccountSrRealize(strSQL, filterSQL string) (*entity.SuperAccountInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountSrRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountSrRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var srAcountInfo = &entity.SuperAccountInfo{}
	//填充数据
	for dataPtr.NextT() {
		srAcountInfo.Address = dataPtr.GetField("address")
		srAcountInfo.GithubLink = dataPtr.GetField("url")
	}

	return srAcountInfo, nil
}
