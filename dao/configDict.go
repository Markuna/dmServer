package dao

import "log"

type ConfigDict struct {
	Id       int64  `gorm:"column:id"`
	ParentId int64  `gorm:"column:parent_id"`
	DKey     string `gorm:"column:d_key"`
	DValue   string `gorm:"column:d_value"`
	DSort    int64  `gorm:"column:d_sort"`
}

func (ConfigDict) TableName() string {
	return "dm_conf_dict"
}
func QueryDictByKey(key string) *[]ConfigDict {
	w := &[]ConfigDict{}
	r := DB.Where("d_key = ?", key).Order("d_sort ASC").Find(w)
	if r.Error != nil {
		log.Println(r.Error)
	}
	return w
}
