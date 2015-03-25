package base

// Tokener 包含了令牌获取、刷新和设置的操作
type Tokener interface {
	Token() (token string, err error)
	RefreshToken() (token string, err error)
	SetTokenInfo(token string, expiresIn int64)
}
