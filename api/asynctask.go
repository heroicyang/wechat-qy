package api

import (
	"encoding/json"
	"net/url"
)

const (
	inviteTaskURI            = "https://qyapi.weixin.qq.com/cgi-bin/batch/inviteuser"
	updateUsersTaskURI       = "https://qyapi.weixin.qq.com/cgi-bin/batch/syncuser"
	replaceUsersTaskURI      = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceuser"
	replaceDepartmentTaskURI = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceparty"
)

// AsyncTaskCallback 异步任务的回调信息
type AsyncTaskCallback struct {
	URL            string `json:"url"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encodingaeskey"`
}

// InviteTask 批量邀请成员关注企业号的异步任务信息
type InviteTask struct {
	ToUser     string            `json:"touser,omitempty"`
	ToParty    string            `json:"toparty,omitempty"`
	ToTag      string            `json:"totag,omitempty"`
	InviteTips string            `json:"invite_tips,omitempty"`
	Callback   AsyncTaskCallback `json:"callback"`
}

// UpdateContactTask 增量更新企业成员、全量更新成员或部门的异步任务信息
type UpdateContactTask struct {
	MediaID  string            `json:"media_id"`
	Callback AsyncTaskCallback `json:"callback"`
}

// PerformInviteTask 方法执行邀请成员关注的任务
func (a *API) PerformInviteTask(task InviteTask) (string, error) {
	return a.performTask(inviteTaskURI, task)
}

// PerformUpdateUsersTask 方法执行增量更新成员的任务
func (a *API) PerformUpdateUsersTask(task UpdateContactTask) (string, error) {
	return a.performTask(updateUsersTaskURI, task)
}

// PerformReplaceUsersTask 方法执行全量更新成员的任务
func (a *API) PerformReplaceUsersTask(task UpdateContactTask) (string, error) {
	return a.performTask(replaceUsersTaskURI, task)
}

// PerformReplaceDepartmentTask 方法执行全量更新部门的任务
func (a *API) PerformReplaceDepartmentTask(task UpdateContactTask) (string, error) {
	return a.performTask(replaceDepartmentTaskURI, task)
}

func (a *API) performTask(baseURI string, task interface{}) (string, error) {
	token, err := a.Tokener.Token()
	if err != nil {
		return "", err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)

	uri := baseURI + "?" + qs.Encode()

	data, err := json.Marshal(task)
	if err != nil {
		return "", err
	}

	body, err := a.Client.PostJSON(uri, data)
	if err != nil {
		return "", err
	}

	result := &struct {
		JobID string `json:"jobid"`
	}{}

	err = json.Unmarshal(body, result)

	return result.JobID, err
}
