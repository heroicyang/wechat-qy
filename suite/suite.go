package suite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/heroicyang/wechat-crypter"
	"github.com/heroicyang/wechat-qy/base"
)

// 应用套件相关操作的 API 地址
const (
	suiteTokenURI    = "https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"
	preAuthCodeURI   = "https://qyapi.weixin.qq.com/cgi-bin/service/get_pre_auth_code"
	authURI          = "https://qy.weixin.qq.com/cgi-bin/loginpage"
	permanentCodeURI = "https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code"
	authInfoURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_auth_info"
	getAgentURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_agent"
	setAgentURI      = "https://qyapi.weixin.qq.com/cgi-bin/service/set_agent"
	corpTokenURI     = "https://qyapi.weixin.qq.com/cgi-bin/service/get_corp_token"
)

// Suite 结构体包含了应用套件的相关操作
type Suite struct {
	id             string
	secret         string
	ticket         string
	token          string
	encodingAESKey string
	msgCrypter     crypter.MessageCrypter
	tokener        *base.Tokener
	client         *base.Client
}

// New 方法用于创建 Suite 实例
func New(suiteID, suiteSecret, suiteToken, suiteEncodingAESKey string) *Suite {
	msgCrypter, _ := crypter.NewMessageCrypter(suiteToken, suiteEncodingAESKey, suiteID)

	suite := &Suite{
		id:             suiteID,
		secret:         suiteSecret,
		token:          suiteToken,
		encodingAESKey: suiteEncodingAESKey,
		msgCrypter:     msgCrypter,
	}

	suite.client = base.NewClient(suite)
	suite.tokener = base.NewTokener(suite)

	return suite
}

// Retriable 方法实现了套件在发起请求遇到 token 错误时，先刷新 token 然后再次发起请求的逻辑
func (s *Suite) Retriable(reqURL string, body []byte) (bool, string, error) {
	u, err := url.Parse(reqURL)
	if err != nil {
		return false, "", nil
	}

	q := u.Query()
	if q.Get("suite_access_token") == "" {
		return false, "", nil
	}

	result := &base.Error{}
	if err := json.Unmarshal(body, result); err != nil {
		return false, "", err
	}

	switch result.ErrCode {
	case base.ErrCodeOk:
		return false, "", nil
	case base.ErrCodeSuiteTokenInvalid, base.ErrCodeSuiteTokenTimeout, base.ErrCodeSuiteTokenFailure:
		if err := s.tokener.RefreshToken(); err != nil {
			return false, "", err
		}

		token, err := s.tokener.Token()
		if err != nil {
			return false, "", err
		}

		q.Set("suite_access_token", token)
		u.RawQuery = q.Encode()
		return true, u.String(), nil
	default:
		return false, "", result
	}
}

// Parse 方法用于解析应用套件的消息回调
func (s *Suite) Parse(body []byte, signature, timestamp, nonce string) (interface{}, error) {
	var err error

	reqBody := &base.RecvHTTPReqBody{}
	if err = xml.Unmarshal(body, reqBody); err != nil {
		return nil, err
	}

	if signature != s.msgCrypter.GetSignature(timestamp, nonce, reqBody.Encrypt) {
		return nil, fmt.Errorf("validate signature error")
	}

	origData, suiteID, err := s.msgCrypter.Decrypt(reqBody.Encrypt)
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
	msgEncrypt, err := s.msgCrypter.Encrypt(string(message))
	if err != nil {
		return nil, err
	}

	nonce := base.GenerateNonce()
	timestamp := base.GenerateTimestamp()
	signature := s.msgCrypter.GetSignature(fmt.Sprintf("%d", timestamp), nonce, msgEncrypt)

	resp := &base.RecvHTTPRespBody{
		Encrypt:      base.StringToCDATA(msgEncrypt),
		MsgSignature: base.StringToCDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        base.StringToCDATA(nonce),
	}

	return xml.MarshalIndent(resp, " ", "  ")
}

// SetTicket 方法用于设置套件的 ticket 信息
func (s *Suite) SetTicket(suiteTicket string) {
	s.ticket = suiteTicket
}

// FetchToken 方法用于向 API 服务器获取套件的令牌信息
func (s *Suite) FetchToken() (token string, expiresIn int64, err error) {
	buf, _ := json.Marshal(map[string]string{
		"suite_id":     s.id,
		"suite_secret": s.secret,
		"suite_ticket": s.ticket,
	})

	body, err := s.client.PostJSON(suiteTokenURI, buf)
	if err != nil {
		return
	}

	tokenInfo := &struct {
		Token     string `json:"suite_access_token"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	if err = json.Unmarshal(body, tokenInfo); err != nil {
		return
	}

	token = tokenInfo.Token
	expiresIn = tokenInfo.ExpiresIn

	return
}

func (s *Suite) getPreAuthCode(appIDs []int) (*preAuthCodeInfo, error) {
	token, err := s.tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := preAuthCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id": s.id,
		"appid":    appIDs,
	})

	body, err := s.client.PostJSON(uri, buf)
	if err != nil {
		return nil, err
	}

	result := &preAuthCodeInfo{}
	err = json.Unmarshal(body, result)

	return result, err
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

	return authURI + "?" + qs.Encode(), nil
}

// GetPermanentCode 方法用于获取企业的永久授权码
func (s *Suite) GetPermanentCode(authCode string) (PermanentCodeInfo, error) {
	token, err := s.tokener.Token()
	if err != nil {
		return PermanentCodeInfo{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := permanentCodeURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]interface{}{
		"suite_id":  s.id,
		"auth_code": authCode,
	})

	body, err := s.client.PostJSON(uri, buf)
	if err != nil {
		return PermanentCodeInfo{}, err
	}

	result := PermanentCodeInfo{}
	err = json.Unmarshal(body, &result)

	return result, err
}

// GetCorpAuthInfo 方法用于获取已授权当前套件的企业号的授权信息
func (s *Suite) GetCorpAuthInfo(corpID, permanentCode string) (CorpAuthInfo, error) {
	token, err := s.tokener.Token()
	if err != nil {
		return CorpAuthInfo{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := authInfoURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := s.client.PostJSON(uri, buf)
	if err != nil {
		return CorpAuthInfo{}, err
	}

	result := CorpAuthInfo{}
	err = json.Unmarshal(body, &result)

	return result, err
}

// GetCropAgent 方法用于获取已授权当前套件的企业号的某个应用信息
func (s *Suite) GetCropAgent(corpID, permanentCode, agentID string) (CorpAgent, error) {
	token, err := s.tokener.Token()
	if err != nil {
		return CorpAgent{}, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := getAgentURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
		"agentid":        agentID,
	})

	body, err := s.client.PostJSON(uri, buf)
	if err != nil {
		return CorpAgent{}, err
	}

	result := CorpAgent{}
	err = json.Unmarshal(body, &result)

	return result, err
}

// UpdateCorpAgent 方法用于设置已授权当前套件的企业号的某个应用信息
func (s *Suite) UpdateCorpAgent(corpID, permanentCode string, agent AgentEditInfo) error {
	token, err := s.tokener.Token()
	if err != nil {
		return err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := setAgentURI + "?" + qs.Encode()

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

	_, err = s.client.PostJSON(uri, buf)
	return err
}

func (s *Suite) fetchCorpToken(corpID, permanentCode string) (*corpTokenInfo, error) {
	token, err := s.tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := url.Values{}
	qs.Add("suite_access_token", token)
	uri := corpTokenURI + "?" + qs.Encode()

	buf, _ := json.Marshal(map[string]string{
		"suite_id":       s.id,
		"auth_corpid":    corpID,
		"permanent_code": permanentCode,
	})

	body, err := s.client.PostJSON(uri, buf)
	if err != nil {
		return nil, err
	}

	result := &corpTokenInfo{}
	err = json.Unmarshal(body, result)

	return result, err
}
