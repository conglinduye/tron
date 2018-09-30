package module

import (
	"fmt"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/entity"
)


func GetApiUser(username string) (*entity.ApiUsers, error) {
	strSQL := fmt.Sprintf(`
		select id, username, password
		from wlcy_api_users where username='%v'`, username)
	log.Sql(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("GetUser error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("GetUser dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}

	apiUsers := &entity.ApiUsers{}

	for dataPtr.NextT() {
		apiUsers.Id = mysql.ConvertDBValueToUint64(dataPtr.GetField("id"))
		apiUsers.Username = dataPtr.GetField("username")
		apiUsers.Password = dataPtr.GetField("password")
	}

	return apiUsers, nil
}

