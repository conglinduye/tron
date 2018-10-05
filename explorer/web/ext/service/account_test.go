package service

import (
	"fmt"
	"testing"
)

//{Key:"cc9fd97198c6072729fa5df0159d5607d4e3da03a92d0c24eb89a9f07f43539d", Address:"TCWkTNVAV5QErki3ERxkYsMjkzWWXaQCB6"}
func TestCreateAccount(t *testing.T) {
	ss, _ := CreateAccount()
	fmt.Printf("%#v", ss)
}
