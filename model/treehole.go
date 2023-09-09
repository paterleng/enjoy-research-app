package model

import "gorm.io/gorm"

type TreeHole struct {
	gorm.Model
	User        User
	UserID      uint              `json:"userID"`           //用户ID
	Content     string            `json:"content"`          //树洞内容
	LikeNum     uint              `json:"likeNum"`          //点赞数
	CommentNum  uint              `json:"comment_num"`      //评论数
	Replays     []TreeHoleComment `json:"replays"`          //回复的一级评论
	IsAnonymity bool              `json:"isAnonymity"`      //是否匿名
	IsLike      bool              `gorm:"-" json:"is_like"` //喜欢或者不喜欢
}

type TreeHoleComment struct {
	gorm.Model
	CommentReplays    []*TreeHoleComment `gorm:"foreignkey:ParentID"`
	User              User
	UserID            uint             `json:"userID" ` //用户ID
	Parent            *TreeHoleComment `gorm:"foreignkey:ParentID;references:ID"`
	ParentID          uint             `json:"parent_id"`           //评论父级ID  0
	AncestorsID       uint             `json:"ancestors_id"`        //祖先评论
	Content           string           `json:"content"`             //内容
	IsAnonymityReplay uint             `json:"is_anonymity_replay"` //是否匿名回复
	TreeHoleID        uint             `json:"tree_hole_id"`        //树洞ID
}

type PageParam struct {
	Size   int `json:"size"`   //每页大小
	OffSet int `json:"offset"` //第几页
}
