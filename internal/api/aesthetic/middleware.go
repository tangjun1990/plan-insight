package aesthetic

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 用户鉴权中间件
type AuthMiddleware struct {
	service *Service
}

// NewAuthMiddleware 创建鉴权中间件实例
func NewAuthMiddleware(service *Service) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
	}
}

func (m *AuthMiddleware) recordOpenAPIFailedCall(ctx *gin.Context, resultCode int) {
	_ = m.service.CreateOpenAPICallLog(0, "", ctx.FullPath(), 1, resultCode)
}

// UserAuth 微信小程序用户鉴权中间件
func (m *AuthMiddleware) UserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "请先登录",
			})
			return
		}

		// 处理Bearer Token
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")

		// 验证token
		user, err := m.service.GetUserByToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "无效的token: " + err.Error(),
			})
			return
		}

		// 将用户ID存入上下文
		ctx.Set("userID", user.ID)
		ctx.Set("userSource", user.Source)
		ctx.Next()
	}
}

// AdminAuth 管理员鉴权中间件
func (m *AuthMiddleware) AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "请先登录",
			})
			return
		}

		// 处理Bearer Token
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")

		// 验证token
		admin, err := m.service.GetAdminByToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "无效的token: " + err.Error(),
			})
			return
		}

		// 将管理员ID存入上下文
		ctx.Set("adminID", admin.ID)
		ctx.Next()
	}
}

// OpenAPIAuth 开放平台鉴权中间件
func (m *AuthMiddleware) OpenAPIAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			m.recordOpenAPIFailedCall(ctx, 401)
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "缺少开放平台鉴权信息",
			})
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")

		if err := m.service.ValidateOpenAccessToken(strings.TrimSpace(authHeader)); err != nil {
			m.recordOpenAPIFailedCall(ctx, 401)
			ctx.AbortWithStatusJSON(401, &Response{
				Code:    401,
				Message: "开放平台鉴权失败: " + err.Error(),
			})
			return
		}

		ctx.Next()
	}
}
