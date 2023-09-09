package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"image/color"
	"web_app/dao/mysql"
	"web_app/model"
)

//获取自己唯一的二维码
func CreateQRCode(c *gin.Context) {
	userid := GetUserIDByToken(c)
	//根据用户id查询用户信息
	user, err := mysql.ReturnDataMysql(4)
	if err != nil {
		zap.L().Error("查询用户信息失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "生成失败")
		return
	}
	codeMessage := model.RQCodeMessage{
		ID:       int(userid),
		UserName: user.Username,
		Mobile:   user.Mobile,
		//HeadPortrait: user.HeadPortrait[0].Url,
	}
	m, err := json.Marshal(codeMessage)
	if err != nil {
		zap.L().Error("序列化失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "生成失败")
		return
	}
	size := 256
	b := color.RGBA{102, 255, 255, 255} //外
	f := color.RGBA{0, 153, 85, 255}    //内
	qrcode.WriteColorFile(string(m), qrcode.High, size, b, f, "./a.png")
}
