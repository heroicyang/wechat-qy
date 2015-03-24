package suite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"

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
	id          string
	secret      string
	accessToken string
	preAuthCode string
	RedirectURI string
	msgCrypt    crypto.WechatMsgCrypt
}

// New 方法用于创建 Suite 实例
func New(suiteID, suiteSecret, suiteToken, suiteEncodingAESKey, redirectURI string) base.RecvHandler {
	msgCrypt, _ := crypto.NewWechatCrypt(suiteToken, suiteEncodingAESKey, suiteID)

	return &Suite{
		id:          suiteID,
		secret:      suiteSecret,
		RedirectURI: redirectURI,
		msgCrypt:    msgCrypt,
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

// Response 方法用于响应应用套件的消息回调（应用套件无需被动响应)
func (s *Suite) Response(message []byte) ([]byte, error) {
	return nil, nil
}

// GetSuiteToken 方法用于获取应用套件令牌
func (s *Suite) GetSuiteToken(suiteTicket string) error {
	buf, _ := json.Marshal(map[string]string{
		"suite_id":     s.id,
		"suite_secret": s.secret,
		"suite_ticket": suiteTicket,
	})

	body, err := utils.SendPostRequest(SuiteTokenURI, buf, headers)
	if err != nil {
		return err
	}

	opResp := &suiteTokenResponse{}
	err = json.Unmarshal(body, opResp)
	if err != nil {
		return err
	}

	s.accessToken = opResp.SuiteAccessToken
	return nil
}

// GetPreAuthCode 方法用于获取应用套件预授权码
func (s *Suite) GetPreAuthCode(appID []string) error {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
	uri := PreAuthCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id": s.id,
		"appid":    appID,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return err
	}

	opResp := &preAuthCodeResponse{}
	err = json.Unmarshal(body, opResp)
	if err != nil {
		return err
	}

	if opResp.ErrCode != "0" {
		return fmt.Errorf("获取预授权码失败：%s", opResp.ErrMsg)
	}

	s.preAuthCode = opResp.PreAuthCode
	return nil
}

// GetAuthURI 返回应用套件的授权地址
func (s *Suite) GetAuthURI(state string) string {
	qs := url.Values{}
	qs.Add("suite_id", s.id)
	qs.Add("pre_auth_code", s.preAuthCode)
	qs.Add("redirect_uri", s.RedirectURI)
	qs.Add("state", state)

	return AuthURI + "?" + qs.Encode()
}

// GetPermanentCode 方法用于获取企业的永久授权码
func (s *Suite) GetPermanentCode(authCode string) (PermanentResponse, error) {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
	uri := PermanentCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id":  s.id,
		"auth_code": authCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return PermanentResponse{}, err
	}

	opResp := PermanentResponse{}
	err = json.Unmarshal(body, &opResp)

	return opResp, err
}

// GetAuthInfo 方法用于获取企业号的授权信息
func (s *Suite) GetAuthInfo(corpID, permanentCode string) (AuthInfoResponse, error) {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
	uri := AuthInfoURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return AuthInfoResponse{}, err
	}

	opResp := AuthInfoResponse{}
	err = json.Unmarshal(body, &opResp)

	return opResp, err
}

// GetAgent 方法用于获取授权方的企业号某个应用的基本信息
func (s *Suite) GetAgent(corpID, permanentCode, agentID string) (AgentResponse, error) {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
	uri := GetAgentURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
		"agentid":        agentID,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return AgentResponse{}, err
	}

	opResp := AgentResponse{}
	err = json.Unmarshal(body, &opResp)

	return opResp, err
}

// SetAgent 方法用于设置企业号应用基本信息
func (s *Suite) SetAgent(corpID, permanentCode string, agent AgentEditInfo) error {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
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

	buf, _ := json.Marshal(data)

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return err
	}

	opResp := ErrorResponse{}
	err = json.Unmarshal(body, &opResp)
	if err != nil {
		return err
	}

	if opResp.ErrCode != "0" {
		return fmt.Errorf("设置企业号应用信息失败：%s", opResp.ErrMsg)
	}

	return nil
}

// GetCorpAccessToken 方法用于获取授权后的企业 access token
func (s *Suite) GetCorpAccessToken(corpID, permanentCode string) (CorpAccessTokenResponse, error) {
	qs := url.Values{}
	qs.Add("suite_access_token", s.accessToken)
	uri := CorpTokenURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := utils.SendPostRequest(uri, buf, headers)
	if err != nil {
		return CorpAccessTokenResponse{}, err
	}

	opResp := CorpAccessTokenResponse{}
	err = json.Unmarshal(body, &opResp)

	return opResp, err
}
