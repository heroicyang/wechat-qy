package api

import (
	"encoding/xml"
	"fmt"

	"github.com/heroicyang/wechat-qy/base"
)

// 接收事件类型
const (
	SubscribeEvent       = "subscribe"
	UnsubscribeEvent     = "unsubscribe"
	LocationEvent        = "LOCATION"
	MenuClickEvent       = "CLICK"
	MenuViewEvent        = "VIEW"
	ScanCodePushEvent    = "scancode_push"
	ScanCodeWaitMsgEvent = "scancode_waitmsg"
	PicSysPhotoEvent     = "pic_sysphoto"
	PicPhotoOrAlbumEvent = "pic_photo_or_album"
	PicWeiXinEvent       = "pic_weixin"
	LocationSelectEvent  = "location_select"
	EnterAgentEvent      = "enter_agent"
	BatchJobResultEvent  = "batch_job_result"
)

// RecvBaseData 描述接收到的各类消息或事件的公共结构
type RecvBaseData struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      MessageType
	AgentID      int64
}

// RecvTextMessage 描述接收到的文本类型消息结构
type RecvTextMessage struct {
	RecvBaseData
	MsgID   uint64 `xml:"MsgId"`
	Content string
}

// RecvImageMessage 描述接收到的图片类型消息结构
type RecvImageMessage struct {
	RecvBaseData
	MsgID   uint64 `xml:"MsgId"`
	PicURL  string `xml:"PicUrl"`
	MediaID string `xml:"MediaId"`
}

// RecvVoiceMessage 描述接收到的语音类型消息结构
type RecvVoiceMessage struct {
	RecvBaseData
	MsgID   uint64 `xml:"MsgId"`
	MediaID string `xml:"MediaId"`
	Format  string
}

// RecvVideoMessage 描述接收到的视频类型消息结构
type RecvVideoMessage struct {
	RecvBaseData
	MsgID        uint64 `xml:"MsgId"`
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}

// RecvLocationMessage 描述接收到的地理位置类型消息结构
type RecvLocationMessage struct {
	RecvBaseData
	MsgID     uint64  `xml:"MsgId"`
	LocationX float64 `xml:"Location_X"`
	LocationY float64 `xml:"Location_Y"`
	Scale     int
	Label     string
}

// RecvSubscribeEvent 描述成员关注/取消关注事件的结构
type RecvSubscribeEvent struct {
	RecvBaseData
	Event string
}

// RecvLocationEvent 描述上报地理位置事件的结构
type RecvLocationEvent struct {
	RecvBaseData
	Event     string
	Latitude  float64
	Longitude float64
	Precision float64
}

// RecvMenuEvent 描述菜单事件的结构
type RecvMenuEvent struct {
	RecvBaseData
	Event    string
	EventKey string
}

// ScanCodeInfo 描述扫码事件的相关内容结构
type ScanCodeInfo struct {
	ScanType   string
	ScanResult string
}

// RecvScanCodeEvent 描述扫码推/扫码推事件且弹出“消息接收中”提示框类型事件的结构
type RecvScanCodeEvent struct {
	RecvBaseData
	Event        string
	EventKey     string
	ScanCodeInfo ScanCodeInfo
}

// SendPicMD5Sum 描述发图事件中单个图片的 MD5 信息
type SendPicMD5Sum struct {
	PicMd5Sum string
}

// SendPicItem 描述发图事件中单个图片信息结构
type SendPicItem struct {
	Item SendPicMD5Sum `xml:"item"`
}

// SendPicsInfo 描述发图事件的图片信息结构
type SendPicsInfo struct {
	Count   int64
	PicList []SendPicItem
}

// RecvPicEvent 描述发图事件的结构
type RecvPicEvent struct {
	RecvBaseData
	Event        string
	EventKey     string
	SendPicsInfo SendPicsInfo
}

// SendLocationInfo 描述弹出地理位置选择器事件中地理位置信息结构
type SendLocationInfo struct {
	LocationX float64 `xml:"Location_X"`
	LocationY float64 `xml:"Location_Y"`
	Scale     int
	Label     string
	PoiName   string `xml:"Poiname"`
}

// RecvLocationSelectEvent 描述弹出地理位置选择器事件的结构
type RecvLocationSelectEvent struct {
	RecvBaseData
	Event            string
	EventKey         string
	SendLocationInfo SendLocationInfo
}

// RecvEnterAgentEvent 描述成员进入应用事件的结构
type RecvEnterAgentEvent struct {
	RecvBaseData
	Event    string
	EventKey string
}

// JobResultInfo 描述异步任务完成事件中任务完成情况信息
type JobResultInfo struct {
	JobID   string `xml:"JobId"`
	JobType string
	ErrCode int
	ErrMsg  string
}

// RecvBatchJobResultEvent 描述异步任务完成事件的结构
type RecvBatchJobResultEvent struct {
	RecvBaseData
	Event    string
	BatchJob JobResultInfo
}

// RespBaseData 描述被动响应消息的公共结构
type RespBaseData struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   base.CDATAText
	FromUserName base.CDATAText
	CreateTime   int
	MsgType      base.CDATAText
}

// RespTextMessage 描述被动响应的文本消息的结构
type RespTextMessage struct {
	RespBaseData
	Content base.CDATAText
}

// RespMedia 描述被动响应的媒体内容结构
type RespMedia struct {
	MediaID base.CDATAText `xml:"MediaId"`
}

// RespImageMessage 描述被动相应的图片消息结构
type RespImageMessage struct {
	RespBaseData
	Image RespMedia
}

// RespVoiceMessage 描述被动相应的语音消息结构
type RespVoiceMessage struct {
	RespBaseData
	Voice RespMedia
}

// RespVideoMedia 描述被动相应的视频媒体内容结构
type RespVideoMedia struct {
	MediaID     base.CDATAText `xml:"MediaId"`
	Title       base.CDATAText
	Description base.CDATAText
}

// RespVideoMessage 描述被动相应的语音消息结构
type RespVideoMessage struct {
	RespBaseData
	Video RespVideoMedia
}

// RespArticle 描述被动响应的图文消息结构
type RespArticle struct {
	Title       base.CDATAText
	Description base.CDATAText
	PicURL      base.CDATAText `xml:"PicUrl"`
	URL         base.CDATAText `xml:"Url"`
}

// RespArticleItem 描述被动响应的图文消息结构的包裹
type RespArticleItem struct {
	Item RespArticle `xml:"item"`
}

// RespNewsMessage 描述被动相应的图文消息结构
type RespNewsMessage struct {
	RespBaseData
	ArticleCount int
	Articles     []RespArticleItem
}

type recvMsgHandler struct {
	api *API
}

func (h *recvMsgHandler) Parse(body []byte, signature, timestamp, nonce string) (interface{}, error) {
	var err error

	reqBody := &base.RecvHTTPReqBody{}
	if err = xml.Unmarshal(body, reqBody); err != nil {
		return nil, err
	}

	if signature != h.api.MsgCrypter.GetSignature(timestamp, nonce, reqBody.Encrypt) {
		return nil, fmt.Errorf("validate signature error")
	}

	origData, corpID, err := h.api.MsgCrypter.Decrypt(reqBody.Encrypt)
	if err != nil {
		return nil, err
	}

	if corpID != h.api.CorpID {
		return nil, fmt.Errorf("the request is from corp[%s], not from corp[%s]", corpID, h.api.CorpID)
	}

	probeData := &struct {
		MsgType MessageType
		Event   string
	}{}

	if err = xml.Unmarshal(origData, probeData); err != nil {
		return nil, err
	}

	var data interface{}
	switch probeData.MsgType {
	case TextMsg:
		data = &RecvTextMessage{}
	case ImageMsg:
		data = &RecvImageMessage{}
	case VoiceMsg:
		data = &RecvVoiceMessage{}
	case VideoMsg:
		data = &RecvVideoMessage{}
	case LocationMsg:
		data = &RecvLocationMessage{}
	case EventMsg:
		switch probeData.Event {
		case SubscribeEvent, UnsubscribeEvent:
			data = &RecvSubscribeEvent{}
		case LocationEvent:
			data = &RecvLocationEvent{}
		case MenuClickEvent, MenuViewEvent:
			data = &RecvMenuEvent{}
		case ScanCodePushEvent, ScanCodeWaitMsgEvent:
			data = &RecvScanCodeEvent{}
		case PicSysPhotoEvent, PicPhotoOrAlbumEvent, PicWeiXinEvent:
			data = &RecvPicEvent{}
		case LocationSelectEvent:
			data = &RecvLocationSelectEvent{}
		case EnterAgentEvent:
			data = &RecvEnterAgentEvent{}
		case BatchJobResultEvent:
			data = &RecvBatchJobResultEvent{}
		default:
			return nil, fmt.Errorf("unknown event type: %s", probeData.Event)
		}
	default:
		return nil, fmt.Errorf("unknown message type: %s", probeData.MsgType)
	}

	if err = xml.Unmarshal(origData, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (h *recvMsgHandler) Response(message []byte) ([]byte, error) {
	msgEncrypt, err := h.api.MsgCrypter.Encrypt(string(message))
	if err != nil {
		return nil, err
	}

	nonce := base.GenerateNonce()
	timestamp := base.GenerateTimestamp()
	signature := h.api.MsgCrypter.GetSignature(fmt.Sprintf("%d", timestamp), nonce, msgEncrypt)

	resp := &base.RecvHTTPRespBody{
		Encrypt:      base.StringToCDATA(msgEncrypt),
		MsgSignature: base.StringToCDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        base.StringToCDATA(nonce),
	}

	return xml.MarshalIndent(resp, " ", "  ")
}

// NewRecvMsgHandler 方法用于创建消息接收处理器的实例
func (a *API) NewRecvMsgHandler() *recvMsgHandler {
	return &recvMsgHandler{a}
}
