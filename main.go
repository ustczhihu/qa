package main

import (
	"flag"
	"fmt"
	"qa/config"
	"qa/dao/mysql"
	"qa/dao/redis"
	"qa/router"
	util "qa/util"
)

func main() {
	// 加载配置文件
	if err := config.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	// 雪花算法
	if err := util.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 数据库连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.DB.Close() // 程序退出关闭数据库连接

	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 是否初始化数据库
	var shouldInitDB = flag.Bool("initDB", false, "initialize database")
	flag.Parse()
	if *shouldInitDB {
		if err := Init(); err != nil {
			fmt.Printf("init database failed, err:%v\n", err)
			return
		}
	}

	// 路由
	r := router.Init()

	// 运行！！！
	if err := r.Run(config.Conf.Address); err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
