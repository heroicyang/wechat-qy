package suite

import (
	"time"
)

type tokenInfo struct {
	Token     string `json:"suite_access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Tokener 应用套件令牌
type Tokener struct {
	tokenInfo *tokenInfo
	suite     *Suite
}

// NewTokener 方法用于创建 Tokener 实例
func NewTokener(suite *Suite) *Tokener {
	return &Tokener{suite: suite}
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
	return t.suite.getToken()
}

// SetTokenInfo 方法用于设置应用套件令牌的信息
func (t *Tokener) SetTokenInfo(token string, expiresIn int64) {
	t.tokenInfo = &tokenInfo{token, expiresIn}
}

func (t *Tokener) isValidToken() bool {
	now := time.Now().Unix()

	if now >= t.tokenInfo.ExpiresIn || t.tokenInfo.Token == "" {
		return false
	}

	return true
}
