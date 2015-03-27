package api

import (
	"encoding/json"
	"net/url"

	"github.com/heroicyang/wechat-crypter"
	"github.com/heroicyang/wechat-qy/base"
)

const (
	fetchTokenURI = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
)

// API 封装了企业号相关的接口操作
type API struct {
	corpSecret string
	CorpID     string
	MsgCrypter crypter.MessageCrypter
	Client     *base.Client
	Tokener    *base.Tokener
}

// New 方法创建 API 实例
func New(corpID, corpSecret, token, encodingAESKey string) *API {
	msgCrypter, _ := crypter.NewMessageCrypter(token, encodingAESKey, corpID)

	api := &API{
		corpSecret: corpSecret,
		CorpID:     corpID,
		MsgCrypter: msgCrypter,
	}

	api.Client = base.NewClient(api)
	api.Tokener = base.NewTokener(api)

	return api
}

// Retriable 方法实现了 API 在发起请求遇到 token 错误时，先刷新 token 然后再次发起请求的逻辑
func (a *API) Retriable(reqURL string, body []byte) (bool, string, error) {
	u, err := url.Parse(reqURL)
	if err != nil {
		return false, "", nil
	}

	q := u.Query()
	if q.Get("access_token") == "" {
		return false, "", nil
	}

	result := &base.Error{}
	if err = json.Unmarshal(body, result); err != nil {
		return false, "", nil
	}

	switch result.ErrCode {
	case base.ErrCodeOk:
		return false, "", nil
	case base.ErrCodeTokenInvalid, base.ErrCodeTokenTimeout:
		if err := a.Tokener.RefreshToken(); err != nil {
			return false, "", err
		}

		token, err := a.Tokener.Token()
		if err != nil {
			return false, "", err
		}

		q.Set("access_token", token)
		u.RawQuery = q.Encode()
		return true, u.String(), nil
	default:
		return false, "", result
	}
}

// FetchToken 方法使用企业管理组的密钥向 API 服务器获取企业号的令牌信息
func (a *API) FetchToken() (token string, expiresIn int64, err error) {
	qs := make(url.Values)
	qs.Add("corpid", a.CorpID)
	qs.Add("corpsecret", a.corpSecret)

	url := fetchTokenURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return
	}

	result := &struct {
		Token     string `json:"access_token"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	if err = json.Unmarshal(body, result); err != nil {
		return
	}

	token = result.Token
	expiresIn = result.ExpiresIn

	return
}
