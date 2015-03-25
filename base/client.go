package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Retrier 是带有重试机制 api 需要实现的接口
type Retrier interface {
	Retry(body []byte) (bool, error)
}

// Client 封装了公共的请求方法
type Client struct {
	*http.Client
	api interface{}
}

// NewClient 方法用于创建 Client 实例
func NewClient(api interface{}) *Client {
	return &Client{http.DefaultClient, api}
}

// GetJSON 方法用于发起 JSON GET 请求
func (c *Client) GetJSON(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.request(req, map[string]string{"Content-Type": "application/json"})
}

// PostJSON 方法用于发起 JSON POST 请求
func (c *Client) PostJSON(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return c.request(req, map[string]string{"Content-Type": "application/json"})
}

func (c *Client) request(req *http.Request, headers map[string]string) ([]byte, error) {
	header := make(http.Header)
	for key, val := range headers {
		header.Set(key, val)
	}
	req.Header = header

	hasRetried := false
	retriable := false
RETRY:
	resp, err := c.Do(req)
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

	switch c.api.(type) {
	case Retrier:
		api := c.api.(Retrier)
		retriable, err = api.Retry(body)
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
