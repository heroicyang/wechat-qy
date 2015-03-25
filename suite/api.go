package suite

import (
	"time"

	"wechat-qy/api"
	"wechat-qy/base"

	"github.com/heroicyang/wechat-crypto"
)

// API 基于套件的接口调用
type API struct {
	api.API
	permanentCode string
	suite         *Suite
}

// NewAPI 方法创建基于该套件的 API 实例
func (s *Suite) NewAPI(corpID, permanentCode string) *API {
	msgCrypt, _ := crypto.NewWechatCrypt(s.token, s.encodingAESKey, corpID)

	suiteAPI := &API{
		permanentCode: permanentCode,
		suite:         s,
	}

	suiteAPI.CorpID = corpID
	suiteAPI.MsgCrypt = msgCrypt
	suiteAPI.Client = base.NewClient(suiteAPI)
	suiteAPI.Tokener = api.NewTokener(suiteAPI)

	return suiteAPI
}

// GetToken 方法用于向 API 服务器获取授权该套件的企业号的令牌信息
func (a *API) GetToken() (token string, expiresIn int64, err error) {
	corpTokenInfo, err := a.suite.getCorpToken(a.CorpID, a.permanentCode)
	if err != nil {
		return
	}

	token = corpTokenInfo.Token
	expiresIn = time.Now().Add(time.Second * time.Duration(corpTokenInfo.ExpiresIn)).Unix()

	return
}
