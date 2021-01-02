package main

import (
	"flag"
	"fmt"
	"qa/config"
	"qa/dao"
	"qa/logic"
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
	if err := dao.InitDB(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer dao.DB.Close()

	// redis连接
	if err := dao.InitRDB(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer dao.RDB.Close()



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

	//异步记录question中view_count变化
	go logic.QuestionViewCount()
	//初始化question中view_count
	logic.InitQuestionViewCountChan<-1

	// 运行！！！
	if err := r.Run(config.Conf.Address); err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
