# Go-QAsystem [![Go](https://img.shields.io/badge/language-Go-00ADD8.svg)](https://github.com/golang/go)

## 设计目标

设计并实现一个类似知乎的问答系统，项目施工中...

- [x] 用户注册
- [x] 用户登录
- [x] 发布问题
- [x] 查看问题列表
- [x] 集成 jwt
- [x] 集成 snowflake 雪花算法
- [ ] 用户修改 profile
- [ ] 用户修改删除问题

## 技术选型

1. web:[gin](https://github.com/gin-gonic/gin)
2. orm:[gorm](https://github.com/jinzhu/gorm)
3. database:[mysql](https://github.com/go-sql-driver/mysql)
4. config manager:[viper](https://github.com/spf13/viper)

## 项目结构

```
-qa
    |-config     项目配置
    |-controller 控制器
    |-dao        数据库连接
    |-middleware 中间件
    |-model      模型
    |-router     路由
    |-util       公用工具
    |-go.mod     项目依赖
    |-initDB.go  数据库初始化文件
    |-main.go    程序执行入口
```

## 安装部署

```shell
git clone https://github.com/funtowin/qa.git
```

本项目使用`go module`管理依赖包，运行项目前请先安装依赖

```shell
go mod tidy
```

如遇下载过慢，请先运行如下命令

```shell
go env -w GOPROXY=https://goproxy.cn,direct
```

## 使用方法

### 使用说明

1. 修改./config/config.yaml

```yaml
# 数据库
db:
  driver: mysql
  addr: root:root@tcp(127.0.0.1:3306)/qasystem?charset=utf8mb4&parseTime=True&loc=Local

# jwt认证密钥
jwtKey: salt20201206

# 端口
address: :9090
```

2. 初次使用需初始化数据库

```shell
go run ./ -initDB
```

3. 之后使用可直接

```shell
go run ./
```

## 效果图

TODO
