package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTB1(*testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	r.POST("/api/transaction-builder/contract/transfer", TBTransfer)
	buff := &bytes.Buffer{}

	buff.WriteString(`{
		"contract": {
		  "ownerAddress": "TPwJS5eC5BPGyMGtYTHNhPTB89sUWjDSSu",
		  "toAddress": "TWxKPGEyGWEP87Z4GrBccQiWQCf5iUHx9E",
		  "amount": 100000
		},
		"key": "FFA5EA61073FB13E1559F182F91E25C3E51C03906428C7BC8C865A335AED7617",
		"broadcast": true
	  }`)

	req, _ := http.NewRequest("POST", "/api/transaction-builder/contract/transfer", buff)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(buff.Len()))

	r.ServeHTTP(w, req)
	respCtx, err := ioutil.ReadAll(w.Body)
	fmt.Printf("%v\n%s\n%v\n", w.Code, respCtx, err)
}
