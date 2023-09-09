package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"web_app/model"
)

/**
"code":10001,     程序状态码
"msg":xx,         提示信息
"data":{},		  数据
*/
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"` //omitempty 当data为空的时候就不会展示该字段
}

// 失败
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// 成功
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// 自定义错误
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseErrorMsg(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Msg:  msg,
		Data: nil,
	})
}

//用于websocket的信息返回
func ResponseWebSocket(socket *websocket.Conn, code ResCode, content string) {
	replyMsg := &model.ReplyMsg{
		Code:    int(code),
		Content: content,
	}
	//序列化为json对象，提高兼容性
	msg, _ := json.Marshal(replyMsg)
	rwLocker.RLock()
	_ = socket.WriteMessage(websocket.TextMessage, msg)
	rwLocker.RUnlock()
}
