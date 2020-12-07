package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/config"
	"qa/util"
	"strings"
	"time"
)

var JwtKey = []byte(config.Conf.JwtKey)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string) (string, util.MyCode) {
	expireTime := time.Now().Add(10 * time.Hour)
	SetClaims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "qa",
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", util.CodeError
	}
	return token, util.CodeSuccess

}

// 验证token
func CheckToken(token string) (*MyClaims, util.MyCode) {
	var claims MyClaims

	setToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, e error) {
		return JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, util.UserTokenWrong
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, util.UserTokenExpired
			} else {
				return nil, util.UserTokenWrong
			}
		}
	}
	if setToken != nil {
		if key, ok := setToken.Claims.(*MyClaims); ok && setToken.Valid {
			return key, util.CodeSuccess
		} else {
			return nil, util.UserTokenWrong
		}
	}
	return nil, util.UserTokenWrong
}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code util.MyCode
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			code = util.UserTokenNotExist
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			})
			c.Abort()
			return
		}
		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) == 0 {
			code = util.UserTokenWrong
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			})
			c.Abort()
			return
		}

		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = util.UserTokenWrong
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			})
			c.Abort()
			return
		}
		key, tCode := CheckToken(checkToken[1])
		if tCode != util.CodeSuccess {
			code = tCode
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			})
			c.Abort()
			return
		}
		c.Set("username", key)
		c.Next()
	}
}
