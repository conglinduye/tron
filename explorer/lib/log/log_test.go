package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	var fileName = "applog.log"
	StartLogRotator(fileName, 500*1024, nil, []int{0, 30}, true)
	for i := 0; i < 100000000; i++ {
		var strMessage = "this is a long message."

		Debugln(strMessage)
		Infoln(strMessage)
		Warnln(strMessage)
		Errorln(strMessage)

		if i%1 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
