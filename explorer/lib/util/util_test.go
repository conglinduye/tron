package util

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	ss := `{"BitTorrent":0,"Bithumb":0,"HuobiToken":0,"IPFS":0,"James":0,"MacCoin":0,"NBACoin":0,"Skypeople":0,"TRXTestCoin":0,"binance":0,"ofoBike":0}`
	tt := ParsingJSONFromString(ss)
	fmt.Print(tt)
}
