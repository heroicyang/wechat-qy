package api

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	createDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/create"
	updateDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/update"
	deleteDepartmentURI = "https://qyapi.weixin.qq.com/cgi-bin/department/delete"
	listDepartmentURI   = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
)

// Department 表示部门信息
type Department struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Order    int64  `json:"order,omitempty"`
	ParentID int64  `json:"parentid,omitempty"`
}

// CreateDepartment 方法用于创建部门
func (a *API) CreateDepartment(department *Department) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := createDepartmentURI + "?" + qs.Encode()
	data, err := json.Marshal(department)
	if err != nil {
		return err
	}

	body, err := a.Client.PostJSON(url, data)
	if err != nil {
		return err
	}

	d := &Department{}
	if err := json.Unmarshal(body, d); err != nil {
		return err
	}

	department.ID = d.ID
	return nil
}

// UpdateDepartment 方法用于更新部门信息
func (a *API) UpdateDepartment(department *Department) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	url := updateDepartmentURI + "?" + qs.Encode()
	data, err := json.Marshal(department)
	if err != nil {
		return err
	}

	_, err = a.Client.PostJSON(url, data)
	return err
}

// DeleteDepartment 方法用于删除部门
func (a *API) DeleteDepartment(id int64) error {
	token, err := a.Tokener.Token()
	if err != nil {
		return err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("id", strconv.FormatInt(id, 10))

	url := deleteDepartmentURI + "?" + qs.Encode()

	_, err = a.Client.GetJSON(url)
	return err
}

// ListDepartment 方法用于获取部门列表，获取根部门时 id 为 1
func (a *API) ListDepartment(id int64) ([]*Department, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return nil, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("id", strconv.FormatInt(id, 10))

	url := listDepartmentURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(url)
	if err != nil {
		return nil, err
	}

	result := &struct {
		Departments []*Department `json:"department"`
	}{}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result.Departments, nil
}
