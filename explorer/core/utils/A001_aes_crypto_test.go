package utils

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
)

func TestAES(*testing.T) {
	passwd := "this iS thE paSswoRD"
	salt, _ := uuid.NewV4()

	msg := "AES crypto data has the same size as the original data in bytes"

	cryptoData := AesCtrEncrypt(passwd, salt.String(), []byte(msg))

	decryptData := AesCtrDecrypt(passwd, salt.String(), cryptoData)

	fmt.Printf("%v\n%s\n", msg, decryptData)
}
