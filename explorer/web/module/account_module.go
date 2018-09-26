package module

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryAccountsRealize 操作数据库
func QueryAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.AccountsResp, error) {
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Sql(strFullSQL)
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
		//account.TokenBalances = buffer.GetAccountTokenBuffer().GetAccountTokenBuffer(account.address)

		//account.TokenBalances = tokenInfo
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

}

//查询某个地址下的token信息
func querytokenBalanceInfo(address string) (map[string]int64, error) {
	strSQL := fmt.Sprintf(`
	select acc.address,acc.asset_name as token_name,acc.creator_address,acc.balance
	from tron.account_asset_balance acc
	where 1=1 order by `)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("querytokenBalanceInfo error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("querytokenBalanceInfo dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	accountTokenMap := make(map[string]int64, 0)
	for dataPtr.NextT() {
		tokenName := dataPtr.GetField("token_name")
		balance := mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		if tokenName != "" {
			if _, ok := accountTokenMap[tokenName]; !ok {
				accountTokenMap[tokenName] = balance
			}
		}
	}
	return accountTokenMap, nil
}

//QueryAccountRealize 操作数据库
func QueryAccountRealize(strSQL, filterSQL, address string) (*entity.AccountDetail, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
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
	var frozenInfo = &entity.Frozen{Total: 0, Balances: apiBalance}
	var represent = &entity.Represent{}
	//var balance = &entity.Balance{}
	var balances = make([]*entity.Balance, 0)
	var bandwidth = &entity.BandwidthInfo{}
	var totalFrozen = int64(0)

	accountTokenMap := make(map[string][]*entity.Balance, 0) //保存每个账户的token信息
	if dataPtr.ResNum() == 0 {
		account = &entity.AccountDetail{
			Representative: represent,
			Name:           "",
			Address:        address,
			Bandwidth:      bandwidth,
			Balances:       balances,
			Balance:        0,
			TokenBalances:  balances,
			Frozen:         frozenInfo,
		}
	} else {
		//填充数据
		for dataPtr.NextT() {
			if account.Address == "" {
				account.Address = dataPtr.GetField("address")
				account.Name = dataPtr.GetField("account_name")
				account.Balance = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBalance"))
				represent.Allowance = mysql.ConvertDBValueToInt64(dataPtr.GetField("allowance"))
				isWitness := dataPtr.GetField("is_witness")
				if isWitness == "1" {
					represent.Enabled = true
				}
				represent.LastWithDrawTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("latest_withdraw_time"))
				represent.URL = dataPtr.GetField("url")
				account.Representative = represent
				//acc.free_net_used,acc.free_net_limit,acc.net_used,acc.net_limit,acc.asset_net_used,acc.asset_net_limit,
				bandwidth.FreeNetUsed = mysql.ConvertDBValueToInt64(dataPtr.GetField("free_net_used"))
				bandwidth.FreeNetLimit = mysql.ConvertDBValueToInt64(dataPtr.GetField("free_net_limit"))
				bandwidth.FreeNetRemaining = bandwidth.FreeNetLimit - bandwidth.FreeNetUsed
				bandwidth.FreeNetPercentage = 0
				if bandwidth.FreeNetLimit > 0 {
					bandwidth.FreeNetPercentage = float64(bandwidth.FreeNetUsed) / float64(bandwidth.FreeNetLimit)
				}
				bandwidth.NetUsed = mysql.ConvertDBValueToInt64(dataPtr.GetField("net_used"))
				bandwidth.NetLimit = mysql.ConvertDBValueToInt64(dataPtr.GetField("net_limit"))
				bandwidth.NetRemaining = bandwidth.NetLimit - bandwidth.NetUsed
				bandwidth.NetPercentage = 0
				if bandwidth.NetLimit > 0 {
					bandwidth.NetPercentage = float64(bandwidth.NetUsed) / float64(bandwidth.NetLimit)
				}
				assetNetUsed := dataPtr.GetField("asset_net_used")
				assetNetLimit := dataPtr.GetField("asset_net_limit")
				bandwidth.Assets = getAssetNetInfo(assetNetUsed, assetNetLimit)
				account.Bandwidth = bandwidth
				//[{"frozen_balance":4306000000,"expire_time":1534794417000}]
				frozen := dataPtr.GetField("frozen")
				if frozen != "" {
					if err := json.Unmarshal([]byte(frozen), &oldBalance); err != nil {
						log.Errorf("Unmarshal data failed:[%v]-[%v]", err, frozen)
					}
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
				account.Frozen = frozenInfo
			}
			balance := &entity.Balance{}
			balance.Name = dataPtr.GetField("token_name")
			balance.Address = dataPtr.GetField("owner_address")
			if strings.ToUpper(balance.Name) == "TRX" { //如果通证是波场币，地址替换为当前持有者账户的地址
				balance.Address = account.Address
			}
			balance.Balance = mysql.ConvertDBValueToFloat64(dataPtr.GetField("balance"))

			if account.Address != "" {
				if tokenInfo, ok := accountTokenMap[account.Address]; ok {
					if balance.Name != "" {
						tokenInfo = append(tokenInfo, balance)
						accountTokenMap[account.Address] = tokenInfo
					}
				} else {
					tokenArr := make([]*entity.Balance, 0)
					if account.Balance > 0 {
						ownbalance := &entity.Balance{}
						ownbalance.Name = "TRX"
						balance.Address = account.Address
						ownbalance.Balance = float64(account.Balance) / 1000000 //单位换算，页面按照TRX显示
						tokenArr = append(tokenArr, ownbalance)
					}
					if balance.Balance > 0 {
						tokenArr = append(tokenArr, balance)
					}

					accountTokenMap[account.Address] = tokenArr
				}
			}
		}

		//拼接tokeninfo列表
		if tokenInfo, ok := accountTokenMap[account.Address]; ok {
			account.TokenBalances = tokenInfo
			account.Balances = tokenInfo
		}
	}
	return account, nil

}

//QueryAccountMediaRealize 操作数据库
func QueryAccountMediaRealize(strSQL, filterSQL string) (*entity.AccountMediaInfo, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
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
		mediaInfo.Image = dataPtr.GetField("url") //TODO 需要从表中配置的url里面扣logo图片
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
	log.Sql(strFullSQL)
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
	strSQL := fmt.Sprintf(`insert into tron.wlcy_sr_account
			(address,github_link) value( '%v','%v')`,
		address, github)

	log.Sql(strSQL)
	instID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("InsertSrAccount result fail:[%v]  sql:%s", err, strSQL)
		return instID, err
	}
	return instID, err
}

//UpdateSrAccount 更新github地址
func UpdateSrAccount(address, github string) (int64, error) {
	strSQL := fmt.Sprintf(`update tron.wlcy_sr_account set github_link='%v' where address='%v'`,
		github, address)

	log.Sql(strSQL)
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
	log.Sql(strFullSQL)
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

//QueryAccountStatsRealize 查询用户的交易统计信息
func QueryAccountStatsRealize(strSQL string) (*entity.AccountTransactionNum, error) {
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("QueryAccountStatsRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountStatsRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var acountTrxInfo = &entity.AccountTransactionNum{}
	//填充数据
	for dataPtr.NextT() {
		acountTrxInfo.TransactionsOut = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxOut"))
		acountTrxInfo.TransactionIn = mysql.ConvertDBValueToInt64(dataPtr.GetField("trxIn"))
		acountTrxInfo.Transactions = acountTrxInfo.TransactionsOut + acountTrxInfo.TransactionIn
	}

	return acountTrxInfo, nil
}

//解析用户asset的带宽使用情况
func getAssetNetInfo(assetNetUsed, assetNetLimit string) map[string]*entity.AssetInfo {
	var assetNetInfo = make(map[string]*entity.AssetInfo, 0)
	if assetNetUsed == "" || assetNetLimit == "" {
		return assetNetInfo
	}
	netUsedMap := util.ParsingJSONFromString(assetNetUsed)
	netLimitMap := util.ParsingJSONFromString(assetNetLimit)
	for param, value := range netUsedMap {
		if param != "" {
			assetInfo := &entity.AssetInfo{}
			assetInfo.NetUsed = value
			if val, ok := netLimitMap[param]; ok {
				assetInfo.NetLimit = val
			} else {
				assetInfo.NetLimit = 0
			}
			assetInfo.NetRemaining = assetInfo.NetLimit - assetInfo.NetUsed
			assetInfo.NetPercentage = 0
			if assetInfo.NetLimit > 0 {
				assetInfo.NetPercentage = float64(assetInfo.NetUsed) / float64(assetInfo.NetLimit)
			}
			assetNetInfo[param] = assetInfo

		}
	}
	return assetNetInfo
}
