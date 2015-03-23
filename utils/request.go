package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
