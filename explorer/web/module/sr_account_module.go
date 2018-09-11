package module

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
)

//QuerySrAccounts
func QuerySrAccountsRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.SrAccountResp, error){
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QuerySrAccounts error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QuerySrAccounts dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	srAccountResp := &entity.SrAccountResp{}
	srAccounts := make([]*entity.SrAccountInfo, 0)

	for dataPtr.NextT() {
		srAccount := &entity.SrAccountInfo{}
		srAccount.Address = dataPtr.GetField("address")
		srAccount.GithubLink = dataPtr.GetField("github_link")
		srAccount.CreateTime = dataPtr.GetField("create_time")
		srAccount.ModifiedTime = dataPtr.GetField("modified_time")

		srAccounts = append(srAccounts, srAccount)
	}

	var total = int64(len(srAccounts))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	srAccountResp.Total = total
	srAccountResp.Data = srAccounts

	return srAccountResp, nil
}
