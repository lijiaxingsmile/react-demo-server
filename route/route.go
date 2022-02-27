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

	auth := v1.Group("/auth")

	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/user", api.GetUser)

		auth.GET("/menus", api.GetMenus)
		auth.GET("/menu", api.GetMenu)
		auth.POST("/menu", api.CreateMenu)
		auth.DELETE("/menus", api.BatchDeleteMenus)
		auth.PUT("/menu", api.UpdateMenu)
	}

}
