package User

import (
	"web_app/dao/mysql"
	"web_app/model"
)

func Login(p *model.ParameLogin) (user model.User, err error) {
	//给两个参数赋值
	user = model.User{
		Mobile:   p.Mobile,
		Password: p.Password,
	}
	//传递一个指针,查用户是否存在
	user, err = mysql.Login(user)
	if err != nil {
		return model.User{}, err
	}
	//用户信息正确,则给用户生成一个标识,证明用户登录过
	return user, err
}
