package middleware

import (
	"log"
	"strings"
	"time"

	"react-demo-server/db"
	"react-demo-server/model"
	"react-demo-server/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JwtInit(r *gin.Engine) *jwt.GinJWTMiddleware {
	log.Println("初始化JWT中间件...")

	var identityKey = "username"

	type login struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "demo",
		Key:         []byte("demo secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// 上下文
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
					"id":        v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				UserName: claims[identityKey].(string),
			}
		},
		// 登陆验证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userRecord := &model.User{
				UserName: loginVals.Username,
				Password: loginVals.Password,
			}

			err := db.DB.Take(&userRecord).Error
			if err == gorm.ErrRecordNotFound {
				return nil, jwt.ErrFailedAuthentication
			}

			return userRecord, nil
		},
		// 权限验证
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok {
				return true
			}
			return false
		},
		// 登录成功后的回调
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(200, util.Response(true, "登录成功", gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			}))
		},
		// 验证失败
		Unauthorized: func(c *gin.Context, code int, message string) {
			if strings.HasSuffix(c.Request.URL.Path, "/login") {
				c.JSON(200, util.Response(false, "用户名或密码错误", nil))
				return
			}
			c.JSON(code, util.Response(false, message, nil))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic(err)
	}

	return authMiddleware
}
