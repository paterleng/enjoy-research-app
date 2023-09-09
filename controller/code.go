package controller

type ResCode int64

const (
	CodeSuccess      ResCode = 200 + iota
	CodeInvalidParam ResCode = 1001 + iota
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	//token错误
	CodeInvalidToken
	CodeNeedLogin
	CodeUploadError
	CodeCreateError
	CodeNotNil
	CodeParamError
	CodeIsAttention
	CodeConnectionSuccess
	CodeConnectionFail
	CodeConnectionBreak
	CodeLimiteTimes
	CodeNoPower
	CodeChat
	CodeNotice
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:           "success",
	CodeInvalidParam:      "请求参数错误",
	CodeUserExist:         "用户名已存在",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "密码错误",
	CodeServerBusy:        "服务繁忙",
	CodeInvalidToken:      "无效token",
	CodeNeedLogin:         "需要登录",
	CodeUploadError:       "文件上传失败",
	CodeCreateError:       "创建失败",
	CodeNotNil:            "选择不能为空",
	CodeParamError:        "参数错误",
	CodeIsAttention:       "已经关注，无法操作",
	CodeConnectionSuccess: "websocket连接成功",
	CodeConnectionFail:    "websocket连接失败",
	CodeConnectionBreak:   "websocket连接断开",
	CodeLimiteTimes:       "限制聊天条数",
	CodeNoPower:           "没有权限",
	CodeChat:              "聊天",
	CodeNotice:            "通知",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
