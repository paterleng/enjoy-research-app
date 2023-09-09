package mysql

import (
	"time"
	"web_app/model"
)

//获取未读消息
func GetMessageUnread(sendid int) (msgs []model.ChatMessage, err error) {
	tx := DB.Where("recipient_id = ?", sendid)
	err = tx.Where("read", false).Order("created_at desc").Find(&msgs).Error
	return
}

//修改消息内容，只能把未读的修改成已读的
func UpdateMessage(msg *model.ChatMessage) (err error) {
	err = DB.Model(&msg).Where("id in ?", msg.ID).Update("read", true).Error
	return
}

//获取历史信息
func GetHistoryMsg(direction string, id string, timeT time.Time, size int) (msgs *[]model.ChatMessage, err error) {
	err = DB.Where("direction = ? and created_at <= ? ", direction, timeT).Or("direction = ? and created_at <= ?", id, timeT).Limit(size).Order("created_at desc").Find(&msgs).Error
	return
}

func GetChatByFrom(from string) (err error, msgs []model.ChatMessage) {
	err = DB.Where("direction", from).Find(&msgs).Error
	return
}
