package serve

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/model"
)

//用于在启动时进行监听
func Start(manager *model.ClientManager) {
	for {
		fmt.Println("<-----监听通信管道----->")
		select {
		//监测model.Manager.Register这个的变化，有新的东西加入管道时会被监听到，从而建立连接
		case conn := <-model.Manager.Register:
			fmt.Println("建立新连接:", conn.ID)
			//将新建立的连接加入到用户管理的map中，用于记录连接对象,以连接人的id为键，以连接对象为值
			model.Manager.Clients[conn.ID] = conn
			//返回成功信息
			controller.ResponseWebSocket(conn.Socket, controller.CodeConnectionSuccess, "已连接至服务器")
		//断开连接,监测到变化，有用户断开连接
		case conn := <-model.Manager.Unregister:
			fmt.Println("连接失败:", conn.ID)
			if _, ok := model.Manager.Clients[conn.ID]; ok {
				controller.ResponseWebSocket(conn.Socket, controller.CodeConnectionBreak, "连接已断开")
			}
			//关闭当前用户使用的管道
			//close(conn.Send)
			//删除用户管理中的已连接的用户
			delete(model.Manager.Clients, conn.ID)
		case broadcast := <-model.Manager.Broadcast: //广播消息
			message := broadcast.Message
			recipientID := broadcast.Client.RecipientID
			//给一个变量用于确定状态
			flag := false
			contentid := createId(strconv.Itoa(broadcast.Client.SendID), strconv.Itoa(recipientID))
			rID := strconv.Itoa(recipientID) + "->"
			//遍历客户端连接map,查找该用户有没有在线,判断的是对方的连接例如:1要向2发消息,我现在是用户1,那么我需要判断2->1是否存在在用户管理中
			for id, conn := range model.Manager.Clients {
				//如果找不到就说明用户不在线,与接收人的id比较
				if id != rID {
					continue
				}
				//走到这一步,就说明用户在线,就把消息放入管道里面
				select {
				case conn.Send <- *broadcast:
					flag = true
				default: //否则就把该连接从用户管理中删除
					close(conn.Send)
					delete(model.Manager.Clients, conn.ID)
				}
			}
			//判断完之后就把将消息发给用户
			if flag {
				fmt.Println("用户在线应答")
				//controller.ResponseWebSocket(model.Manager.Clients[rID].Socket, controller.CodeSuccess, string(message))
				//把消息插到数据库中
				msg := model.ChatMessage{
					Direction:   contentid,
					SendID:      broadcast.Client.SendID,
					RecipientID: recipientID,
					GroupID:     broadcast.Client.GroupID,
					Content:     string(message),
					Read:        true,
				}
				err := mysql.DB.Create(&msg).Error
				if err != nil {
					zap.L().Error("在线发送消息出现了错误", zap.Error(err))
				}
			} else { //如果不在线
				//controller.ResponseWebSocket(broadcast.Client.Socket, controller.CodeConnectionSuccess, "对方不在线")
				//把消息插到数据库中
				msg := model.ChatMessage{
					Direction:   contentid,
					SendID:      broadcast.Client.SendID,
					RecipientID: recipientID,
					GroupID:     broadcast.Client.GroupID,
					Content:     string(message),
					Read:        false,
				}
				err := mysql.DB.Create(&msg).Error
				if err != nil {
					zap.L().Error("不在线发送消息出现了错误", zap.Error(err))
				}
			}
		}

	}

}

func createId(uid, toUid string) string {
	return uid + "->" + toUid
}
