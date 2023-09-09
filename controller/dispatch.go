package controller

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"strconv"
	"web_app/dao/mysql"
	"web_app/model"
)

//聊天的后端调度逻辑
//单聊
func SingleChat(c *model.Client, sendMsg *model.SendMsg) {
	//将消息广播出去
	model.Manager.Broadcast <- &model.Broadcast{
		Client:  c,
		Message: []byte(sendMsg.Content),
		Type:    1,
	}
}

//查看未读消息
func UnreadMessages(c *model.Client) {
	//获取数据库中的未读消息
	msgs, err := mysql.GetMessageUnread(c.SendID)
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	for i, msg := range msgs {
		replyMsg := model.ReplyMsg{
			From:    msg.Direction,
			Content: msg.Content,
		}
		message, _ := json.Marshal(replyMsg)
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		//发送完后将消息设为已读
		msgs[i].Read = true
		err := mysql.UpdateMessage(&msgs[i])
		if err != nil {
			ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
		}
	}
}

//拉取历史消息记录
func HistoryMsg(c *model.Client, sendMsg *model.SendMsg) {
	//拿到传过来的时间
	timeT := TimeStringToGoTime(sendMsg.Content)
	//查找聊天记录
	//做一个分页处理，一次查询十条数据,根据时间去限制次数
	//别人发给当前用户的
	direction := CreateId(strconv.Itoa(c.RecipientID), strconv.Itoa(c.SendID))
	//当前用户发出的
	id := CreateId(strconv.Itoa(c.SendID), strconv.Itoa(c.RecipientID))
	msgs, err := mysql.GetHistoryMsg(direction, id, timeT, 10)
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	//把消息写给用户
	for _, msg := range *msgs {
		replyMsg := model.ReplyMsg{
			From:    msg.Direction,
			Content: msg.Content,
		}
		message, _ := json.Marshal(replyMsg)
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		//发送完后将消息设为已读
		if err != nil {
			ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
		}
	}
}

//群聊消息广播
func GroupChat(c *model.Client, sendMsg *model.SendMsg) {
	//根据消息类型判断是否为群聊消息
	//先去数据库查询该群下的所有用户
	users, err := mysql.GetAllGroupUser(strconv.Itoa(sendMsg.RecipientID))
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	//向群里面的用户广播消息
	for _, user := range users {
		//获取群里每个用户的连接
		if int(user.ID) == c.SendID {
			continue
		}
		c.ID = strconv.Itoa(c.SendID) + "->"
		c.GroupID = strconv.Itoa(sendMsg.RecipientID)
		c.RecipientID = int(user.ID)
		model.Manager.Broadcast <- &model.Broadcast{
			Client:  c,
			Message: []byte(sendMsg.Content),
		}
	}
}
