# Go-QAsystem [![Go](https://img.shields.io/badge/language-Go-00ADD8.svg)](https://github.com/golang/go)

## 设计目标

设计并实现一个类似知乎的问答系统，项目施工中...

User Stroy 1

- [x] 用户注册
- [x] 用户登录
- [x] 发布问题
- [x] 查看问题列表
- [x] 修改问题
- [x] 删除问题

User Story 2

- [x] 修改 profile
- [ ] 回答问题
- [ ] 修改回答
- [ ] 删除回答

## 技术选型

| 名称      | 地址                                   | 作用            |
| --------- | -------------------------------------- | --------------- |
| gin       | https://github.com/gin-gonic/gin       | web 框架        |
| gorm      | https://github.com/jinzhu/gorm         | orm 映射框架    |
| mysql     | https://github.com/go-sql-driver/mysql | 数据库          |
| viper     | https://github.com/spf13/viper         | 配置管理        |
| jwt-go    | https://github.com/dgrijalva/jwt-go    | jwt 认证        |
| snoyflake | https://github.com/sony/sonyflake      | 雪花算法唯一 ID |
| air       | https://github.com/cosmtrek/air        | 开发热加载      |

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

### 使用说明

1. 修改./config/config.yaml

```yaml
# 数据库
db:
  driver: mysql
  addr: yourusername:yourpassword@tcp(127.0.0.1:3306)/qasystem?charset=utf8mb4&parseTime=True&loc=Local

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

TODO
