package suite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	"wechat-qy/base"
	"wechat-qy/utils"

	"github.com/heroicyang/wechat-crypto"
)

// 应用套件相关操作的 API 地址
const (
	SuiteTokenURI    = "https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"
	PreAuthCodeURI   = "https://qyapi.weixin.qq.com/cgi-bin/service/get_pre_auth_code"
	AuthURI          = "https://qy.weixin.qq.com/cgi-bin/loginpage"
	PermanentCodeURI = "https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code"
	AuthInfoURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_auth_info"
	GetAgentURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_agent"
	SetAgentURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/set_agent"
	CorpTokenURI     = "https://qyapi.weixin.qq.com/cgi-bin/service/get_corp_token"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

// Suite 结构体包含了应用套件的相关操作
type Suite struct {
	id        string
	secret    string
	ticket    string
	msgCrypt  crypto.WechatMsgCrypt
	tokenInfo *tokenInfo
}

// New 方法用于创建 Suite 实例
func New(suiteID, suiteSecret, suiteToken, suiteEncodingAESKey string) base.RecvHandler {
	msgCrypt, _ := crypto.NewWechatCrypt(suiteToken, suiteEncodingAESKey, suiteID)

	return &Suite{
		id:       suiteID,
		secret:   suiteSecret,
		msgCrypt: msgCrypt,
	}
}

// Parse 方法用于解析应用套件的消息回调
func (s *Suite) Parse(body []byte, signature, timestamp, nonce string) (interface{}, error) {
	var err error

	reqBody := &base.RecvHTTPReqBody{}
	if err = xml.Unmarshal(body, reqBody); err != nil {
		return nil, err
	}

	if signature != s.msgCrypt.GetSignature(timestamp, nonce, reqBody.Encrypt) {
		return nil, fmt.Errorf("validate signature error")
	}

	origData, suiteID, err := s.msgCrypt.Decrypt(reqBody.Encrypt)
	if err != nil {
		return nil, err
	}

	if suiteID != s.id {
		return nil, fmt.Errorf("the request is from suite[%s], not from suite[%s]", suiteID, s.id)
	}

	probeData := &struct {
		InfoType string
	}{}

	if err = xml.Unmarshal(origData, probeData); err != nil {
		return nil, err
	}

	var data interface{}
	switch probeData.InfoType {
	case "suite_ticket":
		data = &RecvSuiteTicket{}
	case "change_auth", "cancel_auth":
		data = &RecvSuiteAuth{}
	default:
		return nil, fmt.Errorf("unknown message type: %s", probeData.InfoType)
	}

	if err = xml.Unmarshal(origData, data); err != nil {
		return nil, err
	}

	return data, nil
}

// Response 方法用于生成应用套件的被动响应消息
func (s *Suite) Response(message []byte) ([]byte, error) {
	return nil, nil
}

// SetTicket 方法用于设置套件的 ticket 信息
func (s *Suite) SetTicket(suiteTicket string) {
	s.ticket = suiteTicket
}

// Token 方法用于获取当前套件的令牌
func (s *Suite) Token() (token string, err error) {
	if s.isValidToken() {
		token = s.tokenInfo.Token
		return
	}

	return s.RefreshToken()
}

// RefreshToken 方法用于刷新当前套件的令牌
func (s *Suite) RefreshToken() (string, error) {
	tokenInfo, err := s.getToken()
	if err != nil {
		return "", err
	}

	s.tokenInfo = tokenInfo
	return tokenInfo.Token, nil
}

func (s *Suite) isValidToken() bool {
	now := time.Now().Unix()

	if now >= s.tokenInfo.ExpiresIn || s.tokenInfo.Token == "" {
		return false
	}

	return true
}

func (s *Suite) getToken() (*tokenInfo, error) {
	buf, _ := json.Marshal(map[string]string{
		"suite_id":     s.id,
		"suite_secret": s.secret,
		"suite_ticket": s.ticket,
	})

	body, err := utils.SendPostRequest(SuiteTokenURI, buf, headers)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{}
	err = json.Unmarshal(body, tokenInfo)
	if err != nil {
		return nil, err
	}

	tokenInfo.ExpiresIn = time.Now().Add(time.Second * time.Duration(tokenInfo.ExpiresIn)).Unix()

	return tokenInfo, nil
}

func (s *Suite) getPreAuthCode(appIDs []int) (*preAuthCodeInfo, error) {
	token, err := s.Token()
	if err != nil {
		return nil, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := PreAuthCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id": s.id,
		"appid":    appIDs,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return nil, err
	}

	result := &struct {
		base.Error
		preAuthCodeInfo
	}{}

	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	if result.ErrCode != base.ErrCodeOk {
		return nil, &result.Error
	}

	return &result.preAuthCodeInfo, nil
}

// GetAuthURI 方法用于获取应用套件的授权地址
func (s *Suite) GetAuthURI(appIDs []int, redirectURI, state string) (string, error) {
	preAuthCodeInfo, err := s.getPreAuthCode(appIDs)
	if err != nil {
		return "", err
	}

	qs := url.Values{}
	qs.Add("suite_id", s.id)
	qs.Add("pre_auth_code", preAuthCodeInfo.Code)
	qs.Add("redirect_uri", redirectURI)
	qs.Add("state", state)

	return AuthURI + "?" + qs.Encode(), nil
}

// GetPermanentCode 方法用于获取企业的永久授权码
func (s *Suite) GetPermanentCode(authCode string) (PermanentCodeInfo, error) {
	token, err := s.Token()
	if err != nil {
		return PermanentCodeInfo{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := PermanentCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id":  s.id,
		"auth_code": authCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return PermanentCodeInfo{}, err
	}

	resp := PermanentCodeInfo{}
	err = json.Unmarshal(body, &resp)

	return resp, err
}

// GetCorpAuthInfo 方法用于获取已授权当前套件的企业号的授权信息
func (s *Suite) GetCorpAuthInfo(corpID, permanentCode string) (CorpAuthInfo, error) {
	token, err := s.Token()
	if err != nil {
		return CorpAuthInfo{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := AuthInfoURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return CorpAuthInfo{}, err
	}

	corpAuthInfo := CorpAuthInfo{}
	err = json.Unmarshal(body, &corpAuthInfo)

	return corpAuthInfo, err
}

// GetCropAgent 方法用于获取已授权当前套件的企业号的某个应用信息
func (s *Suite) GetCropAgent(corpID, permanentCode, agentID string) (CorpAgent, error) {
	token, err := s.Token()
	if err != nil {
		return CorpAgent{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := GetAgentURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
		"agentid":        agentID,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return CorpAgent{}, err
	}

	result := &struct {
		base.Error
		CorpAgent
	}{}

	if err = json.Unmarshal(body, result); err != nil {
		return CorpAgent{}, err
	}
	if result.ErrCode != base.ErrCodeOk {
		return CorpAgent{}, &result.Error
	}

	return result.CorpAgent, nil
}

// UpdateCorpAgent 方法用于设置已授权当前套件的企业号的某个应用信息
func (s *Suite) UpdateCorpAgent(corpID, permanentCode string, agent AgentEditInfo) error {
	token, err := s.Token()
	if err != nil {
		return err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := SetAgentURI + "?" + qs.Encode()

	data := struct {
		SuiteID       string        `json:"suite_id"`
		AuthCorpID    string        `json:"auth_corpid"`
		PermanentCode string        `json:"permanent_code"`
		Agent         AgentEditInfo `json:"agent"`
	}{
		s.id,
		corpID,
		permanentCode,
		agent,
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return err
	}

	result := &base.Error{}
	if err := json.Unmarshal(body, result); err != nil {
		return err
	}

	if result.ErrCode != base.ErrCodeOk {
		return result
	}

	return nil
}

// GetCorpToken 方法用于获取已授权当前套件的企业号的 access token 信息
func (s *Suite) GetCorpToken(corpID, permanentCode string) (CorpTokenInfo, error) {
	token, err := s.Token()
	if err != nil {
		return CorpTokenInfo{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := CorpTokenURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return CorpTokenInfo{}, err
	}

	corpTokenInfo := CorpTokenInfo{}
	err = json.Unmarshal(body, &corpTokenInfo)

	return corpTokenInfo, err
}
