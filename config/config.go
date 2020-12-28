package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
	"os"
	"qa/util"
)

type Configuration struct {
	DB struct {
		Driver string `json:"driver"`
		Addr   string `json:"addr"`
	} `json:"db"`
	REDIS struct {
		Host         string `json:"host"`
		Password     string `json:"password"`
		Port         int    `json:"port"`
		DB           int    `json:"db"`
		PoolSize     int    `json:"pool_size"`
	} `json:"redis"`
	Address string `json:"address"`
	JwtKey  string `json:"jwtKey"`
	Mode    string `json:"mode"`

	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	Bucket      string `json:"bucket"`
	QiniuServer string `json:"qiniuServer"`
}

var Conf = new(Configuration)

// 加载配置信息
func Init() (err error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err = viper.ReadInConfig(); err != nil {
		os.Exit(1)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		os.Exit(1)
	}

	// 加载七牛云存储的配置文件
	file, err := ini.Load("config/qiniuConfig.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径：", err)
	}
	util.LoadQiniu(file)

	return
}
