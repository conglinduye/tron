package util

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/wlcy/tron/explorer/lib/log"
)

//ErrorMsgShowTranceMessage ErrorMsg开关，是否显示代码调用跟踪信息
var ErrorMsgShowTranceMessage = true

//ErrorMsg 自定义的错误类型，包括了错误代码 和 错误消息
type ErrorMsg struct {
	errorCode    int
	errorMessage error
	fileName     string //受 ErrorMsgShowTranceMessage 控制
	codeLine     int    //受 ErrorMsgShowTranceMessage 控制
	functionName string //受 ErrorMsgShowTranceMessage 控制
}

//Error,the inteface of golang buildin type error
func (err *ErrorMsg) Error() string {
	return err.buildMessage()
}

//String 返回ERROR的描述信息
func (err *ErrorMsg) String() string {
	return err.buildMessage()
}

//ErrorCode 返回ERROR的编码
func (err *ErrorMsg) ErrorCode() int {
	return err.errorCode
}

//ErrorMessage 返回ERROR的消息
func (err *ErrorMsg) ErrorMessage() string {
	if nil != err.errorMessage {
		return err.errorMessage.Error()
	}
	return ""
}

//buildMessage 构造ERROR的消息描述
func (err *ErrorMsg) buildMessage() string {
	if false == ErrorMsgShowTranceMessage {
		return fmt.Sprintf("ErrorCode:[%d],ErrorMsg:[%s]", err.errorCode, err.errorMessage)
	}
	return fmt.Sprintf("ErrorCode:[%d],ErrorMsg:[%s],FileName:[%s],CodeLine:[%d],functionName[%s]",
		err.errorCode, err.errorMessage, err.fileName, err.codeLine, err.functionName)
}

//NewError 创建一个ErrorMsg
func NewError(errCode int, errMsg string) *ErrorMsg {
	return newErrorMessage(errCode, errMsg, 2)
}

//NewErrorMsg 创建一个ErrorMsg
func NewErrorMsg(errCode int) *ErrorMsg {
	errMsg := GetErrorMsgSleek(errCode)
	return newErrorMessage(errCode, errMsg, 2)
}

//GetErrorCode 返回ErrorCode
func GetErrorCode(err error) (int, bool) {
	if err != nil {
		//var errPtr *ErrorMsg
		errPtr, ok := err.(*ErrorMsg)
		if ok && nil != errPtr {
			return errPtr.ErrorCode(), true
		}
	}
	return -1, false
}

//GetErrorMessage 返回ERROR的消息
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	//var errPtr *ErrorMsg
	errPtr, ok := err.(*ErrorMsg)
	if ok && nil != errPtr {
		return errPtr.ErrorMessage()
	}
	return err.Error()
}

//newErrorMessage 构造错误对象
func newErrorMessage(errCode int, errMsg string, callStack int) *ErrorMsg {
	errmsg := new(ErrorMsg)
	errmsg.errorCode = errCode
	errmsg.errorMessage = errors.New(errMsg)

	//捕获错误，保证程序正常结束
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("recover error in ErrorMsg.NewError:", err)
		}
	}()

	//构造函数调用信息
	if ErrorMsgShowTranceMessage {
		pc, file, line, ok := runtime.Caller(callStack)
		if ok {
			errmsg.fileName = file
			errmsg.codeLine = line
			fun := runtime.FuncForPC(pc)
			if nil != fun {
				errmsg.functionName = fun.Name()
			}
		}
	}

	return errmsg
}
