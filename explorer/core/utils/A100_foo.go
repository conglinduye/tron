package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	real := reflect.ValueOf(val)
	if real.IsNil() {
		return ""
	}
	if real.Kind() == reflect.Ptr && !real.Elem().IsValid() {
		return ""
	}
	if (real.Kind() == reflect.Slice || real.Kind() == reflect.Array || real.Kind() == reflect.Map) && real.IsNil() {
		fmt.Printf("list:%#v\n", real)
		return ""
	}
	data, err := json.Marshal(val)
	if nil != err {
		return fmt.Sprintf("%#v", val)
	}
	return string(data)
}
