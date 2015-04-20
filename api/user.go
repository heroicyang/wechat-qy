package api

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	createUserURI      = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	updateUserURI      = "https://qyapi.weixin.qq.com/cgi-bin/user/update"
	deleteUserURI      = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"
	batchDeleteUserURI = "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete"
	getUserURI         = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	listSimpleUserURI  = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	listUserURI        = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	inviteUserURI      = "https://qyapi.weixin.qq.com/cgi-bin/invite/send"
)

// UserAttribute struct 为用户扩展信息
type UserAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UserAttributes struct 为用户扩展信息列表
type UserAttributes struct {
	Attrs []*UserAttribute `json:"attrs,omitempty"`
}

// User struct 为企业用户信息
type User struct {
	UserID        string         `json:"userid"`
	Name          string         `json:"name,omitempty"`
	DepartmentIds []int64        `json:"department,omitempty"`
	Position      string         `json:"position,omitempty"`
	Mobile        string         `json:"mobile,omitempty"`
	Email         string         `json:"email,omitempty"`
	WeixinID      string         `json:"weixinid,omitempty"`
	Enable        *int           `json:"enable,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
	Status        *int           `json:"status,omitempty"`
	ExtAttr       UserAttributes `json:"extattr,omitempty"`
}

// CreateUser 方法用于创建用户
func (a *API) CreateUser(user *User) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := createUserURI + "?" + qs.Encode()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)
	return err
}

// UpdateUser 方法用于更新用户信息
func (a *API) UpdateUser(user *User) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := updateUserURI + "?" + qs.Encode()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)
	return err
}

// DeleteUser 方法用于删除某个用户
func (a *API) DeleteUser(userID string) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("userid", userID)

	url := deleteUserURI + "?" + qs.Encode()

	_, err = a.Client.GetJSON(url)
	return err
}

// BatchDeleteUser 方法用于批量删除用户
func (a *API) BatchDeleteUser(userIds []string) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := batchDeleteUserURI + "?" + qs.Encode()

	data, err := json.Marshal(map[string][]string{
		"useridlist": userIds,
	})
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)
	return err
}

// GetUser 方法用于获取某个用户的信息
func (a *API) GetUser(userID string) (*User, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("userid", userID)

	url := getUserURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(body, user)

	return user, err
}

// ListSimpleUser 方法用于获取部门成员列表（成员仅有简单信息）
func (a *API) ListSimpleUser(departmentID int64, fetchChild *int, status *int) ([]*User, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("department_id", strconv.FormatInt(departmentID, 10))

	if fetchChild != nil {
		qs.Add("fetch_child", strconv.Itoa(*fetchChild))
	}

	if status != nil {
		qs.Add("status", strconv.Itoa(*status))
	}

	url := listSimpleUserURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return nil, err
	}

	result := &struct {
		UserList []*User `json:"userlist"`
	}{}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result.UserList, nil
}

// ListUser 方法用于获取部门成员列表（成员带有详情信息）
func (a *API) ListUser(departmentID int64, fetchChild *int, status *int) ([]*User, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("department_id", strconv.FormatInt(departmentID, 10))

	if fetchChild != nil {
		qs.Add("fetch_child", strconv.Itoa(*fetchChild))
	}

	if status != nil {
		qs.Add("status", strconv.Itoa(*status))
	}

	url := listUserURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return nil, err
	}

	result := &struct {
		UserList []*User `json:"userlist"`
	}{}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result.UserList, nil
}

// InviteUser 方法用于邀请成员关注
func (a *API) InviteUser(userID, inviteTips string) (inviteType int, err error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := inviteUserURI + "?" + qs.Encode()
	data, _ := json.Marshal(map[string]string{
		"userid":      userID,
		"invite_tips": inviteTips,
	})

	body, err := a.Client.PostJSON(url, data)
	if err != nil {
		return
	}

	result := &struct {
		Type int `json:"type"`
	}{}

	if err = json.Unmarshal(body, result); err != nil {
		return
	}

	return result.Type, nil
}
