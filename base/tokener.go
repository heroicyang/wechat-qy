package base

// TokenFetcher 包含向 API 服务器获取令牌信息的操作
type TokenFetcher interface {
	GetToken() (token string, expiresIn int64, err error)
}

// Tokener 包含了令牌获取、刷新和设置的操作
type Tokener interface {
	Token() (token string, err error)
	RefreshToken() (token string, err error)
}
