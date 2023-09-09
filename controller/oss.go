package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hash"
	"io"
	"log"
	"net/http"
	"web_app/settings"
)

// 结构体
type Policy struct {
	Expiration string          `json:"expiration"`
	Conditions [][]interface{} `json:"conditions"`
}

func GetOSSToken(c *gin.Context) {
	// 生成签名代码
	var policy Policy
	policy.Expiration = "9999-12-31T12:00:00.000Z"
	var conditions []interface{}
	conditions = append(conditions, "content-length-range")
	conditions = append(conditions, 0)
	conditions = append(conditions, 1048576000)
	policy.Conditions = append(policy.Conditions, conditions)
	policyByte, err := json.Marshal(policy)
	if err != nil {
		log.Println("序列化失败", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": CodeServerBusy,
			"msg":  "上传失败",
		})
	}
	policyBase64 := base64.StdEncoding.EncodeToString(policyByte)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(settings.Conf.OSSConfig.AccessKeySecret))
	io.WriteString(h, policyBase64)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 主要拿的就是下面这两个玩意
	// 将这两个与前面获取的临时授权参数一起返回就好了
	token := struct {
		Host            string //上传地址
		DurationSeconds int    //过期时间
		Signature       string //签名
		Policy          string
		AccessKeyId     string
		AccessKeySecret string
	}{
		Host:            settings.Conf.OSSConfig.BasePath,
		DurationSeconds: 3600,
		Signature:       signature,
		Policy:          policyBase64,
		AccessKeyId:     settings.Conf.OSSConfig.AccessKeyId,
		AccessKeySecret: settings.Conf.OSSConfig.AccessKeySecret,
	}
	ResponseSuccess(c, token)
}
