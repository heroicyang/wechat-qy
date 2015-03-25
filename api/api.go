package api

import (
	"encoding/json"
	"time"
	"wechat-qy/base"

	"github.com/heroicyang/wechat-crypto"
)

// 企业号相关接口的 API 接口地址
const (
	CreateMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/create"
	DeleteMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete"
	GetMenuURI    = "https://qyapi.weixin.qq.com/cgi-bin/menu/get"
)

// TokenInfo 企业号 API 的令牌信息
type TokenInfo struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// API 封装了企业号相关的接口操作
type API struct {
	corpID     string
	corpSecret string
	msgCrypt   crypto.WechatMsgCrypt
	client     *base.Client
	tokenInfo  *TokenInfo
}

// New 方法创建 API 实例
func New(corpID, corpSecret, token, encodingAESKey string) *API {
	msgCrypt, _ := crypto.NewWechatCrypt(token, encodingAESKey, corpID)

	api := &API{
		corpID:     corpID,
		corpSecret: corpSecret,
		msgCrypt:   msgCrypt,
	}

	client := base.NewClient(api)
	api.client = client

	return api
}

// Retry 方法实现了 API 在发起请求遇到 token 错误时，先刷新 token 然后再次发起请求的逻辑
func (a *API) Retry(body []byte) (bool, error) {
	result := &base.Error{}
	if err := json.Unmarshal(body, result); err != nil {
		return false, err
	}

	switch result.ErrCode {
	case base.ErrCodeOk:
		return false, nil
	case base.ErrCodeTokenInvalid, base.ErrCodeTokenTimeout:
		if _, err := a.RefreshToken(); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, result
	}
}

// Token 方法用于获取企业号的令牌
func (a *API) Token() (token string, err error) {
	if a.isValidToken() {
		token = a.tokenInfo.Token
		return
	}

	return a.RefreshToken()
}

// RefreshToken 方法用于刷新当前企业号的令牌
func (a *API) RefreshToken() (string, error) {
	tokenInfo, err := a.getToken()
	if err != nil {
		return "", err
	}

	a.tokenInfo = tokenInfo
	return tokenInfo.Token, nil
}

func (a *API) isValidToken() bool {
	now := time.Now().Unix()

	if now >= a.tokenInfo.ExpiresIn || a.tokenInfo.Token == "" {
		return false
	}

	return true
}

func (a *API) getToken() (*TokenInfo, error) {
	return nil, nil
}
