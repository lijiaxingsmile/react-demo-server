package initApp

import (
	"log"
	"react-demo-server/config"
	"react-demo-server/database"
	"react-demo-server/db"
	"react-demo-server/route"

	"github.com/gin-gonic/gin"
)

func Start() {
	log.Println("初始化全局配置...")
	config.Start()

	log.Println("初始化数据库...")
	db.DB = database.Start()

	log.Println("初始化Web服务...")
	r := gin.Default()

	r.SetTrustedProxies(config.Config.TrustedProxies)
	gin.SetMode(config.Config.Env)

	log.Println("初始化路由...")
	route.Start(r)

	log.Println("启动服务...")
	r.Run("localhost:" + config.Config.Port)
}
