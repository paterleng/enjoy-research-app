package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"web_app/model"
)

func CreateGroup(group *model.Group) (err error) {
	DB.Transaction(func(tx *gorm.DB) error {
		//先创建群信息
		err = tx.Create(&group).Error
		if err != nil {
			return err
		}
		//向群成员内将群主添加进去
		u := &model.UsersGroup{
			GroupId: group.ID,
			UserId:  group.GroupOwnerId,
		}
		err = tx.Table("users_groups").Create(&u).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

//添加群成员
func AddGroupUser(p *model.UsersGroup) (err error) {
	err = DB.Create(&p).Error
	return
}

//查询用户是否在该群里面
func SelectGroupUser(p *model.UsersGroup) (err error, num int64) {
	err = DB.Where("user_id", p.UserId).First(&p).Count(&num).Error
	return
}

//根据群id查询群成员
func SelectGroupMember(id string) (group model.Group, err error) {
	err = DB.Where("id", id).Preload("Users").Find(&group).Error
	return
}

//以用户为主体获取群所有成员
func GetAllGroupUser(id string) (user []model.User, err error) {
	err = DB.Model(model.Group{ID: id}).Association("Users").Find(&user)
	return
}
func GetGroupById(id string) (group model.Group, err error) {
	err = DB.Where("id", id).First(&group).Error
	return
}

//根据群id删除群
func DeleteGroup(p model.Group) (err error) {
	//同时删除关联表中的内容
	err = DB.Select(clause.Associations).Delete(&p).Error
	return
}

//删除群成员
func DeleteGroupUser(p *model.UsersGroup) (err error) {
	err = DB.Where("user_id = ? and group_id = ?", p.UserId, p.GroupId).Delete(&p).Error
	return
}

//查询用户下的群
func SelectGroupByUserID(userid uint) (groupList []model.Group, err error) {
	DB.Model(model.User{}).Where("id", userid).Association("Groups").Find(&groupList)
	return
}
