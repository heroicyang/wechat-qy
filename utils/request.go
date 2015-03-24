package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"wechat-qy/base"
)

// SendGetRequest 方法用于发起 GET 请求
func SendGetRequest(uri string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	return sendRequest(req, headers)
}

// SendPostRequest 方法用于发起 POST 请求
func SendPostRequest(uri string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return sendRequest(req, headers)
}

func sendRequest(req *http.Request, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	header := make(http.Header)
	for key, val := range headers {
		header.Set(key, val)
	}
	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, status[%d]", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errResp := &base.Error{}
	if err = json.Unmarshal(body, errResp); err != nil {
		return nil, err
	}

	if errResp.ErrCode != base.ErrCodeOk {
		return nil, errResp
	}

	return body, nil
}
