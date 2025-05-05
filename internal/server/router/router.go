package router

import (
	"git.4321.sh/feige/commonapi/docs"
	"git.4321.sh/feige/commonapi/internal/api/aesthetic"
	"git.4321.sh/feige/commonapi/internal/server/middleware"
	"git.4321.sh/feige/flygo/component/server/kin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func AddRouter(server *kin.Component, db *gorm.DB) {
	server.Use(middleware.Recovery) //

	// 注册审美感知应用API路由
	aesthetic.RegisterRouter(server, db)

	//Swagger
	docs.SwaggerInfo.BasePath = ""
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
