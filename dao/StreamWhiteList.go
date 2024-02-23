package dao

import "log"

type WhiteList struct {
	Id  int64  `gorm:"column:id"`
	Uid string `gorm:"column:uid"`
}

func (WhiteList) TableName() string {
	return "dm_stream_white_list"
}

func QueryByUid(uid string) *WhiteList {
	w := &WhiteList{}
	r := DB.Where("uid = ?", uid).First(w)
	if r.Error != nil {
		log.Println(r.Error)
	}
	return w
}
