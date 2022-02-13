package api

import (
	"react-demo-server/model"
	"react-demo-server/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// func Login(c *gin.Context) {

// 	user := model.User{}

// 	err := c.ShouldBindJSON(&user)
// 	if err != nil {
// 		c.JSON(400, util.Response(false, "无效的登录信息", nil))
// 		return
// 	}

// 	if user.Username == "邢立佳" && user.Password == "xlj1401" {

// 		c.JSON(200, util.Response(true, "登录成功", gin.H{
// 			"token":    "123456",
// 			"remember": user.Remeber,
// 		}))
// 		return
// 	}

// 	c.JSON(200, util.Response(false, "用户名或密码错误", nil))
// }

func GetUser(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID := claims["id"].(float64)
	user, err := model.GetUser(uint(userID), false)
	user.Password = ""

	if err != nil {
		c.JSON(200, util.Response(false, "无效的用户信息", nil))
		return
	}

	c.JSON(200, util.Response(true, "获取用户信息成功", gin.H{
		"user": user,
	}))
}
