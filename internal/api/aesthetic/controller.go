package aesthetic

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 审美感知应用API控制器
type Controller struct {
	service *Service
}

// NewController 创建控制器实例
func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

// ResponseSuccess 返回成功响应
func (c *Controller) ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ResponseError 返回错误响应
func (c *Controller) ResponseError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// GetUserID 从上下文中获取用户ID
func (c *Controller) GetUserID(ctx *gin.Context) uint {
	userID, _ := ctx.Get("userID")
	return userID.(uint)
}

// ======== 微信小程序用户相关 ========

// WxAuth 微信小程序用户鉴权
// @Summary 微信小程序用户鉴权
// @Description 基于微信小程序授权的手机号进行鉴权
// @Tags 小程序
// @Accept json
// @Produce json
// @Param body body WxAuthRequest true "鉴权信息"
// @Success 200 {object} Response{data=WxAuthResponse} "成功响应"
// @Router /api/wx/auth [post]
func (c *Controller) WxAuth(ctx *gin.Context) {
	var req WxAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.WxAuth(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "鉴权失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// SaveAestheticData 保存审美数据
// @Summary 保存审美数据
// @Description 用户提交审美数据表单
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param body body AestheticDataRequest true "审美数据表单"
// @Success 200 {object} Response{data=AestheticData} "成功响应，返回保存的数据"
// @Router /api/aesthetic/data [post]
func (c *Controller) SaveAestheticData(ctx *gin.Context) {
	var req AestheticDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUserID(ctx)
	data, err := c.service.SaveAestheticData(userID, &req)
	if err != nil {
		c.ResponseError(ctx, 500, "保存数据失败: "+err.Error())
		return
	}

	// 返回保存的数据
	c.ResponseSuccess(ctx, data)
}

// GetUserAestheticDataList 获取用户审美数据列表
// @Summary 获取用户审美数据列表
// @Description 小程序用户查看自己提交的审美数据列表
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页条数，默认10" default(10)
// @Success 200 {object} Response{data=PageResponse{list=[]AestheticData}} "成功响应"
// @Router /api/aesthetic/data/list [get]
func (c *Controller) GetUserAestheticDataList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	userID := c.GetUserID(ctx)
	resp, err := c.service.GetUserAestheticDataList(userID, page, pageSize)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetAestheticDataDetail 获取审美数据详情
// @Summary 获取审美数据详情
// @Description 小程序用户查看自己提交的审美数据详情
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param id path int true "审美数据ID"
// @Success 200 {object} Response{data=AestheticData} "成功响应"
// @Router /api/aesthetic/data/{id} [get]
func (c *Controller) GetAestheticDataDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.ResponseError(ctx, 400, "无效的ID参数")
		return
	}

	userID := c.GetUserID(ctx)
	data, err := c.service.GetAestheticDataDetail(uint(id), userID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, data)
}

// ======== 管理后台相关 ========

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 基于固定的手机号和密码进行登录验证
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param body body AdminLoginRequest true "登录信息"
// @Success 200 {object} Response{data=AdminLoginResponse} "成功响应"
// @Router /admin/auth/login [post]
func (c *Controller) AdminLogin(ctx *gin.Context) {
	var req AdminLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.AdminLogin(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "登录失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户数据列表，分页返回
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页条数，默认10" default(10)
// @Param phone query string false "手机号过滤"
// @Param status query int false "状态过滤，1正常，0禁用"
// @Success 200 {object} Response{data=PageResponse{list=[]User}} "成功响应"
// @Router /admin/user/list [get]
func (c *Controller) GetUserList(ctx *gin.Context) {
	var req UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resp, err := c.service.GetUserList(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户列表失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// DisableUser 禁用用户
// @Summary 禁用用户
// @Description 将单个用户改为禁用状态
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param id path int true "用户ID"
// @Success 200 {object} Response "成功响应"
// @Router /admin/user/{id}/disable [put]
func (c *Controller) DisableUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.ResponseError(ctx, 400, "无效的ID参数")
		return
	}

	if err := c.service.UpdateUserStatus(uint(id), 0); err != nil {
		c.ResponseError(ctx, 500, "禁用用户失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, nil)
}

// EnableUser 启用用户
// @Summary 启用用户
// @Description 将单个用户改为启用状态
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param id path int true "用户ID"
// @Success 200 {object} Response "成功响应"
// @Router /admin/user/{id}/enable [put]
func (c *Controller) EnableUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.ResponseError(ctx, 400, "无效的ID参数")
		return
	}

	if err := c.service.UpdateUserStatus(uint(id), 1); err != nil {
		c.ResponseError(ctx, 500, "启用用户失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, nil)
}

// GetAestheticDataList 获取审美数据列表
// @Summary 获取审美数据列表
// @Description 分页获取审美数据表中的数据
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页条数，默认10" default(10)
// @Param name query string false "姓名过滤"
// @Param gender query string false "性别过滤"
// @Param age_min query int false "最小年龄过滤"
// @Param age_max query int false "最大年龄过滤"
// @Param city query string false "所在城市过滤"
// @Param phone query string false "手机号过滤"
// @Success 200 {object} Response{data=PageResponse{list=[]AestheticData}} "成功响应"
// @Router /admin/aesthetic/data/list [get]
func (c *Controller) GetAestheticDataList(ctx *gin.Context) {
	var req AestheticDataListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resp, err := c.service.GetAestheticDataList(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取审美数据列表失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetAestheticDataAnalysis 获取审美数据统计分析
// @Summary 获取审美数据统计分析
// @Description 基于审美数据表中的数据，进行数据统计和分析
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param analysis_type query string true "分析类型: color(喜欢的颜色), disliked_color(讨厌的颜色), adjective(喜欢的形容词), image(喜欢的图片)"
// @Param dimension query string true "分析维度: count(数量), top(前N项), percent(百分比)"
// @Param top query int false "取前N条数据，默认10" default(10)
// @Param gender query string false "按性别过滤"
// @Param age_min query int false "最小年龄过滤"
// @Param age_max query int false "最大年龄过滤"
// @Param city query string false "按城市过滤"
// @Success 200 {object} Response{data=[]AnalysisItem} "成功响应"
// @Router /admin/aesthetic/data/analysis [get]
func (c *Controller) GetAestheticDataAnalysis(ctx *gin.Context) {
	var req AestheticAnalysisRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Top <= 0 {
		req.Top = 10
	}

	result, err := c.service.GetAestheticDataAnalysis(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取审美数据分析失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, result)
}
