package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/url"
)

const (
	uploadMediaURI   = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	downloadMediaURI = "https://qyapi.weixin.qq.com/cgi-bin/media/get"
)

type mediaType string

// 媒体文件类型
const (
	ImageMedia mediaType = "image"
	VoiceMedia mediaType = "voice"
	VideoMedia mediaType = "video"
	FileMedia  mediaType = "file"
)

// UploadedMedia 上传成功的媒体文件信息
type UploadedMedia struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

// UploadMedia 方法用于将媒体文件上传至微信服务器
func (a *API) UploadMedia(mediaType mediaType, filename string, reader io.Reader) (UploadedMedia, error) {
	media := UploadedMedia{}

	token, err := a.Tokener.Token()
	if err != nil {
		return media, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("type", string(mediaType))

	url := uploadMediaURI + "?" + qs.Encode()

	body, err := a.Client.PostMultipart(url, "media", filename, reader)
	if err != nil {
		return media, err
	}

	err = json.Unmarshal(body, &media)

	return media, err
}

// DownloadMedia 方法用于从微信服务器获取媒体文件，文件流将写入 writer 中，并返回 filename 或者 error 信息
func (a *API) DownloadMedia(mediaID string, writer io.Writer) (string, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return "", err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("media_id", mediaID)

	url := downloadMediaURI + "?" + qs.Encode()

	resp, err := a.Client.GetMedia(url)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("从微信服务器获取媒体文件失败")
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)
	return params["filename"], err
}
