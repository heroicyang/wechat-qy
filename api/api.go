package api

import (
	"encoding/json"
	"wechat-qy/base"

	"github.com/heroicyang/wechat-crypto"
)

// 企业号相关接口的 API 接口地址
const (
	CreateMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/create"
	DeleteMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete"
	GetMenuURI    = "https://qyapi.weixin.qq.com/cgi-bin/menu/get"
)

// API 封装了企业号相关的接口操作
type API struct {
	corpSecret string
	CorpID     string
	MsgCrypt   crypto.WechatMsgCrypt
	Client     *base.Client
	Tokener    base.Tokener
}

// New 方法创建 API 实例
func New(corpID, corpSecret, token, encodingAESKey string) *API {
	msgCrypt, _ := crypto.NewWechatCrypt(token, encodingAESKey, corpID)

	api := &API{
		CorpID:     corpID,
		corpSecret: corpSecret,
		MsgCrypt:   msgCrypt,
	}

	api.Client = base.NewClient(api)
	api.Tokener = NewTokener(api)

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
		if _, err := a.Tokener.RefreshToken(); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, result
	}
}

func (a *API) GetToken() (token string, expiresIn int64, err error) {
	return "", 0, nil
}
