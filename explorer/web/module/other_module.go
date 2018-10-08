package module

import (
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/wlcy/tron/explorer/lib/mysql"
)

//FindByRecentIP 根据ip获取最近一小时记录，如存在，返回true
func FindByRecentIP(ip string) bool {
	var result bool
	now := time.Now()
	h, _ := time.ParseDuration("-1h")
	hour := now.Add(h).Format(mysql.DATETIMEFORMAT)
	strSQL := fmt.Sprintf(`select address,ip from wlcy_trx_request where ip='%v' and create_time>='%v'`, ip, hour)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("FindByRecentIP error :[%v]\n", err)
		return result
	}
	if dataPtr == nil {
		log.Errorf("FindByRecentIP dataPtr is nil ")
		return result
	}
	if dataPtr.ResNum() > 0 {
		result = true
	}
	return result
}

//FindByAddress 根据address获取记录，如存在，返回true
func FindByAddress(address string) bool {
	var result bool
	strSQL := fmt.Sprintf(`select address,ip from wlcy_trx_request where address='%v'`, address)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("FindByAddress error :[%v]\n", err)
		return result
	}
	if dataPtr == nil {
		log.Errorf("FindByAddress dataPtr is nil ")
		return result
	}
	if dataPtr.ResNum() > 0 {
		result = true
	}
	return result
}

//InsertTrxRequest ...
func InsertTrxRequest(address, ip string) error {
	strSQL := fmt.Sprintf(`
	insert into wlcy_trx_request 
	(ip, address)values('%v', '%v')`, ip, address)
	insID, _, err := mysql.ExecuteSQLCommand(strSQL, false)
	if err != nil {
		log.Errorf("insert TrxRequest  fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("insert TrxRequest success, insert id: [%v]", insID)
	return nil
}
