package api

import (
	"encoding/json"
	"net/url"
)

const (
	sendMessageURI = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// MessageType 消息类型定义
type MessageType string

// 各种消息类型值
const (
	// 发送和接收的消息类型

	TextMsg  MessageType = "text"
	ImageMsg MessageType = "image"
	VoiceMsg MessageType = "voice"
	VideoMsg MessageType = "video"

	// 发送的额外消息类型

	FileMsg   MessageType = "file"
	NewsMsg   MessageType = "news"
	MpNewsMsg MessageType = "mpnews"

	// 接收的额外消息类型

	LocationMsg MessageType = "location"
	EventMsg    MessageType = "event"
)

// TextContent 为文本类型消息的文本内容
type TextContent struct {
	Content string `json:"content"`
}

// TextMessage 为发送的文本类型消息
type TextMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	Text    TextContent `json:"text"`
	Safe    int         `json:"safe"`
}

// Media 为发送消息的媒体文件内容
type Media struct {
	ID string `json:"media_id"`
}

// ImageMessage 为发送的图片类型消息
type ImageMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	Image   Media       `json:"image"`
	Safe    int         `json:"safe"`
}

// VoiceMessage 为发送的声音类型消息
type VoiceMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	Voice   Media       `json:"voice"`
	Safe    int         `json:"safe"`
}

// VideoContent 为视频类型消息的内容
type VideoContent struct {
	ID          string `json:"media_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// VideoMessage 为发送的视频类型消息
type VideoMessage struct {
	ToUser  string       `json:"touser,omitempty"`
	ToParty string       `json:"toparty,omitempty"`
	ToTag   string       `json:"totag,omitempty"`
	MsgType MessageType  `json:"msgtype"`
	AgentID int64        `json:"agentid"`
	Video   VideoContent `json:"video"`
	Safe    int          `json:"safe"`
}

// FileMessage 为发送的文件类型消息
type FileMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	File    Media       `json:"file"`
	Safe    int         `json:"safe"`
}

// Article 为普通图文消息的文章内容
type Article struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	PicURL      string `json:"picurl,omitempty"`
}

// Articles 为普通图文消息的文章列表
type Articles struct {
	Articles []Article `json:"articles"`
}

// NewsMessage 为发送的普通图文类型消息
type NewsMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	News    Articles    `json:"news"`
}

// MpArticle 为特殊图文消息的文章内容
type MpArticle struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	Content          string `json:"content"`
	Digest           string `json:"digest,omitempty"`
	ShowCoverPic     *int   `json:"show_cover_pic,omitempty"`
}

// MpArticles 为特殊图文消息的文章列表
type MpArticles struct {
	Articles []MpArticle `json:"articles"`
}

// MpNewsMessage 为发送的特殊图文类型消息
type MpNewsMessage struct {
	ToUser  string      `json:"touser,omitempty"`
	ToParty string      `json:"toparty,omitempty"`
	ToTag   string      `json:"totag,omitempty"`
	MsgType MessageType `json:"msgtype"`
	AgentID int64       `json:"agentid"`
	MpNews  MpArticles  `json:"mpnews"`
	Safe    int         `json:"safe"`
}

// SendMessage 方法用于主动发送消息给企业成员
func (a *API) SendMessage(message interface{}) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := sendMessageURI + "?" + qs.Encode()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)
	return err
}
