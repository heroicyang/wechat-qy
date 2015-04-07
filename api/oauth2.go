package api

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	oauth2AuthorizeURI = "https://open.weixin.qq.com/connect/oauth2/authorize"
	oauth2GetUserURI   = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
)

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	UserID   string `json:"UserId"`
	DeviceID string `json:"DeviceId"`
}

// GetOAuth2AuthorizeURI 方法用于构建 OAuth2 中企业获取 code 的 URL 地址
func (a *API) GetOAuth2AuthorizeURI(redirectURI, state string) string {
	qs := make(url.Values)
	qs.Add("appid", a.CorpID)
	qs.Add("redirect_uri", redirectURI)
	qs.Add("response_type", "code")
	qs.Add("scope", "snsapi_base")
	qs.Add("state", state)

	return oauth2AuthorizeURI + "?" + qs.Encode() + "#wechat_redirect"
}

// GetOAuth2User 方法用于获取 OAuth2 方式验证登录后的用户信息
func (a *API) GetOAuth2User(agentID int64, code string) (OAuth2UserInfo, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return OAuth2UserInfo{}, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("code", code)
	qs.Add("agentid", strconv.FormatInt(agentID, 10))

	url := oauth2GetUserURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return OAuth2UserInfo{}, err
	}

	result := OAuth2UserInfo{}
	err = json.Unmarshal(body, &result)

	return result, err
}
