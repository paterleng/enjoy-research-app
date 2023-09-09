package mysql

import "web_app/model"

func SearchProvince() (province []model.Province, err error) {
	err = DB.Find(&province).Error
	return province, err
}
func SearchCity(provinceid uint) (citys []model.City, err error) {
	err = DB.Where("province_id", provinceid).Find(&citys).Error
	return citys, err
}
func SearchSchool(cityid uint) (schools []model.School, err error) {
	err = DB.Where("city_id", cityid).Find(&schools).Error
	return schools, err
}

func SearchAcademies(schoolid uint) (academy []model.Academy, err error) {
	err = DB.Where("school_id", schoolid).Find(&academy).Error
	return academy, err
}

func SelectAllAddress() (err error, address []model.Province) {
	err = DB.Model(model.Province{}).Preload("Cities").Find(&address).Error
	return err, address
}

//
//func SelectLikeSchool() (err error, schools []model.School) {
//
//}
