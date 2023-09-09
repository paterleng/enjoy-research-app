package mysql

import (
	"gorm.io/gorm"
	"web_app/model"
)

func CreateTreeHole(treehole *model.TreeHole) (err error) {
	err = DB.Create(treehole).Error
	return
}
func GetTreeHole(page int, size int) (treehole []model.TreeHole, err error) {
	limit := size
	offset := (page - 1) * size
	err = DB.Order("created_at desc").Preload("User").Preload("Replays", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Where("parent_id = ?", 0)
	}).Preload("Replays.User").Preload("Replays.CommentReplays").Limit(limit).Offset(offset).Find(&treehole).Error
	return
}
func GetTreeHoleByUserID(page int, size int, userid int) (treehole []model.TreeHole, err error) {
	limit := size
	offset := (page - 1) * size
	err = DB.Order("created_at desc").Preload("User").Preload("Replays", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Where("parent_id = ?", 0)
	}).Preload("Replays.User").Preload("Replays.CommentReplays").Where("user_id = ?", userid).Limit(limit).Offset(offset).Find(&treehole).Error
	return
}

func GetSecondTreeHoleComments(treeHoleCommentId int, page int, size int) (comment []*model.TreeHoleComment, err error) {
	limit := size
	offset := (page - 1) * size
	err = DB.Preload("User").Preload("Parent").Preload("Parent.User").Where("ancestors_id = ? and parent_id != ?", treeHoleCommentId, 0).Limit(limit).Offset(offset).Find(&comment).Error
	return
}
func GetTreeHoleByID(treeHoleID int) (treeHole *model.TreeHole, err error) {
	err = DB.Preload("User").Preload("Replays", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Where("parent_id = ?", 0)
	}).Preload("Replays.User").Preload("Replays.CommentReplays").Where("id = ?", treeHoleID).First(&treeHole).Error
	return
}
func DeleteTreeHoleByID(treeHoleID int) (err error) {
	err = DB.Delete(&model.TreeHole{}, treeHoleID).Error
	return
}
func UpdateTreeHole(treehole *model.TreeHole) (err error) {
	//更新数据
	err = DB.Save(&treehole).Error
	if err != nil {
		return err
	}
	return
}
func GetTreeHoleIsExist(treeHoleID int) (total int64, err error) {
	err = DB.Model(model.TreeHole{}).Where("id = ?", treeHoleID).Count(&total).Error
	return
}
func GetTreeHoleCommentIsExist(parentID int) (total int64, err error) {
	err = DB.Model(model.TreeHoleComment{}).Where("id = ?", parentID).Count(&total).Error
	return
}
func CreateTreeHoleComment(comment *model.TreeHoleComment) (err error) {
	err = DB.Create(comment).Error
	return
}
func GetTreeHoleComments(commentID int) (comment []*model.TreeHoleComment, err error) {
	err = DB.Model(&comment).Preload("CommentReplays").Preload("User").Where("parent_id", commentID).Find(&comment).Error
	return
}
func GetTreeHoleCommentByID(commentID int) (comment *model.TreeHoleComment, err error) {
	err = DB.Model(&comment).Preload("CommentReplays").Preload("User").Preload("Parent").Where("id = ?", commentID).Find(&comment).Error
	return
}
func DeleteTreeHoleComment(commentID int) (err error) {
	err = DB.Delete(&model.TreeHoleComment{}, commentID).Error
	return
}
