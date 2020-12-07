package util

type MyCode int64

const (
	CodeSuccess MyCode = 200
	CodeError   MyCode = 500

	// code= 1000... 用户模块的错误
	UserInvalidParams   MyCode = 1001
	UserExist           MyCode = 1002
	UserNotExist        MyCode = 1003
	UserInvalidPassword MyCode = 1004
	UserDataBaseError   MyCode = 1005
	UserNotLogin        MyCode = 1006
	UserTokenNotExist   MyCode = 1007
	UserTokenExpired    MyCode = 1008
	UserTokenWrong      MyCode = 1009

	// code= 2000... 问题模块的错误
	QuestionInvalidParams  MyCode = 2001
	QuestionDataBaseError MyCode = 2002
)

var msgFlags = map[MyCode]string{
	CodeSuccess: "OK",
	CodeError:   "FAIL",

	UserInvalidParams:   "用户请求参数错误",
	UserExist:           "用户名重复",
	UserNotExist:        "用户不存在",
	UserInvalidPassword: "用户名或密码错误",
	UserDataBaseError:   "用户数据库错误",
	UserNotLogin:        "用户未登录",
	UserTokenNotExist:   "TOKEN不存在",
	UserTokenExpired:    "TOKEN已过期",
	UserTokenWrong:      "TOKEN格式错误",

	QuestionInvalidParams:  "问题请求参数错误",
	QuestionDataBaseError: "问题数据库错误",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeError]
}
