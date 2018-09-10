package utils

import (
	"encoding/json"
	"fmt"
)

func jsonPrint(val interface{}) {
	data, err := json.Marshal(val)
	fmt.Printf("err:%v\n%s\n----------------------------------\n\n", err, data)
}

// VerifyCall ...
func VerifyCall(val interface{}, err error) {
	if nil != err || nil == val {
		fmt.Printf("Faield, error:%v\n", err)
	} else {
		jsonPrint(val)
	}
}

// ToJSONStr ...
func ToJSONStr(val interface{}) string {
	if nil == val {
		return ""
	}
	data, err := json.Marshal(val)
	if nil != err {
		return fmt.Sprintf("%#v", val)
	}
	return string(data)
}
