package api

import (
	"time"

	"github.com/heroicyang/wechat-qy/base"
)

// TokenInfo 企业号 API 的令牌信息
type TokenInfo struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Tokener 企业号应用的令牌
type Tokener struct {
	tokenInfo    *TokenInfo
	tokenFetcher base.TokenFetcher
}

// NewTokener 方法用于创建 Tokener 实例
func NewTokener(tokenFetcher base.TokenFetcher) *Tokener {
	return &Tokener{tokenFetcher: tokenFetcher}
}

// Token 方法用于获取应用套件令牌
func (t *Tokener) Token() (token string, err error) {
	if t.isValidToken() {
		token = t.tokenInfo.Token
		return
	}

	return t.RefreshToken()
}

// RefreshToken 方法用于刷新应用套件令牌信息
func (t *Tokener) RefreshToken() (token string, err error) {
	var expiresIn int64

	token, expiresIn, err = t.tokenFetcher.FetchToken()
	if err != nil {
		return
	}

	expiresIn = time.Now().Add(time.Second * time.Duration(expiresIn)).Unix()

	t.tokenInfo = &TokenInfo{token, expiresIn}

	return
}

func (t *Tokener) isValidToken() bool {
	now := time.Now().Unix()

	if now >= t.tokenInfo.ExpiresIn || t.tokenInfo.Token == "" {
		return false
	}

	return true
}
