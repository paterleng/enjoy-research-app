package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/mysql"
	"web_app/dao/redis"
)

//获取我喜欢的帖子
func GetMyLike(c *gin.Context) {
	userid := GetUserIDByToken(c)
	str := "Like:" + strconv.Itoa(int(userid))
	postid, err := redis.REDIS.SMembers(context.Background(), str).Result()
	if err != nil {
		zap.L().Error("查询所有帖子redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var postids []int
	for i := 0; i < len(postid); i++ {
		id, err := strconv.Atoi(postid[i])
		if err != nil {
			zap.L().Error("查询所有帖子转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		postids = append(postids, id)
	}
	//根据id查询帖子信息
	posts, err := mysql.SelectPostMessage(postids)
	if err != nil {
		zap.L().Error("查询所有帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	ResponseSuccess(c, posts)
}

//获取我收藏的帖子
func GetMyContent(c *gin.Context) {
	userid := GetUserIDByToken(c)
	str1 := "Collent:" + strconv.Itoa(int(userid))
	collentpostid, err := redis.REDIS.SMembers(context.Background(), str1).Result()
	if err != nil {
		zap.L().Error("查询所有帖子redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var collentpostids []int
	for i := 0; i < len(collentpostid); i++ {
		ids, err := strconv.Atoi(collentpostid[i])
		if err != nil {
			zap.L().Error("查询所有帖子转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		collentpostids = append(collentpostids, ids)
	}
	//根据id查询帖子信息
	posts, err := mysql.SelectPostMessage(collentpostids)
	if err != nil {
		zap.L().Error("查询所有帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	ResponseSuccess(c, posts)
}
