package dao

import (
	"douyinApi/config"
	"fmt"

	"gorm.io/driver/mysql" // gorm mysql 驱动包
	"gorm.io/gorm"         // gorm
)

var DB *gorm.DB

func init() {
	// MySQL 配置信息
	username := config.ConfigData.Mysql.User     // 账号
	password := config.ConfigData.Mysql.Password // 密码
	host := config.ConfigData.Mysql.Host         // 地址
	port := config.ConfigData.Mysql.Port         // 端口
	DBname := config.ConfigData.Mysql.Db         // 数据库名称
	timeout := "10s"                             // 连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
	// Open 连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql.")
	}
	DB = db
}
