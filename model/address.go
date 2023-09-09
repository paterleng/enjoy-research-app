package model

import "gorm.io/gorm"

// 用于建数据库
type Province struct {
	gorm.Model
	Cities   []City
	Province string //省份名
}
type ProvinceParam struct {
	Province string //省份名
	Cities   []City `json:"children"`
}
type City struct {
	gorm.Model
	City       string //市名
	Schools    []School
	ProvinceID uint
}
type School struct {
	gorm.Model
	School    string //学校名称
	CityID    uint   //市名
	Academies []Academy
}

type Academy struct {
	gorm.Model
	Academy  string //学院名
	SchoolID uint   //学校id
	Majors   []Major
}

type Major struct {
	gorm.Model
	Major     string    //专业
	AcademyID uint      //学院id
	Num       uint      //招生人数
	Direction string    //研究方向
	Subjects  []Subject `gorm:"many2many:major_subjects;"` //考研科目
	Notes     string    //备注
}
type Subject struct {
	gorm.Model          //科目id
	SubjectID   uint    //考研科目id
	SubjectName string  //考研科目名称
	Majors      []Major `gorm:"many2many:major_subjects;"`
}

type MajorSubject struct {
	MajorID   uint `gorm:"primaryKey"`
	SubjectID uint `gorm:"primaryKey"`
	IsSelf    bool //0 1
}
