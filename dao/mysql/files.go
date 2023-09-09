package mysql

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"web_app/model"
)

func CreateShuoShuo(shuoshuo *model.ShuoShuo) (err error) {
	err = DB.Create(shuoshuo).Error
	return
}

func GetPostById(p model.ShuoShuo) (post model.ShuoShuo, err error) {
	err = DB.Model(&p).Where("ID", p.ID).First(&post).Error
	return
}
func GetPostOnlyById(postId int) (post model.ShuoShuo, err error) {
	err = DB.Model(&post).Where("ID", postId).First(&post).Error
	return
}

//用于更新帖子的点赞及收藏数
func AddOrCencelLikeNum(p *model.ShuoShuo) (err error) {
	err = DB.Save(&p).Error
	return
}

//根据条件查询帖子
func SeletePostByClassify(p model.ClassIfication) (post []model.ClassIfication, err error) {
	//根据条件查出帖子id
	db := DB.Model(&p)
	if p.SchoolId != 0 {
		db.Where("school_id = ? ", p.SchoolId)
	}
	if p.MajorId != 0 {
		db.Where("major_id = ?", p.MajorId)
	}
	if p.SubjectId != 0 {
		db.Where("subject_id = ?", p.SubjectId)
	}
	err = db.Select("shuo_shuo_id").Find(&post).Error
	//根据帖子id查出帖子详情返回
	return
}

func SeletePostById(postids []uint) (post []model.ShuoShuo, err error) {
	err = DB.Where("id IN ?", postids).Preload("Comments.User", "parent_id", 0).Preload("Comments.User").Preload("Comments.Replays").Find(&post).Error
	return
}

func UpdatePost(p model.ShuoShuo) (err error) {
	//根据文件在数据中的id去维护gorm的关系,文件的内容不会被修改,前端给我传文件id有就会维护关系，没有就会失去与说说的关系
	//开启事务
	err = DB.Transaction(func(tx *gorm.DB) error {
		//先更新帖子信息
		err = tx.Updates(&p).Error
		if err != nil {
			return err
		}
		//更新有关帖子的资料信息，这一步是为了如果用户删除一个文件，让删除的文件与说说断开链接
		files := p.Files
		err = tx.Model(&p).Association("Files").Replace(&files)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func DeletePost(p model.ShuoShuo) (err error) {
	err = DB.Delete(&p).Error
	if err != nil {
		zap.L().Error("删除说说数据库失败", zap.Error(err))
		return
	}
	return
}

func MyShuoShuo(userid uint) (err error, ShuoShuos *[]model.ShuoShuo) {
	err = DB.Model(&model.User{Model: gorm.Model{ID: userid}}).Association("ShuoShuos").Find(&ShuoShuos)
	return
}

func AllPost(page int, size int) (posts []model.ShuoShuo, err error) {
	limit := size
	offset := (page - 1) * size
	err = DB.Preload("Files").Preload("Comments", "parent_id", 0).Preload("Comments.User").Preload("Comments.Replays").Order("created_at desc").Limit(limit).Offset(offset).Find(&posts).Error
	return
}

//创建评论
func CreateComment(p *model.Comment) (err error) {
	err = DB.Create(&p).Error
	return
}
func DeleteAllChildComment(comment model.Comment) (err error) {
	//根据评论找到帖子信息
	var post model.ShuoShuo
	err = DB.Where("id", comment.ShuoShuoID).First(&post).Error
	if err != nil {
		return
	}
	var comments []model.Comment
	//查询他下面的所有的回复
	result := DB.Where("parent_id", comment.ID).Find(&comments)
	if result.RowsAffected != 0 {
		for _, comment := range comments {
			DeleteAllChildComment(comment)
			post.CommentNum--
		}
	}
	err = DB.Delete(&comment).Error
	if err != nil {
		return
	}
	//更新评论数
	err = DB.Updates(&post).Error
	if err != nil {
		return
	}
	return
}

func SeleteSecondComment(page, size, rootid int) (err error, p *[]model.Comment) {
	limit := size
	offset := (page - 1) * size
	err = DB.Preload("User").Where("root_id = ? and parent_id != ?", rootid, 0).Order("created_at desc").Limit(limit).Offset(offset).Find(&p).Error
	return
}

//根据评论id查看评论详情
func SelectCommentById(id uint) (comment model.Comment, err error) {
	err = DB.Where("id", id).Find(&comment).Error
	return
}

//根据id查询帖子信息
func SelectPostMessage(postids []int) (posts []model.ShuoShuo, err error) {
	err = DB.Where("id IN ?", postids).Preload("Files").Preload("User").Find(&posts).Error
	return
}
