package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取当前用户的id
func getCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

//获取分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	//获取分页参数
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = 0
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
