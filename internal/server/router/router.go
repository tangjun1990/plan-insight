package router

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tangjun1990/flygo/component/server/kin"
	"github.com/tangjun1990/plan-insight/docs"
	"github.com/tangjun1990/plan-insight/internal/api/aesthetic"
	"github.com/tangjun1990/plan-insight/internal/server/middleware"
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
