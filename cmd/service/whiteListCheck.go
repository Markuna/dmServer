package service

import "douyinApi/dao"

func CheckUidIsInWhiteList(uid string) bool {
	r := dao.QueryByUid(uid)
	if r != nil && r.Uid != "" {
		return true
	}
	return false
}
