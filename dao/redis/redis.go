package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"qa/config"
)

var client *redis.Client

// Init 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Conf.REDIS.Host, config.Conf.REDIS.Port),
		Password:     config.Conf.REDIS.Password, // no password set
		DB:           config.Conf.REDIS.DB,       // use default DB
		PoolSize:     config.Conf.REDIS.PoolSize,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
