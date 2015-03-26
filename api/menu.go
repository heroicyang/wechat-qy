package api

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	createMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/create"
	deleteMenuURI = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete"
	getMenuURI    = "https://qyapi.weixin.qq.com/cgi-bin/menu/get"
)

// 自定义菜单按钮类型
const (
	MenuButtonTypeClick           = "click"
	MenuButtonTypeView            = "view"
	MenuButtonTypeScanCodePush    = "scancode_push"
	MenuButtonTypeScanCodeWaitMsg = "scancode_waitmsg"
	MenuButtonTypePicSysPhoto     = "pic_sysphoto"
	MenuButtonTypePicPhotoOrAlbum = "pic_photo_or_album"
	MenuButtonTypePicWeixin       = "pic_weixin"
	MenuButtonTypeLocationSelect  = "location_select"
)

// MenuButton 中最多包含 5 个子菜单（二级菜单）
type MenuButton struct {
	Type       string       `json:"type"`
	Name       string       `json:"name"`
	Key        string       `json:"key,omitempty"`
	URL        string       `json:"url,omitempty"`
	SubButtons []MenuButton `json:"sub_button,omitempty"`
}

// Menu 中最多包含 3 个菜单按钮（一级菜单）
type Menu struct {
	Buttons []MenuButton `json:"button"`
}

// CreateMenu 方法用于创建某个应用的菜单
func (a *API) CreateMenu(agentID int64, menu Menu) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("agentid", strconv.FormatInt(agentID, 10))

	url := createMenuURI + "?" + qs.Encode()
	data, err := json.Marshal(menu)
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)

	return err
}

// DeleteMenu 方法用于删除某个应用的菜单
func (a *API) DeleteMenu(agentID int64) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("agentid", strconv.FormatInt(agentID, 10))

	url := deleteMenuURI + "?" + qs.Encode()

	_, err = a.Client.GetJSON(url)

	return err
}

// GetMenu 方法用于获取某个应用的菜单
func (a *API) GetMenu(agentID int64) (Menu, error) {
	var menu Menu

	token, err := a.Tokener.Token()
	if err != nil {
		return menu, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("agentid", strconv.FormatInt(agentID, 10))

	url := getMenuURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return menu, err
	}

	err = json.Unmarshal(body, &menu)

	return menu, err
}
