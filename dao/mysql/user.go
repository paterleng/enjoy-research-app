package mysql

import (
	"gorm.io/gorm"
	"web_app/model"
)

// 注册
func Regist(u *model.User) (err error) {
	//存到数据库
	err = DB.Create(&u).Error
	//创建完数据后我向数据库添加头像信息
	return err
}

func SelectUser(mobile string) (err error, count int64) {
	var user model.User
	err = DB.Table("users").Where("mobile = ?", mobile).Take(&user).Count(&count).Error
	return
}

// 登录
func Login(user model.User) (u model.User, err error) {
	//查询数据库,判断用户是否存在
	err = DB.Preload("HeadPortrait", "is_using", 1).Where("mobile = ?", user.Mobile).First(&u).Error
	//用户不存在,返回一个错误:用户不存在
	if err == gorm.ErrRecordNotFound {
		return model.User{}, ErrorUserNoExist
	}
	//数据库查询失败
	if err != nil {
		return model.User{}, err
	}
	//判断密码是否正确
	if user.Password != u.Password {
		return model.User{}, ErrorInvalidPassword
	}
	return u, err
}

// 用户回显数据
func ReturnDataMysql(userid uint) (user *model.User, err error) {
	//根据手机号向用户表查询数据
	err = DB.Where("id", userid).Preload("HeadPortrait", "is_using = ?", 1).First(&user).Error
	return user, err
}

//根据学校查询学校id
func SelectSchoolID(school string) (schools *model.School, err error) {
	err = DB.Model(model.School{}).Where("school", school).First(&schools).Error
	return
}

//根据专业查询专业id
func SelectMajorID(major string) (majors *model.Major, err error) {
	err = DB.Model(model.Major{}).Where("major", major).First(&majors).Error
	return
}

// 更新用户头像信息
func UpdateUser(u *model.User) (err error) {
	//开启一个事务
	err = DB.Transaction(func(tx *gorm.DB) error {
		//修改用户表中的数据
		err = tx.Updates(u).Error
		if err != nil {
			return err
		}
		//先判断是否需要修改头像,长度为0代表头像没有修改
		if len(u.HeadPortrait) != 0 {
			if u.HeadPortrait[0].Url != "" {
				//先把该用户所有的头像都设为没有使用,
				err = tx.Table("head_portraits").Where("user_id=?", u.ID).Update("is_using", 0).Error
				if err != nil {
					return err
				}
				//不为空我去更新头像
				//创建一条新的头像信息
				portrait := model.HeadPortrait{
					UserID:  u.ID,
					Url:     u.HeadPortrait[0].Url,
					IsUsing: 1,
				}
				err = tx.Create(&portrait).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	return err
}

// 删除用户
func DeleteUser(userid uint) (err error) {
	DB.Where("id", userid).Delete(&model.User{})
	return err
}

func HeadPortrait(file *model.HeadPortrait) (err error, id int) {
	err = DB.Table("head_portraits").Create(file).Error
	var ids []int
	DB.Raw("select LAST_INSERT_ID() as id").Pluck("id", &ids)
	id = ids[0]
	return err, id
}

//查询匿名信息
func ReturnAnonymous(userid uint) (anonymousUser *model.UserParam, err error) {
	err = DB.Table("users").Where("id = ?", userid).First(&anonymousUser).Error
	return
}

//修改匿名名称
func UpdateAnonymous(u *model.UserParam) (err error) {
	err = DB.Table("users").Where("id", u.ID).Update("anonymous_name", u.AnonymousName).Error
	return
}

func GetUserMessageById(userid uint) (user model.User, err error) {
	err = DB.Table("users").Where("id = ?", userid).First(&user).Error
	return
}

//更改用户信息
func UpdateUserMessage(user model.User) (err error) {
	err = DB.Save(user).Error
	return
}

//根据一组用户id获取用户的详细信息
func GetUserAllMsg(userids []int) (users []model.User, err error) {
	err = DB.Where("id in ?", userids).Find(&users).Error
	return
}
