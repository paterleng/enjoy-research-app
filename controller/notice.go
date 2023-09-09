package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"web_app/dao/mysql"
	"web_app/model"
)

func SendNotice(c *gin.Context, noticeT string, sendID int, recipientID int, likeContentID int, commentContent string) (err error) {
	//点赞帖子操作
	noticeType, err := mysql.GetNoticeType(noticeT)
	if err != nil {
		zap.L().Error("查询数据库失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	notice := model.Notice{
		SendID:        sendID,
		RecipientID:   recipientID,
		NoticeTy:      noticeT,
		LikeContentID: likeContentID,
		Content:       noticeType.Content,
		Read:          false,
	}
	user, err := mysql.ReturnDataMysql(uint(sendID))
	if err != nil {
		zap.L().Error("回显查数据失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	controllerMessage := model.ControllerMessage{}
	if strings.Contains(noticeT, "TreeHole") {
		treeHole, err := mysql.GetTreeHoleByID(likeContentID)
		if err != nil {
			zap.L().Error("查询数据库失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return err
		}
		if strings.Contains(noticeT, "omment") {
			treeHole.Content = commentContent
			likeContentID = int(uuid.New().ID())
		}
		controllerMessage = model.ControllerMessage{
			UserID:       strconv.Itoa(sendID),
			HeadPortrait: user.AnonymousHeadPortrait,
			Username:     user.AnonymousName,
			Time:         time.Now(),
			Content:      treeHole.Content,
		}
	} else if strings.Contains(noticeT, "Post") {
		post, err := mysql.GetPostOnlyById(likeContentID)
		if err != nil {
			zap.L().Error("查询数据库失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return err
		}
		if strings.Contains(noticeT, "omment") {
			post.Tittle = commentContent
			likeContentID = int(uuid.New().ID())
		}
		controllerMessage = model.ControllerMessage{
			UserID:       strconv.Itoa(sendID),
			HeadPortrait: user.HeadPortrait[0].Url,
			Username:     user.Username,
			Time:         time.Now(),
			Content:      post.Tittle,
		}
	} else {
		controllerMessage = model.ControllerMessage{
			UserID:       strconv.Itoa(sendID),
			HeadPortrait: user.HeadPortrait[0].Url,
			Username:     user.Username,
			Time:         time.Now(),
			Content:      "关注了你",
		}
	}

	total, err := mysql.GetNotice(notice.SendID, notice.RecipientID, notice.Content, noticeT, likeContentID)
	if total == 0 {
		err := mysql.CreateNotice(notice)
		if err != nil {
			zap.L().Error("创建通知失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeInvalidParam, "创建通知失败")
			return err
		}
		client := model.Manager.Clients[strconv.Itoa(sendID)+"->"]
		client.RecipientID = notice.RecipientID
		model.Manager.Broadcast <- &model.Broadcast{
			Client:            client,
			Message:           []byte(noticeType.Content),
			Type:              2,
			ControllerMessage: controllerMessage,
		}
	}
	return
}
