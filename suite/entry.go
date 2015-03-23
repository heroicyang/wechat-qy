package suite

// ErrorResponse 为请求出错的响应结果
type ErrorResponse struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type suiteTokenResponse struct {
	SuiteAccessToken string  `json:"suite_access_token"`
	ExpiresIn        float64 `json:"expires_in"`
}

type preAuthCodeResponse struct {
	ErrorResponse
	PreAuthCode string  `json:"pre_auth_code"`
	ExpiresIn   float64 `json:"expires_in"`
}

// CorpInfo 为授权方企业信息
type CorpInfo struct {
	ID            string `json:"corpid"`
	Name          string `json:"corp_name"`
	Type          string `json:"corp_type"`
	RoundLogoURI  string `json:"corp_round_logo_url"`
	SquareLogoURI string `json:"corp_square_logo_url"`
	UserMax       string `json:"corp_user_max"`
	AgentMax      string `json:"corp_agent_max"`
	QRCode        string `json:"corp_wxqrcode"`
}

// AgentInfo 为授权方应用信息
type AgentInfo struct {
	ID                   string `json:"agentid"`
	Name                 string `json:"name"`
	RoundLogoURI         string `json:"round_logo_url"`
	SquareLogoURI        string `json:"square_logo_url"`
	Description          string `json:"description"`
	RedirectDomain       string `json:"redirect_domain"`
	RedirectLocationFlag int64  `json:"report_location_flag"`
	IsReportUser         int64  `json:"isreportuser"`
	IsReportEnter        int64  `json:"isreportenter"`
}

type authAgentInfo struct {
	AgentInfo
	AppID    string   `json:"appid"`
	APIGroup []string `json:"api_group"`
}

// Department 为授权的通讯录部门
type department struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentid"`
	Writable string `json:"writable"`
}

// AuthInfo 为授权信息
type AuthInfo struct {
	Agent      []*authAgentInfo `json:"agent"`
	Department []*department    `json:"department"`
}

// AuthUserInfo 为授权的管理员信息
type AuthUserInfo struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

type allowUserInfo struct {
	UserID string `json:"userid"`
	Status string `json:"status"`
}

type allowUserInfos struct {
	User []*allowUserInfo `json:"user"`
}

type allowPartys struct {
	PartyID []int64 `json:"partyid"`
}

type allowTags struct {
	TagID []int64 `json:"tagid"`
}

// PermanentResponse 用于存储获取永久授权码时的响应结果
type PermanentResponse struct {
	AccessToken   string    `json:"access_token"`
	ExpiresIn     float64   `json:"expires_in"`
	PermanentCode string    `json:"permanent_code"`
	AuthCorpInfo  *CorpInfo `json:"auth_corp_info"`
	AuthInfo      *AuthInfo `json:"auth_info"`
}

// AuthInfoResponse 用于存储获取企业号的授权信息时的响应结果
type AuthInfoResponse struct {
	AuthCorpInfo *CorpInfo     `json:"auth_corp_info"`
	AuthInfo     *AuthInfo     `json:"auth_info"`
	AuthUserInfo *AuthUserInfo `json:"auth_user_info"`
}

// AgentResponse 用于存储获取授权方的企业号某个应用的基本信息
type AgentResponse struct {
	ErrorResponse
	AgentInfo
	AllowUserInfos []*allowUserInfos `json:"allow_userinfos"`
	AllowPartys    *allowPartys      `json:"allow_partys"`
	AllowTags      *allowTags        `json:"allow_tags"`
	Close          int64             `json:"close"`
}

// AgentEditInfo 为设置应用时的应用信息
type AgentEditInfo struct {
	AgentInfo
	LogoMediaID string `json:"logo_mediaid"`
}

// CorpAccessTokenResponse 用于存储获取企业号 access token 的响应结果
type CorpAccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}
