package base

import "time"

// TokenFetcher 包含向 API 服务器获取令牌信息的操作
type TokenFetcher interface {
	FetchToken() (token string, expiresIn int64, err error)
}

// Tokener 用于管理应用套件或企业号的令牌信息
type Tokener struct {
	token        string
	expiresIn    int64
	tokenFetcher TokenFetcher
}

// NewTokener 方法用于创建 Tokener 实例
func NewTokener(tokenFetcher TokenFetcher) *Tokener {
	return &Tokener{tokenFetcher: tokenFetcher}
}

// Token 方法用于获取应用套件令牌
func (t *Tokener) Token() (token string, err error) {
	if !t.isValidToken() {
		if err = t.RefreshToken(); err != nil {
			return "", err
		}
	}

	return t.token, nil
}

// RefreshToken 方法用于刷新令牌信息
func (t *Tokener) RefreshToken() error {
	token, expiresIn, err := t.tokenFetcher.FetchToken()
	if err != nil {
		return err
	}

	expiresIn = time.Now().Add(time.Second * time.Duration(expiresIn)).Unix()

	t.token = token
	t.expiresIn = expiresIn

	return nil
}

func (t *Tokener) isValidToken() bool {
	now := time.Now().Unix()

	if now >= t.expiresIn || t.token == "" {
		return false
	}

	return true
}
