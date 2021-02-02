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
	QuestionDataBaseError  MyCode = 2002
	QuestionExist          MyCode = 2003
	QuestionNotExist       MyCode = 2004
	QuestionUserIdNotMatch MyCode = 2005

	// code= 3000... 个人简介模块的错误
	ProfileSaveFail       MyCode = 3001

	// code= 4000... 上传文件模块的错误
	UploadFail            MyCode = 4001

	// code= 5000... 回答模块的错误
	AnswerInvalidParams  MyCode = 5001
	AnswerDataBaseError  MyCode = 5002
	AnswerUserIdNotMatch MyCode = 5003

	// code= 6000... 投票模块的错误
	VoteInvalidParams  MyCode = 6001
	VoteDataBaseError  MyCode = 6002
)

var msgFlags = map[MyCode]string{
	CodeSuccess: "OK",
	CodeError:   "FAIL",

	UserInvalidParams:   "用户请求参数错误",
	UserExist:           "用户名重复",
	UserNotExist:        "用户不存在",
	UserInvalidPassword: "密码错误",
	UserDataBaseError:   "用户数据库错误",
	UserNotLogin:        "用户未登录",
	UserTokenNotExist:   "TOKEN不存在",
	UserTokenExpired:    "TOKEN已过期",
	UserTokenWrong:      "TOKEN格式错误",

	QuestionInvalidParams:  "问题请求参数错误",
	QuestionDataBaseError:  "问题数据库错误",
	QuestionExist:          "问题题目重复",
	QuestionNotExist:       "问题题目不存在",
	QuestionUserIdNotMatch: "提问者ID和当前用户ID不匹配",

	ProfileSaveFail:         "个人简介保存失败",

	UploadFail:              "上传失败",


	AnswerInvalidParams:  "回答请求参数错误",
	AnswerDataBaseError:  "回答数据库错误",
	AnswerUserIdNotMatch: "回答者ID和当前用户ID不匹配",

	VoteInvalidParams:"投票请求参数错误",
	VoteDataBaseError:"投票数据库错误",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeError]
}
