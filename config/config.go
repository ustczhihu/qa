package config

import (
	"github.com/spf13/viper"
	"os"
)

type Configuration struct {
	DB struct {
		Driver string `json:"driver"`
		Addr   string `json:"addr"`
	} `json:"db"`
	Address string `json:"address"`
	JwtKey  string `json:"jwtKey"`
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
	return
}
