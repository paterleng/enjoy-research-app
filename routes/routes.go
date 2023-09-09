package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_app/controller"
	"web_app/middlewares"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	//处理跨域请求问题
	r.Use(middlewares.Core())
	v1 := r.Group("/api/v1")
	v1.POST("/register", controller.Register)
	v1.POST("/login", controller.Login)
	v1.GET("getMessageCode", controller.GetMessageCode)
	//应用jwt做是否登录验证
	v1.Use(middlewares.JWTAuthMiddleware())
	v1.GET("/getOSSToken", controller.GetOSSToken)
	u := v1.Group("/user")
	{
		//回显用户信息
		u.GET("/returnData", controller.ReturnData)
		//修改用户信息
		u.PUT("/updateUser", controller.UpdateUser)
		//用户注销
		u.DELETE("/deleteUser", controller.DeleteUser)
		//用户头像上传
		u.POST("/headPortrait", controller.HeadPortrait)
		//获得匿名名称，匿名头像
		u.GET("/getAnonymousName", controller.GetAnonymousName)
		//修改匿名
		u.POST("/updateAnonymousName", controller.UpdateAnonymousName)
		//用户关注或取消关注
		u.PUT("attentionUser", controller.AttentionUser)
	}
	post := v1.Group("/post")
	{
		//查看所有的说说
		post.GET("allPost", controller.AllPost)
		//文件用户头像
		post.POST("/myFile", controller.UpdateToOSS)
		post.POST("/createPost", controller.CreatePost)
		//查看自己说说
		post.GET("/myPost", controller.MyPost)
		//删除说说
		post.DELETE("/deletePost", controller.DeletePost)
		//修改说说
		post.PUT("updatePost", controller.UpdatePost)
		//说说点赞
		post.PUT("postLike", controller.PostLike)
		//说说收藏
		post.PUT("postCollect", controller.PostCollect)
		//帖子分类列表:根据科目分类，方向分类，学校分类
		post.GET("seletePostByClassify", controller.SeletePostByClassify)
		//创建帖子评论
		post.POST("createComment", controller.CreateComment)
		//查询该帖子下面的评论，一级包括二级评论
		post.GET("seleteSecondComment", controller.SeleteSecondComment)
		//删除评论
		post.DELETE("deleteComment", controller.DeleteComment)
	}
	treeHole := v1.Group("/treeHole")
	{
		//发表树洞
		treeHole.POST("/createTreeHole", controller.CreateTreeHole)
		//获取树洞
		treeHole.GET("/getTreeHole", controller.GetTreeHole)
		//删除树洞
		treeHole.DELETE("/deleteTreeHole", controller.DeleteTreeHole)
		//树洞回显以及一级评论
		treeHole.GET("/getTreeHoleByID", controller.GetTreeHoleByID)
		//修改树洞
		//treeHole.GET("/updateTreeHole", controller.UpdateTreeHole)
		//树洞点赞
		treeHole.POST("/likeTreeHole", controller.LikeTreeHole)
		//树洞评论
		treeHole.POST("/treeHoleComments", controller.TreeHoleComments)
		//查询树洞回复的评论
		treeHole.GET("/getTreeHoleComments", controller.GetTreeHoleComments)
		//查询树洞二级评论
		treeHole.GET("/getSecondTreeHoleComments", controller.GetSecondTreeHoleComments)
		//删除评论
		treeHole.DELETE("/deleteTreeHoleComment", controller.DeleteTreeHoleComment)
		//获取个人树洞
		treeHole.GET("/getTreeHoleByUserID", controller.GetTreeHoleByUserID)
	}
	address := v1.Group("/address")
	{
		//查找所有省
		address.GET("/searchProvince", controller.SearchProvince)
		//查找所有市
		address.GET("/searchCity", controller.SearchCity)
		//查找市下面的学校
		address.GET("/searchSchool", controller.SearchSchool)
		//查找学校下面的学院
		address.GET("/searchAcademies", controller.SearchAcademies)
		//查询所有的地址
		address.GET("/selectAllAddress", controller.SelectAllAddress)
		//模糊查询学校
		address.GET("/selectLikeSchool", controller.SelectLikeSchool)
		//根据学校查询该学校下面的所有学院
		address.GET("selectMajorBySchool", controller.SelectMajorBySchool)
		//根据学院id查询学院下专业
		address.GET("selectMajorByAcademyId", controller.SelectMajorByAcademyId)
		//根据专业id获得考试科目
		address.GET("seleteSubjectByAcademiyid", controller.SeleteSubjectByAcademiyid)
	}
	mySelf := v1.Group("/mySelf")
	{
		//获取我的喜欢
		mySelf.GET("/getMyLike", controller.GetMyLike)
		//获取我的收藏
		mySelf.GET("/getMyContent", controller.GetMyContent)
	}
	//第二版
	v2 := r.Group("/v2")

	Enter(v2)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
