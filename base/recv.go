package base

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"
)

// CDATAText 用于在 xml 解析时避免转义
type CDATAText struct {
	Text string `xml:",innerxml"`
}

// RecvHTTPReqBody 为回调数据
type RecvHTTPReqBody struct {
	ToUserName string
	AgentID    string
	Encrypt    string
}

// RecvHTTPRespBody 为被动响应给微信的数据
type RecvHTTPRespBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    int
	Nonce        CDATAText
}

// StringToCDATA 方法用于将普通文本变换为 CDATAText 类型
func StringToCDATA(text string) CDATAText {
	return CDATAText{"<![CDATA[" + text + "]]>"}
}

// GenerateNonce 方法生成随机数
func GenerateNonce() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%d", r.Int31())
}

// GenerateTimestamp 方法生成时间戳
func GenerateTimestamp() int {
	return int(time.Now().Unix())
}

// RecvHandler 为微信消息回调模式需要实现的接口
type RecvHandler interface {
	Parse(body []byte, signature, timestamp, nonce string) (interface{}, error)
	Response(message []byte) ([]byte, error)
}
