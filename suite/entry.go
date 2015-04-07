package suite

type preAuthCodeInfo struct {
	Code      string `json:"pre_auth_code"`
	ExpiresIn int64  `json:"expires_in"`
}

type corpTokenInfo struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Corporation 用于表示授权方企业信息
type Corporation struct {
	ID            string `json:"corpid"`
	Name          string `json:"corp_name"`
	Type          string `json:"corp_type"`
	RoundLogoURI  string `json:"corp_round_logo_url"`
	SquareLogoURI string `json:"corp_square_logo_url"`
	UserMax       int    `json:"corp_user_max"`
	AgentMax      int    `json:"corp_agent_max"`
	QRCode        string `json:"corp_wxqrcode"`
}

// Agent 用于表示应用基本信息
type Agent struct {
	ID                   int64  `json:"agentid"`
	Name                 string `json:"name,omitempty"`
	RoundLogoURI         string `json:"round_logo_url,omitempty"`
	SquareLogoURI        string `json:"square_logo_url,omitempty"`
	Description          string `json:"description,omitempty"`
	RedirectDomain       string `json:"redirect_domain,omitempty"`
	RedirectLocationFlag int64  `json:"report_location_flag,omitempty"`
	IsReportUser         int64  `json:"isreportuser,omitempty"`
	IsReportEnter        int64  `json:"isreportenter,omitempty"`
}

type authorizedAgent struct {
	Agent
	AppID    int64    `json:"appid"`
	APIGroup []string `json:"api_group"`
}

type authorizedDepartment struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parentid"`
	Writable bool   `json:"writable"`
}

// AuthInfo 表示授权基本信息
type AuthInfo struct {
	Agent      []*authorizedAgent      `json:"agent"`
	Department []*authorizedDepartment `json:"department"`
}

// PermanentCodeInfo 代表获取企业号永久授权码时的响应信息
type PermanentCodeInfo struct {
	AccessToken   string       `json:"access_token"`
	ExpiresIn     int64        `json:"expires_in"`
	PermanentCode string       `json:"permanent_code"`
	AuthCorpInfo  *Corporation `json:"auth_corp_info"`
	AuthInfo      *AuthInfo    `json:"auth_info"`
}

type operator struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// CorpAuthInfo 代表企业号的授权信息
type CorpAuthInfo struct {
	AuthCorpInfo *Corporation `json:"auth_corp_info"`
	AuthInfo     *AuthInfo    `json:"auth_info"`
	AuthUserInfo *operator    `json:"auth_user_info"`
}

type allowUser struct {
	UserID string `json:"userid"`
	Status string `json:"status"`
}

type allowUsers struct {
	User []*allowUser `json:"user"`
}

type allowPartys struct {
	PartyID []int64 `json:"partyid"`
}

type allowTags struct {
	TagID []int64 `json:"tagid"`
}

// CorpAgent 用于表示授权方企业号某个应用的基本信息
type CorpAgent struct {
	Agent
	AllowUsers  *allowUsers  `json:"allow_userinfos"`
	AllowPartys *allowPartys `json:"allow_partys"`
	AllowTags   *allowTags   `json:"allow_tags"`
	Close       int64        `json:"close"`
}

// AgentEditInfo 代表设置授权方企业号某个应用时的应用信息
type AgentEditInfo struct {
	Agent
	LogoMediaID string `json:"logo_mediaid,omitempty"`
}

// RecvSuiteTicket 用于记录应用套件 ticket 的被动响应结果
type RecvSuiteTicket struct {
	SuiteId     string
	InfoType    string
	TimeStamp   float64
	SuiteTicket string
}

// RecvSuiteAuth 用于记录应用套件授权变更和授权撤销的被动响应结果
type RecvSuiteAuth struct {
	SuiteId    string
	InfoType   string
	TimeStamp  float64
	AuthCorpId string
}
