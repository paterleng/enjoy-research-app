package model

import (
	"gorm.io/gorm"
)

// 用于建数据库
type User struct {
	gorm.Model
	Account               string         `json:"account"`                 //账号
	Password              string         `json:"password"`                //密码
	Mobile                string         `json:"mobile"`                  //手机号
	Username              string         `json:"username"`                //用户名
	Sex                   string         `json:"sex"`                     //性别
	Address               string         `json:"address"`                 //地址
	School                string         `json:"school"`                  //所在院校
	SchoolID              int            `gorm:"-"`                       //所在学校id
	Identity              string         `json:"identity"`                //身份
	Major                 string         `json:"major"`                   //专业
	MajorID               int            `gorm:"-"`                       //所在专业id
	Introduce             string         `json:"introduce"`               //个人简介
	AnonymousName         string         `json:"anonymous_name"`          //匿名名称
	AnonymousHeadPortrait string         `json:"anonymous_head_portrait"` //匿名名称
	HeadPortrait          []HeadPortrait `json:"HeadPortrait"`
	ShuoShuos             []ShuoShuo
	Comment               []Comment
	TreeHoles             []TreeHole
	TreeHoleComment       []TreeHoleComment
	Groups                []Group `gorm:"many2many:users_groups;"`
	AttentionNum          int     `json:"attention_num"` //该用户的关注数
	BeFocused             int     `json:"be_focused"`    //该用户被多少人关注
	AttentionorCencel     bool    `gorm:"-" json:"attentionor_cencel"`
}

//匿名信息
type UserParam struct {
	ID                    uint   `json:"id"`
	AnonymousName         string `json:"anonymous_name"`          //匿名名称
	AnonymousHeadPortrait string `json:"anonymous_head_portrait"` //匿名名称
}

type ParamAttention struct {
	ID                int  `json:"id"`
	AttentionorCencel bool `json:"attentionor_cencel"`
}

//用户头像表
type HeadPortrait struct {
	gorm.Model
	UserID  uint   `json:"userid"`
	Name    string //文件名字
	Type    string //文件类型
	Size    string //大小
	Url     string //用户头像在服务器上的地址
	IsUsing int    `json:"isUsing"` //是否正在使用,1正常 0使用
}
type HeadPortraitParam struct {
	gorm.Model
	Name string //文件名字
	Type string //文件类型
	Size string //大小
	Url  string //用户头像在服务器上的地址
}

// 用于接收前端传来的参数
type ParameRegist struct {
	gorm.Model
	Password string `json:"password" binding:"required"` //密码
	Mobile   string `json:"mobile" binding:"required"`   //手机号
	Username string `json:"username" binding:"required"` //用户名
	School   string `json:"school" binding:"required"`   //所在院校
}

// 登录参数
type ParameLogin struct {
	Mobile   string `json:"mobile" binding:"required"`   //手机号
	Password string `json:"password" binding:"required"` //密码
}
