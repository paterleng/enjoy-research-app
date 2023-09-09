package routes

import (
	"github.com/gin-gonic/gin"
	"web_app/controller"
	"web_app/middlewares"
)

func Enter(v2 *gin.RouterGroup) {
	v2.GET("/text", controller.Text)
	//聊天室路由
	chat := v2.Group("/chatRoom")
	{
		chat.GET("/wsHandle", controller.WsHandle)
	}
	v2.Use(middlewares.JWTAuthMiddleware())
	v2.GET("createQRCode", controller.CreateQRCode)
	group := v2.Group("/group")
	group.Use(middlewares.JWTAuthMiddleware())
	{
		//创建群
		group.POST("/createGroup", controller.CreateGroup)
		//向群中添加用户
		group.POST("/addGroupUser", controller.AddGroupUser)
		//删除群聊
		group.DELETE("/deleteGroup", controller.DeleteGroup)
		//移除群聊成员
		group.POST("/deleteGroupUser", controller.DeleteGroupUser)
		//查询所有群用户信息
		group.GET("/selectAllGroupUser", controller.SelectAllGroupUser)
		//查询该用户的聊天列表
		group.GET("/selectAttentionUserList", controller.SelectAttentionUserList)
	}
	singleChat := v2.Group("/singleChat")
	{
		singleChat.PUT("setUnReadToRead", controller.SetUnReadToRead)
	}

}
