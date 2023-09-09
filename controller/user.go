package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/DanPlayer/randomname"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/middlewares/jwt"
	"web_app/model"
	"web_app/serve/User"
)

//根据token获取用户id
func GetUserIDByToken(c *gin.Context) (userid uint) {
	uid, bool := c.Get(CtxUserIDKey)
	if bool == false {
		ResponseErrorWithMsg(c, 10000, "登录过期，请重新登录")
		return
	}
	user := fmt.Sprintf("%v", uid)
	users, _ := strconv.Atoi(user)
	userid = uint(users)
	return
}

// 用户注册
func Register(c *gin.Context) {
	//获取注册表单的参数
	u := new(model.ParameRegist)
	//看字段数据是否都有
	if err := c.ShouldBindJSON(&u); err != nil {
		//不符合要求，则不让提交
		zap.L().Error("注册参数：", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查看用户上是否已经注册,根据手机号
	err, count := mysql.SelectUser(u.Mobile)
	if count > 0 {
		ResponseErrorWithMsg(c, 10000, "用户已存在")
		return
	}
	//在用户注册的时候生成用户的默认头像,并保存道数据库
	head := make([]model.HeadPortrait, 1)
	head[0].Url = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/defaultimg.png"
	head[0].IsUsing = 1
	//生成匿名名字
	anonymousName := randomname.GenerateName()
	anonymousHeadPortraits := make([]string, 5)
	anonymousHeadPortraits[0] = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/1689836278IMG_20230720_101916.jpg"
	anonymousHeadPortraits[1] = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/1689836278IMG_20230720_101938.jpg"
	anonymousHeadPortraits[2] = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/1689836278unname1.jpeg"
	anonymousHeadPortraits[3] = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/1689836279unname2.jpg"
	anonymousHeadPortraits[4] = "https://enjoyresearch.oss-cn-hangzhou.aliyuncs.com/HeadPortrait/1689836279unname3.jpg"
	index := rand.Intn(5)
	user := &model.User{
		Mobile:                u.Mobile,
		Password:              u.Password,
		Username:              u.Username,
		School:                u.School,
		HeadPortrait:          head,
		AnonymousName:         anonymousName,
		AnonymousHeadPortrait: anonymousHeadPortraits[index],
		Introduce:             "快填写专属座右铭叭！",
	}
	//参数正确,存到数据库里面
	err = mysql.Regist(user)
	if err != nil {
		zap.L().Error("服务繁忙", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "注册成功")
}

// 用户登录
func Login(c *gin.Context) {
	p := new(model.ParameLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//业务逻辑处理
	user, err := User.Login(p)
	if err != nil {
		zap.L().Error("login.UserLogin failed", zap.String("username", p.Mobile), zap.Error(err))
		//判断出错类型-用户不存在,查不到
		if errors.Is(err, mysql.ErrorUserNoExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		} else {
			//其他问题返回服务繁忙
			ResponseError(c, CodeServerBusy)
			return
		}
	}
	//生成JWT
	token, err := jwt.GenToken(int64(int(user.ID)), user.Username)
	if err != nil {
		return
	}
	//查询没有问题则返回用户部分数据给前端
	ResponseSuccess(c, gin.H{
		"userdata": user,
		"token":    token,
	})
}

//返回用户验证码
func GetMessageCode(c *gin.Context) {
	messageCode := 123456
	ResponseSuccess(c, messageCode)
}

// 回显用户数据
func ReturnData(c *gin.Context) {
	userid := GetUserIDByToken(c)
	//拿到参数查数据库
	user, err := mysql.ReturnDataMysql(userid)
	if err != nil {
		zap.L().Error("回显查数据失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	if user.School != "" {
		//查询用户学校id
		school, err := mysql.SelectSchoolID(user.School)
		if err != nil {
			zap.L().Error("回显查数据失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		user.SchoolID = int(school.ID)

	}

	if user.Major != "" {
		//查询用户专业id
		major, err := mysql.SelectMajorID(user.Major)
		if err != nil {
			zap.L().Error("回显查数据失败", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		user.MajorID = int(major.ID)
	}
	ResponseSuccess(c, user)
}

// 用户修改信息
func UpdateUser(c *gin.Context) {
	userid := GetUserIDByToken(c)
	u := new(model.User)
	if err := c.ShouldBindJSON(&u); err != nil {
		//不符合要求，则不让提交
		zap.L().Error("参数绑定失败：", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//前端传过来用户的头像信息
	u.ID = userid
	//没有问题修改数据库
	err := mysql.UpdateUser(u)
	if err != nil {
		zap.L().Error("用户更新数据失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "更新成功")
}

// 用户注销
func DeleteUser(c *gin.Context) {
	userid := GetUserIDByToken(c)
	//修改数据库信息
	err := mysql.DeleteUser(userid)
	if err != nil {
		zap.L().Error("服务繁忙", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "删除成功")
}

/*
用户头像上传
*/
func HeadPortrait(c *gin.Context) {
	file, err := c.FormFile("headPortrait")
	userid := GetUserIDByToken(c)
	if err != nil {
		zap.L().Error("头像上传失败", zap.Error(err))
		ResponseError(c, CodeUploadError)
		return
	}
	f, err := file.Open()
	if err != nil {
		zap.L().Error("接收文件失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "上传失败")
		return
	}
	//调用阿里云服务上传代码
	url, err := UploadByOSS(f, file.Filename)
	if err != nil {
		zap.L().Error("上传失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "上传失败")
		return
	}
	//用户id会在用户修改信息后自动放入表中
	portrait := model.HeadPortrait{
		UserID: userid,
		Name:   file.Filename,
		Type:   file.Header.Get("Content-Type"),
		Size:   strconv.Itoa(int(file.Size)),
		Url:    url,
	}
	//把数据存到数据库
	err, id := mysql.HeadPortrait(&portrait)
	if err != nil {
		zap.L().Error("数据库插入头像数据失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	portrait.ID = uint(id)
	//封装了一个file类型的对象返回给前端
	ResponseSuccess(c, portrait)
}

//获取请求用户的匿名信息
func GetAnonymousName(c *gin.Context) {
	userID := GetUserIDByToken(c)
	if userID == 0 {
		zap.L().Error("参数为空")
		ResponseErrorWithMsg(c, CodeInvalidParam, "您没有权限，不能查询此人匿名信息！")
		return
	}
	//拿到参数查数据库
	user, err := mysql.ReturnAnonymous(userID)
	if err != nil {
		zap.L().Error("查匿名信息失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查匿名信息失败")
		return
	}
	ResponseSuccess(c, user)
}
func UpdateAnonymousName(c *gin.Context) {
	userid := GetUserIDByToken(c)
	u := new(model.UserParam)
	if err := c.ShouldBindJSON(&u); err != nil {
		//不符合要求，则不让提交
		zap.L().Error("参数绑定失败：", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	u.ID = userid

	//没有问题修改数据库
	err := mysql.UpdateAnonymous(u)
	if err != nil {
		zap.L().Error("更新匿名名称失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "更新成功")
}

//用户关注处理
func AttentionUser(c *gin.Context) {
	var p *model.User
	userid := GetUserIDByToken(c)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("用户关注处理参数错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数错误")
		return
	}
	//查询当前用户有没有关注过该用户
	s := "Attention:" + strconv.Itoa(int(userid))
	userids, err2 := GetAttentionIdByRedis(s)
	if err2 != nil {
		zap.L().Error("查询所有关注redis失败", zap.Error(err2))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	err2, isattention := IsAttention(userids, int(p.ID))
	if err2 != nil {
		zap.L().Error("查询所有关注redis失败", zap.Error(err2))
		ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
		return
	}
	if isattention == true && p.AttentionorCencel == true {
		//表示该用户已经关注过该用户
		ResponseErrorWithMsg(c, CodeIsAttention, "你已经关注过该用户，不能继续操作")
		return
	}
	//根据用户id查询当前登录用户的信息
	//被关注的用户信息
	befocused, err1 := mysql.GetUserMessageById(p.ID)
	if err1 != nil {
		zap.L().Error("查询用户错误", zap.Error(err1))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	//当前登录的用户信息
	user, err := mysql.GetUserMessageById(userid)
	if err != nil {
		zap.L().Error("用户关注处理查询当前登录用户错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	//判断该用户是进行关注操作还是取消关注操作
	if p.AttentionorCencel == false && user.AttentionNum > 0 {
		//表明要去进行取消关注操作
		//进行添加关注操作
		//对数据库中的值进行操作
		user.AttentionNum = user.AttentionNum - 1
		err := mysql.UpdateUserMessage(user)
		if err != nil {
			zap.L().Error("修改数据库关注数失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		befocused.BeFocused = befocused.BeFocused - 1
		err = mysql.UpdateUserMessage(befocused)
		if err != nil {
			zap.L().Error("修改数据库被关注数失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		//添加谁被当前用户关注了,在被关注的用户下面添加
		str := "BeFocused:" + strconv.Itoa(int(p.ID))
		err = redis.REDIS.SRem(context.Background(), str, int(userid)).Err()
		if err != nil {
			zap.L().Error("用户关注处理redis错误", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		//添加当前登录用户关注了谁
		str1 := "Attention:" + strconv.Itoa(int(userid))
		err = redis.REDIS.SRem(context.Background(), str1, int(p.ID)).Err()
		if err != nil {
			zap.L().Error("用户关注处理redis错误", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}

	} else {
		//进行添加关注操作
		//对数据库中的值进行操作
		user.AttentionNum = user.AttentionNum + 1
		err := mysql.UpdateUserMessage(user)
		if err != nil {
			zap.L().Error("修改数据库关注数失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		befocused.BeFocused = befocused.BeFocused + 1
		err = mysql.UpdateUserMessage(befocused)
		if err != nil {
			zap.L().Error("修改数据库被关注数失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		//添加谁被当前用户关注了,在被关注的用户下面添加
		str := "BeFocused:" + strconv.Itoa(int(p.ID))
		err = redis.REDIS.SAdd(context.Background(), str, int(userid)).Err()
		if err != nil {
			zap.L().Error("用户关注处理redis错误", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		//添加当前登录用户关注了谁
		str1 := "Attention:" + strconv.Itoa(int(userid))
		err = redis.REDIS.SAdd(context.Background(), str1, int(p.ID)).Err()
		if err != nil {
			zap.L().Error("用户关注处理redis错误", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "关注失败")
			return
		}
		err = SendNotice(c, "likeYou", int(userid), int(p.ID), int(p.ID), "")
		if err != nil {
			zap.L().Error("发送通知失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
	}
	ResponseSuccess(c, CodeSuccess)
}

//判断用户是否关注过该用户
func IsAttention(userids []int, id int) (err error, isattention bool) {
	for _, user := range userids {
		if user == id {
			isattention = true
			return
		}
		isattention = false
	}
	return
}

//取出redis中该用户的关注列表
func GetAttentionIdByRedis(s string) (userids []int, err error) {
	attentionid, err2 := redis.REDIS.SMembers(context.Background(), s).Result()
	if err2 != nil {
		zap.L().Error("查询所有关注redis失败", zap.Error(err2))
		return
	}
	for i := 0; i < len(attentionid); i++ {
		id, err := strconv.Atoi(attentionid[i])
		if err != nil {
			zap.L().Error("查询所有关注redis转换失败", zap.Error(err))
			break
		}
		userids = append(userids, id)
	}
	return
}
