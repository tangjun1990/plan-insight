package aesthetic

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"git.4321.sh/feige/flygo/component/server/kin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRouter 注册审美感知应用的API路由
func RegisterRouter(server *kin.Component, db *gorm.DB) {
	// 创建服务和中间件实例
	service := NewService(db)
	controller := NewController(service)
	authMiddleware := NewAuthMiddleware(service)

	server.LoadHTMLGlob("./*.html")
	server.Static("/img", "./insimg")
	server.Static("/boximg", "./boximg")
	server.Static("/colorimg", "./colorimg")

	// 添加根路由，提供index.html页面
	server.GET("/", func(c *gin.Context) {

		indexImages := controller.service.GetIndexImage()
		imageSlice := make([]string, 0)
		for _, v := range indexImages {
			subImageSlice := make([]string, 0)
			for _, sv := range v.SubItems {
				subImageSlice = append(subImageSlice, fmt.Sprintf("{ id: '%s', url: '%s', categoryName: '%s' }", sv.Name, sv.URL, sv.CategoryName))
			}
			subImageString := strings.Join(subImageSlice, ",")

			imageSlice = append(imageSlice, fmt.Sprintf("%s: [%s]", v.CategoryEnglishName, subImageString))
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"color_map":    controller.service.GetIndexColor(),
			"word_map":     controller.service.GetIndexWord(),
			"image_string": template.JS("const imageData = {" + strings.Join(imageSlice, ",") + "};"),
		})
	})

	server.GET("/admin/login", func(c *gin.Context) {
		c.File("./admin/login.html")
	})
	server.GET("/admin/main.js", func(c *gin.Context) {
		c.File("./admin/main.js")
	})
	server.GET("/admin/index", func(c *gin.Context) {
		c.File("./admin/index.html")
	})
	server.GET("/admin/users", func(c *gin.Context) {
		c.File("./admin/users.html")
	})

	// 添加静态文件服务
	//server.Static("/static", "./web/static")

	// 注册微信小程序相关API
	wxGroup := server.Group("/api")
	{
		// 无需鉴权的接口
		wxGroup.POST("/wx/auth", controller.WxAuth)

		// 需要用户鉴权的接口
		userGroup := wxGroup.Group("/aesthetic")
		userGroup.Use(authMiddleware.UserAuth())
		{
			userGroup.POST("/data", controller.SaveAestheticData)
			userGroup.GET("/data/list", controller.GetUserAestheticDataList)
			userGroup.GET("/data/:id", controller.GetAestheticDataDetail)
		}

		// 用户信息相关接口
		wxUserGroup := wxGroup.Group("/wx/user")
		wxUserGroup.Use(authMiddleware.UserAuth())
		{
			wxUserGroup.GET("/info", controller.GetUserInfo)
			wxUserGroup.PUT("/update", controller.UpdateUserInfo)
		}
	}

	// 注册管理后台相关API
	adminGroup := server.Group("/admin")
	{
		// 无需鉴权的接口
		adminGroup.POST("/auth/login", controller.AdminLogin)

		// 需要管理员鉴权的接口
		authAdminGroup := adminGroup.Group("")
		authAdminGroup.Use(authMiddleware.AdminAuth())
		{

			// 用户管理
			authAdminGroup.GET("/user/list", controller.GetUserList)
			authAdminGroup.PUT("/user/:id/disable", controller.DisableUser)
			authAdminGroup.PUT("/user/:id/enable", controller.EnableUser)

			// 审美数据管理
			authAdminGroup.GET("/aesthetic/data/list", controller.GetAestheticDataList)
			authAdminGroup.GET("/aesthetic/data/analysis", controller.GetAestheticDataAnalysis)
		}
	}
}

// InitAdminUser 初始化管理员账号
func InitAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Admin{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果没有管理员账号，则创建一个默认管理员
	if count == 0 {
		admin := Admin{
			Phone:    "13800138000", // 默认管理员手机号
			Password: "admin123",    // 默认管理员密码
		}
		return db.Create(&admin).Error
	}

	return nil
}

// AutoMigrate 自动创建表结构
func AutoMigrate(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(&User{}, &AestheticData{}, &Admin{}); err != nil {
		return err
	}

	// 初始化管理员账号
	return InitAdminUser(db)
}
