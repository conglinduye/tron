package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/wlcy/tron/explorer/lib/log"
)

// SendRequest 发送请求
func SendRequest(urlStr, method, queryParam string, postData io.Reader) (buffer *bytes.Buffer, err error) {

	client := &http.Client{}
	//reqQueryParam := ""
	URL := urlStr
	Method := method
	if len(Method) == 0 {
		Method = "GET"
	}
	if len(queryParam) != 0 {
		//reqQueryParam = queryParam
		queryParam = strings.Replace(queryParam, " ", "%20", -1) //转义请求参数中所有的空格
		queryParam = strings.Replace(queryParam, "/", "%2F", -1) //转义请求参数中所有的斜线
		queryParam = strings.Replace(queryParam, "(", "%28", -1) //转义请求参数中所有的左括号
		queryParam = strings.Replace(queryParam, ")", "%29", -1) //转义请求参数中所有的右括号
		queryParam = strings.Replace(queryParam, ",", "%2C", -1) //转义请求参数中所有的逗号
		queryParam = strings.Replace(queryParam, ";", "%3B", -1) //转义请求参数中所有的分号
		l, err := url.Parse("?" + queryParam)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		param := l.Query().Encode()
		URL = fmt.Sprintf("%s?%s", URL, param)
		log.Debugf("req to nubia url:[%v]", URL)
	}
	reqBuffer := &bytes.Buffer{}
	if postData != nil {
		reqBuffer.ReadFrom(postData)
	}
	//reqData := reqBuffer.Bytes()
	req, err := http.NewRequest(Method, URL, reqBuffer)
	if nil != err {
		return nil, err
	}
	/*
		reqMsg := ""

		//log.Debugf("%v request method:[%v], URL:%v, queryParameter:[%v]", b.name, method, urlStr, queryParam)
		if strings.Compare(Method, "GET") == 0 && len(queryParam) != 0 { //记录get请求参数
			reqMsg = fmt.Sprintf("[%v][%s]", time.Now().Format("20060102150405.000000"), reqQueryParam)
		} else {
			reqMsg = fmt.Sprintf("[%v][%s]", time.Now().Format("20060102150405.000000"), reqData)
		}
		log.Debugf(" %v", reqMsg)
	*/
	resp, err := client.Do(req)
	if nil != err {
		log.Errorf("request failed:%v", err)
		return nil, err
	}

	// DEBUG("req:%v", req)

	var data io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		data, err = gzip.NewReader(resp.Body)
		if nil != err {
			return nil, err
		}
	default:
		data = resp.Body
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	buffer = &bytes.Buffer{}
	n, err := buffer.ReadFrom(data)
	if nil != err {
		return nil, err
	}
	log.Debugf("Response status:[%v], body read size:[%v]", resp.StatusCode, n)
	log.Debugf("Response body:[%s]", buffer.Bytes())

	return buffer, nil
}
