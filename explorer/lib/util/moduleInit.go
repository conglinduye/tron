package util

//InitErrors 初始化ERROR code 和 消息
func init() {
	initErrorMsg()
}

//initErrorMsg 初始化错误码
func initErrorMsg() {
	errorMessageLock.Lock()
	defer errorMessageLock.Unlock()

	//重新申请一块新的内存
	errorMessageMap = make(map[int]string, 100)

	//common
	errorMessageMap[Error_common_failure] = "操作失败"
	errorMessageMap[Error_common_no_data] = "没有匹配的数据"
	errorMessageMap[Error_common_db_not_connected] = "数据库连接失败"
	errorMessageMap[Error_common_internal_error] = "系统内部计算或者逻辑错误"
	errorMessageMap[Error_common_parameter_invalid] = "参数无效"
	errorMessageMap[Error_common_json_object_nil] = "JSON 对象为空"
	errorMessageMap[Error_common_request_json_convert_error] = "请求参数无法转化为有效对象"
	errorMessageMap[Error_common_request_json_no_data] = "请求参数不完整或缺少请求参数"
	errorMessageMap[Error_common_request_url_not_suport] = "无法识别请求的URL或者不支持该URL"
	errorMessageMap[Error_common_data_exist] = "该对象已经存在"
	errorMessageMap[Error_common_add_exist_data] = "不可以添加已经存在的对象"
	errorMessageMap[Error_common_organization_different] = "不可以操作其他账户的数据"
	errorMessageMap[Error_common_object_name_duplicate] = "操作对象的名称不可以重复"
	errorMessageMap[Error_common_not_suport_parameter] = "不支持的请求参数"
	errorMessageMap[Error_common_not_suport_request_url] = "不支持该请求接口"
	errorMessageMap[Error_common_send_mail] = "发送邮件失败"
	errorMessageMap[Error_common_send_sms] = "发送短消息失败"

	//user
	errorMessageMap[Error_user_token_invalid] = "用户登录标识无效"
	errorMessageMap[Error_user_object_empty] = "用户信息为空"
	errorMessageMap[Error_user_object_invalid] = "用户信息中存在无效数据"
	errorMessageMap[Error_user_object_exist] = "该用户已经存在"
	errorMessageMap[Error_user_passwd_error] = "用户密码错误"
	errorMessageMap[Error_user_passwd_invalid] = "用户密码不符合安全策略"
	errorMessageMap[Error_user_role_exist] = "该角色已经存在"

	//redis
	errorMessageMap[Error_redis_not_connected] = "redis连接断开"
	errorMessageMap[Error_redis_data_invalid] = "redis中数据不合法数据"
	errorMessageMap[Error_redis_read_write_error] = "读写redis时出错"

}
