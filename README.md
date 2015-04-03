# wechat-qy
> 易于使用的微信企业号通用 SDK (Golang)

## 特性

* 支持第三方应用提供商的应用套件相关接口
* 企业号相关的所有 API 同时支持基于应用套件级别的调用和企业号单独调用
* 支持企业号最新的异步任务 API
* Access Token 自动管理和续期
* 支持 Access Token 超期或失效导致接口调用错误时重新获取并自动重试一次当前调用的 API
* 提供被动接收消息（事件）的解析方法，以及生成被动响应消息的方法

## 安装
```bash
$ go get -u github.com/heroicyang/wechat-qy
```

*安利一下通用的微信开放平台加解密库：[github.com/heroicyang/wechat-crypter](https://github.com/heroicyang/wechat-crypter)*

## 使用

### 应用套件级别（适用于第三方应用提供商）
```go
import "github.com/heroicyang/wechat-qy/suite"

wechatSuite := suite.New(suiteID, suiteSecret, suiteToken, suiteEncodingAESKey)

// 解析应用套件的被动回调（ticket, change_auth, cancel_auth）
wechatSuite.Parse(body, signature, timestamp, nonce)

// 设置应用套件的 ticket
wechatSuite.SetTicket(ticket)

// 获取应用套件的引导授权地址
wechatSuite.GetAuthURI(appIDs []int, redirectURI, state)

// 获取企业号的永久授权码
wechatSuite.GetPermanentCode(authCode)

// 创建基于应用套件的 API 调用实例
api := wechatSuite.NewAPI(corpID, permanentCode)

// 基于应用套件调用相关 API 与下面企业号级别的 API 调用一致
```

### 企业号级别（适用于企业自己开发应用）
```go
import "github.com/heroicyang/wechat-qy/api"

wechatAPI := api.New(corpID, corpSecret, token, encodingAESKey)

// 创建被动消息解析器
recvMsgHandler := wechatAPI.NewRecvMsgHandler()
// 解析回调模式被动接收的消息（事件）
recvMsgHandler.Parse(body, signature, timestamp, nonce)
// 生成被动响应消息的加密消息体
recvMsgHandler.Response(message []byte)

// 其它 API 调用
wechatAPI.UploadMedia(...)
wechatAPI.PerformReplaceDepartmentTask(...)
// ...
```

## 文档
[http://godoc.org/github.com/heroicyang/wechat-qy](http://godoc.org/github.com/heroicyang/wechat-qy)
