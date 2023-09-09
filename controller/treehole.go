package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/model"
)

func CreateTreeHole(c *gin.Context) {
	userid := GetUserIDByToken(c)
	treeHole := new(model.TreeHole)
	err := c.ShouldBindJSON(&treeHole)
	if err != nil {
		zap.L().Error("发表失败", zap.Error(err))
		ResponseError(c, CodeCreateError)
		return
	}
	treeHole.UserID = userid
	err = mysql.CreateTreeHole(treeHole)
	if err != nil {
		zap.L().Error("插入数据库失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

func GetTreeHole(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//分页查询树洞和一级评论
	treeHoles, err := mysql.GetTreeHole(int(page), int(size))
	if err != nil {
		zap.L().Error("查询树洞错误", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	userid := GetUserIDByToken(c)
	//从redis中拿到登录用户给谁点赞了
	redisRoad := "TreeHoleLike:" + strconv.Itoa(int(userid))
	likeTreeHoleID, err := redis.REDIS.SMembers(context.Background(), redisRoad).Result()
	if err != nil {
		zap.L().Error("查询所有点赞树洞redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var likeTreeHoleIDs []int
	for i := 0; i < len(likeTreeHoleID); i++ {
		id, err := strconv.Atoi(likeTreeHoleID[i])
		if err != nil {
			zap.L().Error("查询所有树洞转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		likeTreeHoleIDs = append(likeTreeHoleIDs, id)
	}
	p := TreeHoleIntersection(likeTreeHoleIDs, treeHoles)
	ResponseSuccess(c, p)
}

//求两个数组的交集
func TreeHoleIntersection(likeTreeHoleIDs []int, treeHoles []model.TreeHole) (p []model.TreeHole) {
	var likeTreeHoleids []int
	for _, treeHole := range treeHoles {
		likeTreeHoleids = append(likeTreeHoleids, int(treeHole.ID))
	}
	var res []int
	set1 := make(map[int]int)
	for _, v := range likeTreeHoleIDs {
		//以数组中的值为键
		set1[v] = 1
	}
	for _, v := range likeTreeHoleids {
		if count, ok := set1[v]; ok && count > 0 {
			res = append(res, v)
			set1[v]--
		}
	}
	for i := 0; i < len(treeHoles); i++ {
		for _, id := range res {
			if int(treeHoles[i].ID) == id {
				treeHoles[i].IsLike = true
				break
			}
		}
	}
	return treeHoles
}
func DeleteTreeHole(c *gin.Context) {
	userID := GetUserIDByToken(c)
	treeHoleIDStr := c.Query("treeHoleID")
	if strings.Trim(treeHoleIDStr, " ") == "" {
		zap.L().Error("参数为空")
		ResponseErrorWithMsg(c, CodeInvalidParam, "您没有权限，不能删除该树洞！")
		return
	}
	treeHoleID, err := strconv.Atoi(treeHoleIDStr)
	treeHole, err := mysql.GetTreeHoleByID(treeHoleID)
	if err != nil {
		zap.L().Error("通过id查询树洞错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "没有该树洞")
		return
	}
	//确定这个树洞是此用户发表的
	if treeHole.UserID == userID {
		err := mysql.DeleteTreeHoleByID(treeHoleID)
		if err != nil {
			zap.L().Error("通过id删除树洞错误", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseSuccess(c, "树洞删除成功")
	} else {
		ResponseErrorWithMsg(c, CodeServerBusy, "此树洞不是您发表的")
	}
}

func GetTreeHoleByID(c *gin.Context) {
	treeHoleIDStr := c.Query("treeHoleID")
	userid := GetUserIDByToken(c)
	if strings.Trim(treeHoleIDStr, " ") == "" {
		zap.L().Error("参数为空")
		ResponseErrorWithMsg(c, CodeInvalidParam, "您没有权限，不能删除该树洞！")
		return
	}
	treeHoleID, err := strconv.Atoi(treeHoleIDStr)
	treeHole, err := mysql.GetTreeHoleByID(treeHoleID)
	if err != nil {
		zap.L().Error("通过id查询树洞错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "没有该树洞")
		return
	}
	//从redis中拿到登录用户给谁点赞了
	redisRoad := "TreeHoleLike:" + strconv.Itoa(int(userid))
	ok, err := redis.REDIS.SIsMember(context.Background(), redisRoad, treeHoleIDStr).Result()
	if err != nil {
		zap.L().Error("查询所有点赞树洞redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	if ok {
		treeHole.IsLike = true
	}
	ResponseSuccess(c, treeHole)
}

//func UpdateTreeHole(c *gin.Context) {
//	treeHole := new(model.TreeHole)
//	err := c.ShouldBindJSON(&treeHole)
//	if err != nil {
//		zap.L().Error("绑定Json失败", zap.Error(err))
//		ResponseError(c, CodeCreateError)
//		return
//	}
//	userid := GetUserIDByToken(c)
//	treeHoleByID, err := mysql.GetTreeHoleByID(int(treeHole.ID))
//	if err != nil {
//		zap.L().Error("数据库查询树洞失败", zap.Error(err))
//		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
//		return
//	}
//	if userid != treeHoleByID.UserID {
//		zap.L().Error("不是该用户的帖子，该用户想去修改", zap.Error(err))
//		ResponseErrorWithMsg(c, CodeServerBusy, "不是你的帖子")
//		return
//	} else {
//		err = mysql.UpdateTreeHole(treeHole)
//		if err != nil {
//			zap.L().Error("数据库更新帖子失败", zap.Error(err))
//			ResponseErrorMsg(c, "更新失败")
//			return
//		}
//		ResponseSuccess(c, CodeSuccess)
//	}
//}
//点赞
func LikeTreeHole(c *gin.Context) {
	userid := GetUserIDByToken(c)
	treeHole := new(model.TreeHole)
	err := c.ShouldBindJSON(&treeHole)
	if err != nil {
		zap.L().Error("绑定失败", zap.Error(err))
		ResponseError(c, CodeCreateError)
		return
	}
	treeHoleByID, err := mysql.GetTreeHoleByID(int(treeHole.ID))
	if err != nil {
		zap.L().Error("通过id查询树洞错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "没有该树洞")
		return
	}
	//点赞
	if treeHole.IsLike == true {
		treeHoleByID.LikeNum++
		//我给谁点过赞
		redisRoad := "TreeHoleLike:" + strconv.Itoa(int(userid))
		err := redis.REDIS.SAdd(context.Background(), redisRoad, treeHoleByID.ID).Err()
		if err != nil {
			zap.L().Error("我点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeCreateError, "我点赞redis操作失败")
			return
		}
		//谁给我点赞了
		redisRoad = "TreeHoleLike:" + strconv.Itoa(int(treeHoleByID.UserID)) + ":" + strconv.Itoa(int(treeHoleByID.ID))
		err = redis.REDIS.SAdd(context.Background(), redisRoad, userid).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeCreateError, "点赞redis操作失败")
			return
		}
		err = SendNotice(c, "likeTreeHole", int(userid), int(treeHoleByID.UserID), int(treeHole.ID), "")
		if err != nil {
			zap.L().Error("发送通知失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}

	} else if treeHole.IsLike == false && treeHoleByID.LikeNum > 0 {
		treeHoleByID.LikeNum--
		//我给谁取消了点赞
		redisRoad := "TreeHoleLike:" + strconv.Itoa(int(userid))
		err = redis.REDIS.SRem(context.Background(), redisRoad, treeHoleByID.ID).Err()
		if err != nil {
			zap.L().Error("我取消点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeCreateError, "我点赞redis操作失败")
			return
		}
		//谁给我取消了点赞
		redisRoad = "TreeHoleLike:" + strconv.Itoa(int(treeHoleByID.UserID)) + ":" + strconv.Itoa(int(treeHoleByID.ID))
		err = redis.REDIS.SRem(context.Background(), redisRoad, userid).Err()
		if err != nil {
			zap.L().Error("点赞取消redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeCreateError, "点赞redis操作失败")
			return
		}
	} else {
		ResponseErrorWithMsg(c, CodeCreateError, "点赞redis或取消操作失败")
		return
	}

	//存数据库
	err = mysql.UpdateTreeHole(treeHoleByID)
	if err != nil {
		zap.L().Error("数据库更新树洞失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "更新失败")
		return
	}
	ResponseSuccess(c, treeHoleByID)
}

//评论树洞
func TreeHoleComments(c *gin.Context) {
	userid := GetUserIDByToken(c)
	var comment *model.TreeHoleComment
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		zap.L().Error("绑定JSON失败", zap.Error(err))
		ResponseError(c, CodeParamError)
		return
	}
	//判断是否是本人发布的评论
	if userid != comment.UserID {
		zap.L().Error("userid与发布评论的用户不一致", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "非本人操作")
		return
	}
	//判断是否有评论的树洞
	treeHoleIsExist, err := mysql.GetTreeHoleIsExist(int(comment.TreeHoleID))
	if err != nil {
		zap.L().Error("数据库查询树洞失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "查询树洞失败")
		return
	}
	if treeHoleIsExist == 0 {
		zap.L().Error("数据库中没有该树洞", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库中没有该树洞")
		return
	}
	//判断是否有该评论的父级评论
	treeHoleCommentIsExist, err := mysql.GetTreeHoleCommentIsExist(int(comment.ParentID))
	if err != nil {
		zap.L().Error("数据库查询父级评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库查询父级评论失败")
		return
	}
	if treeHoleCommentIsExist == 0 && (comment.ParentID) != 0 {
		zap.L().Error("数据库中没有该评论的父级评论", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库中没有该评论的父级评论")
		return
	}
	//添加评论
	err = mysql.CreateTreeHoleComment(comment)
	if err != nil {
		zap.L().Error("创建树洞评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "创建树洞评论失败")
		return
	}
	//评论数增加
	treeHoleByID, err := mysql.GetTreeHoleByID(int(comment.TreeHoleID))
	if err != nil {
		zap.L().Error("通过id查询树洞错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "没有该树洞")
		return
	}
	treeHoleByID.CommentNum++
	//存数据库
	err = mysql.UpdateTreeHole(treeHoleByID)
	if err != nil {
		zap.L().Error("数据库更新树洞失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "更新失败")
		return
	}
	if comment.ParentID == 0 {
		err = SendNotice(c, "commentTreeHole", int(userid), int(treeHoleByID.UserID), int(comment.TreeHoleID), comment.Content)
	} else {
		err = SendNotice(c, "replyCommentsInTreeHole", int(userid), int(treeHoleByID.UserID), int(comment.TreeHoleID), comment.Content)
	}
	if err != nil {
		zap.L().Error("发送通知失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
		return
	}
	ResponseSuccess(c, comment)
}

//查询回复的评论
func GetTreeHoleComments(c *gin.Context) {
	commentID := c.Query("commentsID")
	commentIDNum, err := strconv.Atoi(commentID)
	if err != nil {
		zap.L().Error("转数字失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "转数字失败")
		return
	}
	comments, err := mysql.GetTreeHoleCommentByID(commentIDNum)
	if err != nil {
		zap.L().Error("数据库查询树洞回复失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库查询树洞回复失败")
		return
	}
	ResponseSuccess(c, comments)
}

//查询二级评论
func GetSecondTreeHoleComments(c *gin.Context) {
	page, size := getPageInfo(c)
	treeHoleComment := c.Query("treeHoleComment_id")
	treeHoleCommentId, _ := strconv.Atoi(treeHoleComment)
	if treeHoleCommentId == 0 || page == 0 || size == 0 {
		zap.L().Error("查询二级评论错误")
		ResponseErrorWithMsg(c, CodeServerBusy, "查询二级评论错误")
		return
	}
	secondTreeHoleComments, err := mysql.GetSecondTreeHoleComments(treeHoleCommentId, int(page), int(size))
	if err != nil {
		zap.L().Error("查询二级评论错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询二级评论错误")
		return
	}
	ResponseSuccess(c, secondTreeHoleComments)
}

//删除评论
func DeleteTreeHoleComment(c *gin.Context) {
	userid := GetUserIDByToken(c)
	commentID := c.Query("commentsID")
	commentIDNum, err := strconv.Atoi(commentID)
	if err != nil {
		zap.L().Error("转数字失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "转数字失败")
		return
	}
	comment, err := mysql.GetTreeHoleCommentByID(commentIDNum)
	if err != nil {
		zap.L().Error("数据库查询回复失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库查询回复失败")
		return
	}
	//查询发表该树洞的楼主
	topTreeHole, err := mysql.GetTreeHoleByID(int(comment.TreeHoleID))
	//判断是否是本人发布的评论或者楼主删除
	if userid != comment.UserID && userid != topTreeHole.UserID {
		zap.L().Error("userid与发布评论的用户不一致", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "非本人或楼主操作")
		return
	}
	//判断是否有评论的树洞
	treeHoleIsExist, err := mysql.GetTreeHoleIsExist(int(comment.TreeHoleID))
	if err != nil {
		zap.L().Error("数据库查询树洞失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "查询树洞失败")
		return
	}
	if treeHoleIsExist == 0 {
		zap.L().Error("数据库中没有该树洞", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库中没有该树洞")
		return
	}
	//判断是否有该评论的父级评论
	treeHoleCommentIsExist, err := mysql.GetTreeHoleCommentIsExist(int(comment.ParentID))
	if err != nil {
		zap.L().Error("数据库查询父级评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库查询父级评论失败")
		return
	}
	if treeHoleCommentIsExist == 0 && (comment.ParentID) != 0 {
		zap.L().Error("数据库中没有该评论的父级评论", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "数据库中没有该评论的父级评论")
		return
	}

	err = DeleteCommentAndReplays(comment)
	if err != nil {
		zap.L().Error("递归删除所有评论的子评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "递归删除所有评论的子评论失败")
		return
	}

	//删除评论
	err = mysql.DeleteTreeHoleComment(int(comment.ID))
	if err != nil {
		zap.L().Error("删除评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "删除评论失败")
		return
	}
	//评论数减少
	//查询顶部树洞
	treeHoleByID, err := mysql.GetTreeHoleByID(int(comment.TreeHoleID))
	if err != nil {
		zap.L().Error("通过id查询树洞错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "没有该树洞")
		return
	}
	treeHoleByID.CommentNum--
	//存数据库
	err = mysql.UpdateTreeHole(treeHoleByID)
	if err != nil {
		zap.L().Error("数据库更新树洞失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "更新失败")
		return
	}
	ResponseSuccess(c, comment)
}

//递归删除所有评论的子评论
func DeleteCommentAndReplays(comment *model.TreeHoleComment) (err error) {
	comment, err = mysql.GetTreeHoleCommentByID(int(comment.ID))
	if err != nil {
		zap.L().Error("数据库查询回复失败", zap.Error(err))
		return err
	}
	for _, replay := range comment.CommentReplays {
		err := DeleteCommentAndReplays(replay)
		if err != nil {
			zap.L().Error("递归查询评论失败", zap.Error(err))
			return err
		}
		//删除评论
		err = mysql.DeleteTreeHoleComment(int(replay.ID))
		if err != nil {
			zap.L().Error("删除评论失败", zap.Error(err))
			return err
		}
		//查询顶部树洞
		treeHoleByID, err := mysql.GetTreeHoleByID(int(comment.TreeHoleID))
		if err != nil {
			zap.L().Error("通过id查询树洞错误", zap.Error(err))
			return err
		}
		treeHoleByID.CommentNum--
		//存数据库
		err = mysql.UpdateTreeHole(treeHoleByID)
		if err != nil {
			zap.L().Error("数据库更新树洞失败", zap.Error(err))
			return err
		}
	}
	return nil
}

//获取个人树洞
func GetTreeHoleByUserID(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	userid := GetUserIDByToken(c)
	//分页查询树洞和一级评论
	treeHoles, err := mysql.GetTreeHoleByUserID(int(page), int(size), int(userid))
	if err != nil {
		zap.L().Error("查询树洞错误", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//从redis中拿到登录用户给谁点赞了
	redisRoad := "TreeHoleLike:" + strconv.Itoa(int(userid))
	likeTreeHoleID, err := redis.REDIS.SMembers(context.Background(), redisRoad).Result()
	if err != nil {
		zap.L().Error("查询所有点赞树洞redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var likeTreeHoleIDs []int
	for i := 0; i < len(likeTreeHoleID); i++ {
		id, err := strconv.Atoi(likeTreeHoleID[i])
		if err != nil {
			zap.L().Error("查询所有树洞转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		likeTreeHoleIDs = append(likeTreeHoleIDs, id)
	}
	p := TreeHoleIntersection(likeTreeHoleIDs, treeHoles)
	ResponseSuccess(c, p)
}
