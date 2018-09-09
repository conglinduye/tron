package util

import (
	"fmt"
	"strconv"
	"sync"
)

//
// 此文件中用于定义常用的消息字符串
//

const Error_code_subsystem_step = 1000
const (
	Error_code_module_common  = 0 * Error_code_subsystem_step
	Error_code_module_user    = 10 * Error_code_subsystem_step
	Error_code_module_network = 20 * Error_code_subsystem_step

	Error_code_module_base = 88 * Error_code_subsystem_step
)

// 错误码以及错误信息的对应字典
var errorMessageMap map[int]string

//消息的锁
var errorMessageLock sync.Mutex

/*
//定义返回给客户端的操作结果标记
const (
	Response_json_status_success = "0"
	Response_json_status_error   = "1"
)*/

//定义返回给客户端的操作结果消息字符
const (
	Error_common_failure                    = Error_code_module_common + 1
	Error_common_no_data                    = Error_code_module_common + 2
	Error_common_db_not_connected           = Error_code_module_common + 3
	Error_common_internal_error             = Error_code_module_common + 4
	Error_common_parameter_invalid          = Error_code_module_common + 5
	Error_common_json_object_nil            = Error_code_module_common + 6
	Error_common_request_json_convert_error = Error_code_module_common + 7
	Error_common_request_json_no_data       = Error_code_module_common + 8
	Error_common_request_url_not_suport     = Error_code_module_common + 9
	Error_common_data_exist                 = Error_code_module_common + 10
	Error_common_add_exist_data             = Error_code_module_common + 11
	Error_common_organization_different     = Error_code_module_common + 12
	Error_common_object_name_duplicate      = Error_code_module_common + 13
	Error_common_not_suport_parameter       = Error_code_module_common + 14
	Error_common_not_suport_request_url     = Error_code_module_common + 15
	Error_common_send_mail                  = Error_code_module_common + 16
	Error_common_send_sms                   = Error_code_module_common + 17

	Error_user_token_invalid  = Error_code_module_user + 1
	Error_user_object_empty   = Error_code_module_user + 2
	Error_user_object_invalid = Error_code_module_user + 3
	Error_user_object_exist   = Error_code_module_user + 4
	Error_user_passwd_error   = Error_code_module_user + 5
	Error_user_passwd_invalid = Error_code_module_user + 6
	Error_user_role_exist     = Error_code_module_user + 7

	Error_redis_not_connected    = Error_code_module_network + 1
	Error_redis_data_invalid     = Error_code_module_network + 2
	Error_redis_read_write_error = Error_code_module_network + 3
)

// 根据错误ID返回错误信息
func GetErrorMsg(errCode int) (string, error) {
	ret, ok := errorMessageMap[errCode]
	if !ok {
		return "", fmt.Errorf("Unkown error code:[%v]", strconv.Itoa(errCode))
	}
	return ret, nil
}

//GetErrorMsgSleek 根据错误ID返回错误信息
func GetErrorMsgSleek(errCode int) string {
	strMsg, _ := GetErrorMsg(errCode)
	return strMsg
}

//GetErrorMessages 返回所有的消息
func GetErrorMessages() map[int]string {
	return errorMessageMap
}
