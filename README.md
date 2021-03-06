# Go-QAsystem [![Go](https://img.shields.io/badge/language-Go-00ADD8.svg)](https://github.com/golang/go)

## 设计目标

设计并实现一个类似知乎的问答系统，项目基本完成

User Stroy 1

- [x] 用户注册
- [x] 用户登录
- [x] 发布问题
- [x] 查看问题列表
- [x] 修改问题
- [x] 删除问题

User Story 2

- [x] 修改 profile
- [x] 回答问题
- [x] 修改回答
- [x] 删除回答

User Story 3

- [x] 热门问题列表
- [x] 对回答点赞和踩
- [x] 个人中心查看自己已发布的问题
- [x] 个人中心查看自己的回答
- [x] 个人中心查看自己点赞的回答

## 技术选型

| 名称      | 地址                                           | 作用            |
| --------- | ---------------------------------------------- | --------------- |
| gin       | https://github.com/gin-gonic/gin               | web 框架        |
| gorm      | https://github.com/jinzhu/gorm                 | orm 映射框架    |
| mysql     | https://github.com/go-sql-driver/mysql         | mysql 数据库    |
| viper     | https://github.com/spf13/viper                 | 配置管理        |
| jwt-go    | https://github.com/dgrijalva/jwt-go            | jwt 认证        |
| snoyflake | https://github.com/sony/sonyflake              | 雪花算法唯一 ID |
| air       | https://github.com/cosmtrek/air                | 开发热加载      |
| logrus    | https://github.com/sirupsen/logrus             | 日志处理        |
| bcrypt    | https://golang.org/x/crypto/bcrypt             | 密码处理        |
| validator | https://github.com/go-playground/validator/v10 | 数据校验        |
| go-redis  | https://github.com/go-redis/redis              | redis 数据库    |

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
git clone https://github.com/ustczhihu/qa.git
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

1. 修改./config/config.yaml

```yaml
# 数据库
db:
  driver: mysql
  addr: yourusername:yourpassword@tcp(127.0.0.1:3306)/qasystem?charset=utf8mb4&parseTime=True&loc=Local

# redis
redis:
  host: localhost
  port: 6379
  password: yourpassword
  db: 0
  pool_size: 100

# jwt认证密钥
jwtKey: yourjwtkey

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

增加 air 热加载功能，运行程序可使用

```shell
air
```

## 效果图

![homepage](/preview/homepage.png)

![profile](/preview/profile.png)

![profile1](/preview/profile1.png)

![profile2](/preview/profile2.png)

![question](/preview/question.png)

![answer](/preview/answer.png)

![answer2](/preview/answer2.png)

![answer3](/preview/answer3.png)
