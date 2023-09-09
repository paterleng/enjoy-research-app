package model

import "gorm.io/gorm"

type ShuoShuo struct {
	gorm.Model
	UserID            uint             `json:"userID"` //用户id
	Tittle            string           `json:"tittle"` //文章标题
	Content           string           `json:"content"`
	Files             []File           `json:"files"`                         //图片
	LikeNum           int              `json:"like_num"`                      //点赞数
	CollentNum        int              `json:"collentNum"`                    //收藏数
	Address           string           `json:"address"`                       //发布该文章的地址
	CommentNum        int              `json:"comment_num"`                   //评论数
	AddOrCancelLike   bool             `gorm:"-" json:"add_or_cancel_like"`   //是增加还是删除点赞
	IsLikeThisPost    bool             `gorm:"-"`                             //该用户是否给这个帖子点过赞
	IsCollectOrCencel bool             `gorm:"-" json:"is_collect_or_cencel"` //用户是增加收藏还是取消收藏
	IsCollentThisPost bool             `gorm:"-"`                             //该用户是否收藏过这个帖子
	Comments          []Comment        `json:"comments"`                      //评论
	ClassIfications   []ClassIfication `json:"class_ifications"`              //考研资料分类
	User              User
}

type Comment struct {
	gorm.Model
	Replays    []*Comment `gorm:"foreignkey:ParentID;association_foreignkey:id"` //设置自引用
	ParentID   uint       `json:"parent_id"`                                     //评论父级ID  0
	Content    string     `json:"content"`                                       //内容
	RootID     uint       `json:"root_id"`                                       //根id
	ShuoShuoID uint       `json:"shuo_shuo_id"`
	UserID     uint       `json:"userID" ` //用户ID
	User       User
}

type ClassIfication struct {
	gorm.Model
	SchoolId   int  `json:"school_id"`  //学校id
	MajorId    int  `json:"major_id"`   //专业id
	SubjectId  int  `json:"subject_id"` //科目id
	ShuoShuoID uint //关联的说说id
}

type ParamCreateShuos struct {
	Tittle          string           `json:"tittle"` //文章标题
	Content         string           `json:"content"`
	Address         string           `json:"address"` //发布该文章的地址
	FileMsg         []File           `json:"fileMsg"`
	ClassIfications []ClassIfication `json:"class_ifications"` //考研资料分类
}

type ParamUpdatePost struct {
	ID      uint     `json:"id"`      //帖子id
	Tittle  string   `json:"tittle"`  //文章标题
	Content string   `json:"content"` //帖子内容
	Address string   `json:"address"` //发布该文章的地址
	URL     []string `json:"url"`     //文件路径
}

type File struct {
	gorm.Model
	Name       string //文件名字
	Type       string //文件类型
	Size       string //大小
	Url        string //文件路径
	ShuoShuoID uint
}

type Filea struct {
	gorm.Model
	Name string //文件名字
	Type string //文件类型
	Size string //大小
}
