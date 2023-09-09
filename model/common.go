package model

import "gorm.io/gorm"

//通知的消息内容
type NoticeType struct {
	gorm.Model
	Type    string //通知的类型
	Content string //通知的通用内容
}

//通知
type Notice struct {
	gorm.Model
	SendID        int    //发送者ID
	RecipientID   int    //接受者ID
	NoticeTy      string //通知的类型
	LikeContentID int    //操作类型的ID
	Content       string //内容
	Read          bool   //是否读了这条消息
}
