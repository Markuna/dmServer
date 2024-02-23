package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Http struct {
			Port string
		}
	}
	Mysql struct {
		Host     string
		User     string
		Password string
		Db       string
		Port     int32
	}

	Redis struct {
		Addr     string
		Password string
		Db       int32
	}
	Rabbitmq struct {
		Addr string
	}

	Douyin struct {
		Url            string
		ToutiaoUrl     string
		AppId          string
		AppSecret      string
		CallbackSecret string
		SkipAuth       bool
	}

	ForTest bool
}

var ConfigData *Config

func init() {
	viper.SetConfigType("yml")
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&ConfigData)
	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	return ConfigData
}
