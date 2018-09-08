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
