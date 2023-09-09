package mysql

import (
	"web_app/model"
)

//通过通知的类型获取数据库中的字符串
func GetNoticeType(t string) (ty model.NoticeType, err error) {
	err = DB.Model(&ty).Where("type = ?", t).Find(&ty).Error
	return
}

//创建通知
func CreateNotice(notice model.Notice) (err error) {
	err = DB.Create(&notice).Error
	return
}

//查询数据库中是否有该通知
func GetNotice(sendID int, recipientID int, content string, noticeTy string, likeContentID int) (total int64, err error) {
	err = DB.Model(&model.Notice{}).Where("send_id = ? and recipient_id = ? and content = ? and notice_ty= ? and like_content_id = ?", sendID, recipientID, content, noticeTy, likeContentID).Count(&total).Error
	return
}

//查询未读的通知
func GetUnreadNotice(sendID int) (notices []model.Notice, err error) {
	a := DB.Where("recipient_id = ?", sendID)
	err = a.Where("read", false).Order("created_at desc").Find(&notices).Error
	return
}

//获取历史通知
func GetHistoryNotices(sendID int, page int, size int) (notices []model.Notice, err error) {
	limit := size
	offset := (page - 1) * size
	a := DB.Where("recipient_id = ?", sendID)
	err = a.Order("created_at desc").Limit(limit).Offset(offset).Find(&notices).Error
	return
}
