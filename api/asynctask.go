package api

import (
	"encoding/json"
	"net/url"

	"github.com/heroicyang/wechat-qy/base"
)

const (
	inviteUsersTaskURI       = "https://qyapi.weixin.qq.com/cgi-bin/batch/inviteuser"
	updateUsersTaskURI       = "https://qyapi.weixin.qq.com/cgi-bin/batch/syncuser"
	replaceUsersTaskURI      = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceuser"
	replaceDepartmentTaskURI = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceparty"
	getTaskResultURI         = "https://qyapi.weixin.qq.com/cgi-bin/batch/getresult"
)

// 异步任务操作类型
const (
	SyncUserTask          = "sync_user"
	ReplaceUserTask       = "replace_user"
	InviteUserTask        = "invite_user"
	ReplaceDepartmentTask = "replace_party"
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

// InviteUserTaskResult 邀请成员关注任务的执行结构信息
type InviteUserTaskResult struct {
	base.Error
	UserID     string `json:"userid"`
	InviteType int    `json:"invitetype"`
}

// UpdateUserTaskResult 增量、全量更新成员任务的执行结果信息
type UpdateUserTaskResult struct {
	base.Error
	Action int    `json:"action"`
	UserID string `json:"userid"`
}

// UpdateDepartmentTaskResult 全量更新部门任务的执行结果信息
type UpdateDepartmentTaskResult struct {
	base.Error
	Action       int   `json:"action"`
	DepartmentID int64 `json:"partyid"`
}

// AsyncTaskResultInfo 为异步任务完成结果信息
type AsyncTaskResultInfo struct {
	Status     int         `json:"status"`
	Type       string      `json:"type"`
	Total      int         `json:"total"`
	Percentage float64     `json:"percentage"`
	RemainTime float64     `json:"remaintime"`
	Result     interface{} `json:"result"`
}

// GetTaskResult 方法用于获取异步任务的完成结果
func (a *API) GetTaskResult(taskID string) (AsyncTaskResultInfo, error) {
	result := AsyncTaskResultInfo{}

	token, err := a.Tokener.Token()
	if err != nil {
		return result, err
	}

	qs := make(url.Values)
	qs.Add("access_token", token)
	qs.Add("jobid", taskID)

	uri := getTaskResultURI + "?" + qs.Encode()

	body, err := a.Client.GetJSON(uri)
	if err != nil {
		return result, err
	}

	probeResult := &struct {
		Type string `json:"type"`
	}{}
	if err = json.Unmarshal(body, probeResult); err != nil {
		return result, err
	}

	switch probeResult.Type {
	case InviteUserTask:
		result.Result = []InviteUserTaskResult{}
	case SyncUserTask, ReplaceUserTask:
		result.Result = []UpdateUserTaskResult{}
	case ReplaceDepartmentTask:
		result.Result = []UpdateDepartmentTaskResult{}
	}

	err = json.Unmarshal(body, &result)

	return result, err
}

// PerformInviteUsersTask 方法执行邀请成员关注的任务
func (a *API) PerformInviteUsersTask(task InviteTask) (string, error) {
	return a.performTask(inviteUsersTaskURI, task)
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
