package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/mysql"
	"web_app/model"
)

//创建群聊
func CreateGroup(c *gin.Context) {
	userid := GetUserIDByToken(c)
	p := new(model.Group)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("创建群聊参数错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "参数错误")
		return
	}
	//生成一个随机的id
	u := uuid.New()
	id := u.ID()
	p.ID = strconv.Itoa(int(id))
	p.GroupNum = 1
	p.GroupOwnerId = int(userid)
	err := mysql.CreateGroup(p)
	if err != nil {
		zap.L().Error("创建群聊失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "创建失败")
		return
	}
	ResponseSuccess(c, "创建成功")
}

//添加群用户
func AddGroupUser(c *gin.Context) {
	//群id和用户id
	userid := GetUserIDByToken(c)
	p := new(model.UsersGroup)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("添加群聊用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "参数错误")
		return
	}
	//先判断该用户是否在这个群里面
	p.UserId = int(userid)
	err, num := mysql.SelectGroupUser(p)
	if num > 0 {
		zap.L().Error("用户已在群里面", zap.Error(err))
		ResponseErrorWithMsg(c, CodeUserExist, "用户已在群里")
		return
	}
	if err != nil && err.Error() != "record not found" {
		zap.L().Error("添加群聊用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "创建失败")
		return
	}
	err = mysql.AddGroupUser(p)
	if err != nil {
		zap.L().Error("添加群聊用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "创建失败")
		return
	}
	ResponseSuccess(c, "添加成功")
}

//删除群聊
func DeleteGroup(c *gin.Context) {
	id := c.Query("id")
	userid := GetUserIDByToken(c)
	//判断该群聊是否属于这个人
	group, err := mysql.GetGroupById(id)
	if err != nil {
		zap.L().Error("查询群聊信息失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
		return
	}
	if int(userid) != group.GroupOwnerId {
		zap.L().Error("不是本用户的群", zap.Error(err))
		ResponseErrorWithMsg(c, CodeNoPower, "删除失败")
		return
	}
	err = mysql.DeleteGroup(group)
	if err != nil {
		zap.L().Error("删除群数据库出错", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//移除群聊成员
func DeleteGroupUser(c *gin.Context) {
	p := new(model.UsersGroup)
	userid := GetUserIDByToken(c)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("添加群聊用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "参数错误")
		return
	}
	//判断该群聊是否属于这个人
	group, err := mysql.GetGroupById(p.GroupId)
	if err != nil {
		zap.L().Error("查询群聊信息失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
		return
	}
	if int(userid) != group.GroupOwnerId {
		zap.L().Error("不是本用户的群", zap.Error(err))
		ResponseErrorWithMsg(c, CodeNoPower, "删除失败")
		return
	}
	//如果群主要移除自己，则调用删除群聊的接口
	if int(userid) == p.UserId {
		m := model.Group{
			ID: p.GroupId,
		}
		err := mysql.DeleteGroup(m)
		if err != nil {
			zap.L().Error("删除群数据库出错", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
			return
		}
		ResponseSuccess(c, CodeSuccess)
	}
	err = mysql.DeleteGroupUser(p)
	if err != nil {
		zap.L().Error("删除群成员数据库出错", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//查询所有群用户
func SelectAllGroupUser(c *gin.Context) {
	group_id := c.Query("group_id")
	group, err := mysql.SelectGroupMember(group_id)
	if err != nil {
		zap.L().Error("查询群聊用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询群聊用户失败")
		return
	}
	ResponseSuccess(c, group)
}

//查询用户关注列表
func SelectAttentionUserList(c *gin.Context) {
	userid := GetUserIDByToken(c)
	s := "Attention:" + strconv.Itoa(int(userid))
	userids, err := GetAttentionIdByRedis(s)
	if err != nil {
		zap.L().Error("查询用户关注列表失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询用户关注列表失败")
		return
	}
	//根据userid获取用户信息
	users, err := mysql.GetUserAllMsg(userids)
	//获取加入的群
	//根据userid获取群id
	groupList, err := mysql.SelectGroupByUserID(userid)
	if err != nil {
		zap.L().Error("查询用户群列表失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询用户群列表失败")
		return
	}
	list := struct {
		Users  []model.User
		Groups []model.Group
	}{
		Users:  users,
		Groups: groupList,
	}
	ResponseErrorWithMsg(c, CodeSuccess, list)
}
