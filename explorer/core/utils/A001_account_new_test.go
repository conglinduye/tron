package utils

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	privKey, pubKey, hexAddr, base58Addr, err := newAccount()
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Printf("%v\n%v\n%v\n%v\n", privKey, pubKey, hexAddr, base58Addr)

	password := "Very1Strange2Pass3Word4"

	storageStr, err := genPrivateKeyStoreage(password, privKey)
	if nil != err {
		t.Error(err)
		return
	}

	storage, result, err := readPrivateKeyStorage(password, storageStr)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Printf("%v\n%#v\n", result, storage)
}
