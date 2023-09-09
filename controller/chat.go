package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"web_app/dao/mysql"
	"web_app/model"
)

//读写锁
var rwLocker sync.RWMutex

func CreateId(uid, toUid string) string {
	return uid + "->" + toUid
}

func Text(c *gin.Context) {
	user, err := mysql.GetAllGroupUser("2215706439")
	if err != nil {
		fmt.Println(err)
	}
	ResponseSuccess(c, user)
}

//当用户进入app后就将websocket进行连接
//用于聊天的接口
func WsHandle(c *gin.Context) {
	get := c.Request.Header.Get("Cookie")
	fmt.Println(get)
	myid := c.Query("myid")
	userid, err := strconv.Atoi(myid)
	if err != nil {
		zap.L().Error("转换失败", zap.Error(err))
		ResponseError(c, CodeParamError)
	}
	//将http协议升级为ws协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	},
		Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户客户端实例，用于记录该用户的连接信息
	client := new(model.Client)
	client = &model.Client{
		ID:     myid + "->",
		SendID: userid,
		Socket: conn,
		Send:   make(chan model.Broadcast),
	}
	//使用管道将实例注册到用户管理上
	model.Manager.Register <- client
	//开启两个协程用于读写消息
	go Read(client)
	go Write(client)
}

//用于读管道中的数据
func Read(c *model.Client) {
	//结束把通道关闭
	defer func() {
		model.Manager.Unregister <- c
		//关闭连接
		_ = c.Socket.Close()
	}()
	for {
		//先测试一下连接能不能连上
		c.Socket.PongHandler()
		sendMsg := new(model.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		c.RecipientID = sendMsg.RecipientID
		if err != nil {
			zap.L().Error("数据格式不正确", zap.Error(err))
			model.Manager.Unregister <- c
			_ = c.Socket.Close()
			return
		}
		//根据要发送的消息类型去判断怎么处理
		//消息类型的后端调度
		switch sendMsg.Type {
		case 1: //私信
			SingleChat(c, sendMsg)
		case 2: //获取未读消息
			UnreadMessages(c)
		case 3: //拉取历史消息记录
			HistoryMsg(c, sendMsg)
		case 4: //群聊消息广播
			GroupChat(c, sendMsg)
		}
	}
}

//用于将数据从管道中读出来，写给用户
func Write(c *model.Client) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		//读取管道里面的信息
		case message, ok := <-c.Send:
			//连接不到就返回消息
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			code := CodeSuccess
			if message.Type == 1 {
				code = CodeChat
			} else if message.Type == 2 {
				code = CodeNotice
			}
			fmt.Println(c.ID+"接收消息：", string(message.Message))
			replyMsg := model.ReplyMsg{
				Code:              int(code),
				Content:           fmt.Sprintf("%s", string(message.Message)),
				ControllerMessage: message.ControllerMessage,
			}
			msg, _ := json.Marshal(replyMsg)
			//将接收的消息发送到对应的websocket连接里
			rwLocker.Lock()
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			rwLocker.Unlock()
		}
	}
}

//将未读消息设置为已读消息
func SetUnReadToRead(c *gin.Context) {
	userid := GetUserIDByToken(c)
	id := c.Query("send_id")
	from := CreateId(id, strconv.Itoa(int(userid)))
	//查询未读消息
	err, msgs := mysql.GetChatByFrom(from)
	if err != nil {
		zap.L().Error("查询未读消息失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	for _, msg := range msgs {
		msg.Read = true
		err := mysql.UpdateMessage(&msg)
		if err != nil {
			zap.L().Error("修改未读消息失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}
}
