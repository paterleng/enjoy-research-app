package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/mysql"
	"web_app/model"
)

// 查找所有省
func SearchProvince(c *gin.Context) {
	//直接去查数据库
	province, err := mysql.SearchProvince()
	if err != nil {
		zap.L().Error("查询省出错了", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//成功返回省的信息
	ResponseSuccess(c, province)
}

// 查找所有市
func SearchCity(c *gin.Context) {
	province, _ := strconv.Atoi(c.Query("provinceid"))
	provinceid := uint(province)
	citys, err := mysql.SearchCity(provinceid)
	if err != nil {
		zap.L().Error("查询市失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, citys)

}

// 查找市下面的学校
func SearchSchool(c *gin.Context) {
	city, _ := strconv.Atoi(c.Query("cityid"))
	cityid := uint(city)
	citys, err := mysql.SearchSchool(cityid)
	if err != nil {
		zap.L().Error("查询学校失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, citys)
}

func SearchAcademies(c *gin.Context) {
	school, _ := strconv.Atoi(c.Query("schoolid"))
	schoolid := uint(school)
	academies, err := mysql.SearchAcademies(schoolid)
	if err != nil {
		zap.L().Error("查询专业失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, academies)
}

//查询所有的地址
func SelectAllAddress(c *gin.Context) {
	err, address := mysql.SelectAllAddress()
	if err != nil {
		zap.L().Error("查询所有地址", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, address)
}

//模糊查询学校
func SelectLikeSchool(c *gin.Context) {

}

//根据学校id查询该学校下面所有专业
func SelectMajorBySchool(c *gin.Context) {
	id := c.Query("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("查询学校下所有专业失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "参数错误")
		return
	}
	var academies []model.Academy
	var academiesid []uint
	err = mysql.DB.Where("school_id", ids).Find(&academies).Error
	if err != nil {
		zap.L().Error("查询学校下所有专业失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	for _, academy := range academies {
		academiesid = append(academiesid, academy.ID)
	}
	//去查学院下面的各个专业
	var majors []model.Major
	err = mysql.DB.Where("academy_id IN ?", academiesid).Find(&majors).Error
	if err != nil {
		zap.L().Error("查询学校下所有专业失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, majors)
}

//根据学院id查询学院下专业
func SelectMajorByAcademyId(c *gin.Context) {
	id := c.Query("id")
	academyid, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("查询学院下所有专业失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "参数错误")
		return
	}
	var p *[]model.Major
	err = mysql.DB.Where("academy_id", academyid).Find(&p).Error
	if err != nil {
		zap.L().Error("查询专业下所有科目失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	ResponseSuccess(c, p)
}

//根据专业id查询考试科目
func SeleteSubjectByAcademiyid(c *gin.Context) {
	id := c.Query("id")
	marjorid, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("查询专业下所有科目失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "参数错误")
		return
	}
	var p model.Major
	p.ID = uint(marjorid)
	err = mysql.DB.Preload("Subjects").Find(&p).Error
	if err != nil {
		zap.L().Error("查询专业下所有科目失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, "查询失败")
		return
	}
	ResponseSuccess(c, p)
}

//根据学校及专业id查询考研科目
