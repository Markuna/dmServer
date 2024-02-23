package service

import "douyinApi/dao"

type ConfigDict struct {
	Value string `json:"value"`
	Sort  int64  `json:"sort"`
}

func QueryDictByKey(key string) *[]ConfigDict {
	d := dao.QueryDictByKey(key)
	r := make([]ConfigDict, 0)
	for _, v := range *d {
		r = append(r, ConfigDict{
			Value: v.DValue,
			Sort:  v.DSort,
		})
	}
	return &r
}
