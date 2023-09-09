package controller

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strconv"
	"strings"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/model"
	"web_app/settings"
)

//查看所有的说说
func AllPost(c *gin.Context) {
	userid := GetUserIDByToken(c)
	//获取分页参数
	page, size := getPageInfo(c)
	posts, err := mysql.AllPost(int(page), int(size))
	if err != nil {
		zap.L().Error("查询所有帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	//从redis里面拿到我给谁点过赞
	str := "Like:" + strconv.Itoa(int(userid))
	postid, err := redis.REDIS.SMembers(context.Background(), str).Result()
	if err != nil {
		zap.L().Error("查询所有帖子redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var postids []int
	for i := 0; i < len(postid); i++ {
		id, err := strconv.Atoi(postid[i])
		if err != nil {
			zap.L().Error("查询所有帖子转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		postids = append(postids, id)
	}
	//拿到之后与数据库中查出来的做对比，查看是否有用户点过赞的帖子
	p := Intersection(postids, posts, 1)
	//展示的时候加上我收藏过哪篇帖子
	str1 := "Collent:" + strconv.Itoa(int(userid))
	collentpostid, err := redis.REDIS.SMembers(context.Background(), str1).Result()
	if err != nil {
		zap.L().Error("查询所有帖子redis失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	var collentpostids []int
	for i := 0; i < len(collentpostid); i++ {
		ids, err := strconv.Atoi(collentpostid[i])
		if err != nil {
			zap.L().Error("查询所有帖子转换失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
			return
		}
		collentpostids = append(collentpostids, ids)
	}
	p = Intersection(collentpostids, posts, 2)
	//封装好数据后返回给前端
	ResponseSuccess(c, p)
}

//求两个数组的交集
func Intersection(postids []int, posts []model.ShuoShuo, state int) (p []model.ShuoShuo) {
	var nowpostids []int
	for _, post := range posts {
		nowpostids = append(nowpostids, int(post.ID))
	}
	var res []int
	set1 := make(map[int]int)
	for _, v := range postids {
		//以数组中的值为键
		set1[v] = 1
	}
	for _, v := range nowpostids {
		if count, ok := set1[v]; ok && count > 0 {
			res = append(res, v)
			set1[v]--
		}
	}
	if state == 1 {
		for i := 0; i < len(posts); i++ {
			for _, id := range res {
				if int(posts[i].ID) == id {
					posts[i].IsLikeThisPost = true
					break
				}
			}
		}
	}
	if state == 2 {
		for i := 0; i < len(posts); i++ {
			for _, id := range res {
				if int(posts[i].ID) == id {
					posts[i].IsCollentThisPost = true
					break
				}
			}
		}
	}
	return posts
}

//创建一个说说
func CreatePost(c *gin.Context) {
	userid := GetUserIDByToken(c)
	p := new(model.ParamCreateShuos)
	err := c.ShouldBindJSON(&p)
	if err != nil {
		zap.L().Error("发表失败", zap.Error(err))
		ResponseError(c, CodeCreateError)
		return
	}
	shuoshuos := model.ShuoShuo{
		UserID:          userid,
		Tittle:          p.Tittle,
		Content:         p.Content,
		Address:         p.Address,
		LikeNum:         0,
		Files:           p.FileMsg,
		ClassIfications: p.ClassIfications,
	}
	err = mysql.CreateShuoShuo(&shuoshuos)
	if err != nil {
		zap.L().Error("插入数据库失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//查看自己的说说
func MyPost(c *gin.Context) {
	userid := GetUserIDByToken(c)
	err, shuos := mysql.MyShuoShuo(userid)
	if err != nil {
		zap.L().Error("查询发表资料失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, shuos)
}

//删除贴子
func DeletePost(c *gin.Context) {
	shuoId := c.Query("id")
	if strings.Trim(shuoId, "") == "" {
		zap.L().Error("参数为空")
		ResponseError(c, CodeInvalidParam)
		return
	}
	var p model.ShuoShuo
	id, err := strconv.Atoi(shuoId)
	if err != nil {
		zap.L().Error("参数错误")
	}
	p.ID = uint(id)
	//先判断用户有没有这个帖子
	post, err := mysql.GetPostById(p)
	if err != nil {
		zap.L().Error("说说查询失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "说说查询失败")
		return
	}
	if GetUserIDByToken(c) != post.UserID {
		zap.L().Error("用户没有权限", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "不可删除非自己的帖子")
		return
	}
	err = mysql.DeletePost(p)
	if err != nil {
		zap.L().Error("删除说说失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除说说失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//更新用户说说
func UpdatePost(c *gin.Context) {
	//先绑定参数
	var p model.ShuoShuo
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("更新说说绑定参数失败", zap.Error(err))
		ResponseErrorMsg(c, "参数不足")
		return
	}
	userid := GetUserIDByToken(c)
	post, err := mysql.GetPostById(p)
	if err != nil {
		zap.L().Error("数据库查询帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	if userid != post.UserID {
		zap.L().Error("不是该用户的帖子，该用户想去删除", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "不是你的帖子")
		return
	}
	//直接拿着前端传来的数据去数据库更新
	err = mysql.UpdatePost(p)
	if err != nil {
		zap.L().Error("数据库更新帖子失败", zap.Error(err))
		ResponseErrorMsg(c, "更新失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//用户说说点赞
func PostLike(c *gin.Context) {
	userid := GetUserIDByToken(c)
	var p model.ShuoShuo
	//接收前端传来的帖子id
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("更新说说绑定参数失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数不足")
		return
	}
	//查数据库中这个帖子属于谁
	post, err := mysql.GetPostById(p)
	if err != nil {
		zap.L().Error("查帖子信息错误", zap.Error(err))
		ResponseErrorMsg(c, "操作失败")
		return
	}
	//向数据库该帖子的点赞字段修改
	if p.AddOrCancelLike == false && post.LikeNum > 0 {
		//取消点赞
		post.LikeNum = post.LikeNum - 1
		//以被点赞用户的信息作键，当前用户id做值存入redis中,即记录谁给我点了赞，我给谁点了赞，点赞的是哪个帖子
		//Like+点赞用户id：[帖子id]
		//1.删除我给哪个帖子点赞了
		str := "Like:" + strconv.Itoa(int(userid))
		err := redis.REDIS.SRem(context.Background(), str, p.ID).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		//2.存储哪个人给我哪个帖子点赞了,谁给我点了赞
		//Like+当前帖子的用户id+当前帖子的id：[给这个帖子点赞的用户id]
		str1 := "Like:" + strconv.Itoa(int(post.UserID)) + ":" + strconv.Itoa(int(post.ID))
		err = redis.REDIS.SRem(context.Background(), str1, userid).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
	} else {
		//增加点赞
		post.LikeNum = post.LikeNum + 1
		//以被点赞用户的信息作键，当前用户id做值存入redis中,即记录谁给我点了赞，我给谁点了赞，点赞的是哪个帖子
		//Like+点赞用户id：[帖子id]
		//1.存储我给哪个帖子点赞了
		str := "Like:" + strconv.Itoa(int(userid))
		err = redis.REDIS.SAdd(context.Background(), str, p.ID).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		//2.存储哪个人给我哪个帖子点赞了,谁给我点了赞
		//Like+当前帖子的用户id+当前帖子的id：[给这个帖子点赞的用户id]
		str1 := "Like:" + strconv.Itoa(int(post.UserID)) + ":" + strconv.Itoa(int(post.ID))
		err = redis.REDIS.SAdd(context.Background(), str1, userid).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		//通过消息队列通知用户
		mq := NewSendRabbitMQ(post.UserID)
		SendData("我怕来了", mq)
		//向消息对列中取数据
		//rabbitMQ := NewReceiveRabbitMQ("->1")
		//ReceiveData(rabbitMQ, "我收到了")
		err := SendNotice(c, "likePost", int(userid), int(post.UserID), int(post.ID), "")
		if err != nil {
			zap.L().Error("发送通知失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
	}
	err = mysql.AddOrCencelLikeNum(&post)
	if err != nil {
		zap.L().Error("点赞量增加失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
		return
	}
	//向前端返回点赞后的数据
	ResponseSuccess(c, CodeSuccess)
}

//用户说说收藏
func PostCollect(c *gin.Context) {
	userid := GetUserIDByToken(c)
	var p model.ShuoShuo
	//接收前端传来的帖子id
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("更新说说绑定参数失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数不足")
		return
	}
	//查数据库中这个帖子属于谁
	post, err := mysql.GetPostById(p)
	if err != nil {
		zap.L().Error("查帖子信息错误", zap.Error(err))
		ResponseErrorMsg(c, "操作失败")
		return
	}
	//向数据库该帖子的收藏字段修改
	if p.IsCollectOrCencel == false && post.CollentNum > 0 {
		//取消点赞
		post.CollentNum = post.CollentNum - 1
		//以被点赞用户的信息作键，当前用户id做值存入redis中,即记录谁给我点了赞，我给谁点了赞，点赞的是哪个帖子
		//Collent+收藏用户id：[帖子id]
		//1.删除我给哪个帖子收藏了
		str := "Collent:" + strconv.Itoa(int(userid))
		err := redis.REDIS.SRem(context.Background(), str, p.ID).Err()
		if err != nil {
			zap.L().Error("收藏redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		//2.存储哪个人给我哪个帖子收藏了,谁给我收藏
		//BeCollent+当前帖子的用户id+当前帖子的id：[给这个帖子收藏的用户id]
		str1 := "BeCollent" + strconv.Itoa(int(post.UserID)) + ":" + strconv.Itoa(int(post.ID))
		err = redis.REDIS.SRem(context.Background(), str1, userid).Err()
		if err != nil {
			zap.L().Error("收藏redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
	} else {
		//增加收藏
		post.CollentNum = post.CollentNum + 1
		//以被点赞用户的信息作键，当前用户id做值存入redis中,即记录谁给我点了赞，我给谁点了赞，点赞的是哪个帖子
		//Collent+收藏用户id：[帖子id]
		//1.存储我给哪个帖子收藏了
		str := "Collent:" + strconv.Itoa(int(userid))
		err = redis.REDIS.SAdd(context.Background(), str, p.ID).Err()
		if err != nil {
			zap.L().Error("点赞redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		//2.存储哪个人给我哪个帖子收藏了,谁给我点了赞
		//BeCollent+当前帖子的用户id+当前帖子的id：[给这个帖子点赞的用户id]
		str1 := "BeCollent:" + strconv.Itoa(int(post.UserID)) + ":" + strconv.Itoa(int(post.ID))
		err = redis.REDIS.SAdd(context.Background(), str1, userid).Err()
		if err != nil {
			zap.L().Error("收藏redis操作失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
		err := SendNotice(c, "collectPost", int(userid), int(post.UserID), int(post.ID), "")
		if err != nil {
			zap.L().Error("发送通知失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
			return
		}
	}
	err = mysql.AddOrCencelLikeNum(&post)
	if err != nil {
		zap.L().Error("收藏量增加失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
		return
	}
	//向前端返回点赞后的数据
	ResponseSuccess(c, CodeSuccess)
}

//根据分类查询帖子
func SeletePostByClassify(c *gin.Context) {
	var p model.ClassIfication
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("根据条件查询帖子绑定参数失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数绑定失败")
		return
	}
	var postids []uint
	classification, err := mysql.SeletePostByClassify(p)
	if err != nil {
		zap.L().Error("根据条件查询帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "服务繁忙")
		return
	}
	for _, class := range classification {
		postids = append(postids, class.ShuoShuoID)
	}
	//根据帖子id查询帖子信息
	post, err := mysql.SeletePostById(postids)
	if err != nil {
		zap.L().Error("根据条件查询帖子失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "服务繁忙")
		return
	}
	ResponseSuccess(c, post)
}

//创建评论
func CreateComment(c *gin.Context) {
	//拿到当前登陆的用户id
	userid := GetUserIDByToken(c)
	var p *model.Comment
	//绑定参数
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("评论绑定参数失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数不足")
		return
	}
	//根据当前id获取根id
	var comment *model.Comment
	err := mysql.DB.Where("id", p.ParentID).Find(&comment).Error
	if err != nil {
		zap.L().Error("评论创建失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "评论创建失败")
		return
	}
	//判断他的父级评论是否含有根id，不含有则他的父评论就是一级评论
	if comment.RootID == 0 {
		p.RootID = comment.ID
	} else {
		p.RootID = comment.RootID
	}
	p.UserID = userid
	//创建评论
	err = mysql.CreateComment(p)
	if err != nil {
		zap.L().Error("评论创建失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "评论创建失败")
		return
	}
	//将该帖子的评论数++
	var post model.ShuoShuo
	post.ID = p.ShuoShuoID
	err = mysql.DB.First(&post).Error
	if err != nil {
		zap.L().Error("帖子信息失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "帖子信息失败")
		return
	}
	post.CommentNum = post.CommentNum + 1
	err = mysql.UpdatePost(post)
	if err != nil {
		zap.L().Error("评论数记录失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "评论数记录失败")
		return
	}
	if p.ParentID == 0 {
		err = SendNotice(c, "commentPost", int(userid), int(post.UserID), int(post.ID), p.Content)
	} else {
		err = SendNotice(c, "replyCommentsInPost", int(userid), int(post.UserID), int(post.ID), p.Content)
	}
	if err != nil {
		zap.L().Error("发送通知失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "操作失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//查看二级评论，只需要把数据返回
func SeleteSecondComment(c *gin.Context) {
	page, size := getPageInfo(c)
	id := c.Query("id")
	rootid, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("评论查询失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "评论查询失败")
		return
	}
	err, p := mysql.SeleteSecondComment(int(page), int(size), rootid)
	if err != nil {
		zap.L().Error("评论查询失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "评论查询失败")
		return
	}
	ResponseSuccess(c, p)
}

func DeleteComment(c *gin.Context) {
	userid := GetUserIDByToken(c)
	id := c.Query("commentid")
	commentid, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("删除帖子评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数错误")
		return
	}
	//根据帖子id查询帖子数据
	comment, err := mysql.SelectCommentById(uint(commentid))
	if err != nil {
		zap.L().Error("删除帖子评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询数据库失败")
		return
	}
	if userid != comment.UserID {
		zap.L().Error("非本用户操作", zap.Error(err))
		ResponseErrorWithMsg(c, CodeParamError, "你无权操作该评论")
		return
	}
	err = mysql.DeleteAllChildComment(comment)
	if err != nil {
		zap.L().Error("删除评论失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "删除失败")
		return
	}
	ResponseSuccess(c, CodeSuccess)
}

//上传头像文件至阿里云
func UpdateToOSS(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		zap.L().Error("接收文件失败", zap.Error(err))
		ResponseError(c, CodeUploadError)
		return
	}
	// 处理文件
	files := form.File["file"]
	//创建一个File类型的切片
	filemsg := make([]model.File, len(files))
	paths := make([]string, len(files))
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			zap.L().Error("接收文件失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeInvalidParam, "上传失败")
			return
		}
		//获得文件后缀
		uploadPath, err := UploadByOSS(f, file.Filename)
		paths = append(paths, uploadPath)
		//拿到路径后封装该文件的信息
		filemsg = append(filemsg, model.File{
			Name: file.Filename,
			Size: strconv.Itoa(int(file.Size)),
			Url:  uploadPath,
		})
		if err != nil {
			zap.L().Error("上传失败", zap.Error(err))
			ResponseErrorWithMsg(c, CodeInvalidParam, "上传失败")
			return
		}
	}
	ResponseSuccess(c, filemsg)
}

// OSS对象存储
func UploadByOSS(f io.Reader, fileName string) (s string, err error) {
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	endpoint := settings.Conf.OSSConfig.EndPoint
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := settings.Conf.OSSConfig.AccessKeyId
	accessKeySecret := settings.Conf.OSSConfig.AccessKeySecret

	// yourBucketName填写存储空间名称。
	bucketName := settings.Conf.OSSConfig.BucketName
	// uploadFileName填写文件上传的位置及名字。
	//在文件名后面加一个时间戳，防止用户重复上传相同的头像
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	url := fileNameStr + fileName
	uploadFileName := url
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.SecurityToken("abc"))
	if err != nil {
		zap.L().Error("上传失败", zap.Error(err))
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		zap.L().Error("上传失败", zap.Error(err))
		return
	}

	// 上传文件。
	err = bucket.PutObject("HeadPortrait"+"/"+uploadFileName, f)
	if err != nil {
		zap.L().Error("上传失败", zap.Error(err))
		return
	}
	uploadPath := settings.Conf.OSSConfig.BasePath + "/" + uploadFileName
	return uploadPath, err
}
