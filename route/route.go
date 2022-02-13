package route

import (
	"react-demo-server/api"
	"react-demo-server/middleware"

	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine) {

	authMiddleware := middleware.JwtInit(r)

	v1 := r.Group("/v1")
	v1.GET("/ping", api.Ping)

	v1.POST("/login", authMiddleware.LoginHandler)

	// r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
	// 	claims := jwt.ExtractClaims(c)
	// 	log.Printf("NoRoute claims: %#v\n", claims)
	// 	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	// })

	auth := v1.Group("/auth")

	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/user", api.GetUser)
	}

}
