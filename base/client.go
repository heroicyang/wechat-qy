package base

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
)

// Retrier 是带有重试机制 api 需要实现的接口
type Retrier interface {
	Retriable(url string, body []byte) (retriable bool, newURL string, err error)
}

// Client 封装了公共的请求方法
type Client struct {
	httpClient *http.Client
	api        interface{}
}

// NewClient 方法用于创建 Client 实例
func NewClient(api interface{}) *Client {
	return &Client{http.DefaultClient, api}
}

// GetJSON 方法用于发起 JSON GET 请求
func (c *Client) GetJSON(url string) ([]byte, error) {
	reqURL := url
	hasRetried := false
	retriable := false
RETRY:
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.request(req, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	switch c.api.(type) {
	case Retrier:
		api := c.api.(Retrier)
		retriable, reqURL, err = api.Retriable(url, body)
		if err != nil {
			return nil, err
		}
	}

	if !hasRetried && retriable {
		hasRetried = true
		goto RETRY
	}

	return body, nil
}

// PostJSON 方法用于发起 JSON POST 请求
func (c *Client) PostJSON(url string, data []byte) ([]byte, error) {
	reqURL := url
	hasRetried := false
	retriable := false
RETRY:
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	body, err := c.request(req, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}

	switch c.api.(type) {
	case Retrier:
		api := c.api.(Retrier)
		retriable, reqURL, err = api.Retriable(url, body)
		if err != nil {
			return nil, err
		}
	}

	if !hasRetried && retriable {
		hasRetried = true
		goto RETRY
	}

	return body, nil
}

// PostMultipart 方法用于发起 multipart/form-data POST 请求
func (c *Client) PostMultipart(url, fieldname, filename string, dataReader io.Reader) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(bodyBuf)

	disposition, err := multipartWriter.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(disposition, dataReader); err != nil {
		return nil, err
	}

	if err = multipartWriter.Close(); err != nil {
		return nil, err
	}

	bodyBytes := bodyBuf.Bytes()

	reqURL := url
	hasRetried := false
	retriable := false
RETRY:
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.ContentLength = int64(len(bodyBytes))

	body, err := c.request(req, map[string]string{
		"Content-Type": multipartWriter.FormDataContentType(),
	})
	if err != nil {
		return nil, err
	}

	switch c.api.(type) {
	case Retrier:
		api := c.api.(Retrier)
		retriable, reqURL, err = api.Retriable(url, body)
		if err != nil {
			return nil, err
		}
	}

	if !hasRetried && retriable {
		hasRetried = true
		goto RETRY
	}

	return body, nil
}

// GetMedia 方法专用于从微信服务器获取媒体文件
func (c *Client) GetMedia(url string) (*http.Response, error) {
	reqURL := url
	hasRetried := false
	retriable := false
RETRY:
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, status[%d]", resp.StatusCode)
	}

	contentType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if contentType != "text/plain" && contentType != "application/json" {
		return resp, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch c.api.(type) {
	case Retrier:
		api := c.api.(Retrier)
		retriable, reqURL, err = api.Retriable(url, body)
		if err != nil {
			return nil, err
		}
	}

	if !hasRetried && retriable {
		hasRetried = true
		goto RETRY
	}

	return nil, nil
}

func (c *Client) request(req *http.Request, headers map[string]string) ([]byte, error) {
	header := make(http.Header)
	for key, val := range headers {
		header.Set(key, val)
	}
	req.Header = header

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, status[%d]", resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
