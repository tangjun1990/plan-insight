package aesthetic

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tangjun1990/flygo/core/kcfg"
)

// Controller 审美感知应用API控制器
type Controller struct {
	service                  *Service
	saveAestheticDataLimiter chan struct{}
}

// NewController 创建控制器实例
func NewController(service *Service) *Controller {
	return &Controller{
		service:                  service,
		saveAestheticDataLimiter: make(chan struct{}, 1),
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

func (c *Controller) GetUserSource(ctx *gin.Context) int {
	return ctx.GetInt("userSource")
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

// GetOpenUserAuthToken 开放平台获取用户token
func (c *Controller) GetOpenUserAuthToken(ctx *gin.Context) {
	var req OpenUserAuthTokenRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.IssueOpenUserToken(req.Phone, req.WxOpenID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户token失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetOpenAestheticDataList 开放平台获取用户审美报告列表
func (c *Controller) GetOpenAestheticDataList(ctx *gin.Context) {
	var req OpenAestheticListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.GetOpenAestheticDataList(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取审美报告列表失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetOpenUserSolutionList 开放平台获取用户方案列表
func (c *Controller) GetOpenUserSolutionList(ctx *gin.Context) {
	var req OpenSolutionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.GetOpenUserSolutionList(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取方案列表失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// OpenCreateUserSolution 开放平台生成用户方案
func (c *Controller) OpenCreateUserSolution(ctx *gin.Context) {
	var req OpenCreateSolutionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	us, err := c.service.OpenCreateUserSolution(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "生成方案失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, gin.H{"solution_id": us.ID})
}

// GetOpenUserSolutionDetail 开放平台获取用户方案详情
func (c *Controller) GetOpenUserSolutionDetail(ctx *gin.Context) {
	var req OpenSolutionDetailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	data, err := c.service.GetOpenUserSolutionDetail(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取方案详情失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, data)
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

	// 本地串行化高开销的保存流程，避免瞬时高并发压垮当前服务进程。
	select {
	case c.saveAestheticDataLimiter <- struct{}{}:
		defer func() {
			<-c.saveAestheticDataLimiter
		}()
	case <-ctx.Request.Context().Done():
		c.ResponseError(ctx, 499, "请求已取消，请稍后重试")
		return
	}

	userID := c.GetUserID(ctx)
	data, err := c.service.SaveAestheticData(userID, c.GetUserSource(ctx), &req)
	if err != nil {
		c.ResponseError(ctx, 500, "保存数据失败: "+err.Error())
		return
	}

	// 返回保存的数据
	c.ResponseSuccess(ctx, data)
}

func (c *Controller) AddCollection(ctx *gin.Context) {
	var req AddCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUserID(ctx)
	err := c.service.AddCollection(userID, &req)
	if err != nil {
		c.ResponseError(ctx, 500, "保存数据失败: "+err.Error())
		return
	}

	// 返回保存的数据
	c.ResponseSuccess(ctx, gin.H{})
}

func (c *Controller) CancelCollection(ctx *gin.Context) {
	var req CancelCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUserID(ctx)
	err := c.service.CancelCollection(userID, &req)
	if err != nil {
		c.ResponseError(ctx, 500, "保存数据失败: "+err.Error())
		return
	}

	// 返回保存的数据
	c.ResponseSuccess(ctx, gin.H{})
}

// CreateUserSolution 生成用户方案
// @Summary 生成用户方案
// @Description 基于审美报告ID创建用户方案，避免重复
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param body body CreateSolutionRequest true "生成方案请求体"
// @Success 200 {object} Response{data=map[string]interface{}} "成功响应"
// @Router /api/aesthetic/solution/create [post]
func (c *Controller) CreateUserSolution(ctx *gin.Context) {
	var req CreateSolutionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUserID(ctx)
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户信息失败: "+err.Error())
		return
	}
	if !c.service.IsProActive(user) {
		c.ResponseError(ctx, 403, "您未开通pro版，请联系PLAN客服开通！")
		return
	}

	aid, convErr := strconv.Atoi(req.AID)
	if convErr != nil || aid <= 0 {
		c.ResponseError(ctx, 400, "无效的aid参数")
		return
	}

	us, createErr := c.service.CreateUserSolution(userID, uint(aid))
	if createErr != nil {
		c.ResponseError(ctx, 500, createErr.Error())
		return
	}

	c.ResponseSuccess(ctx, gin.H{"solution_id": us.ID})
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

func (c *Controller) GetUserCollectionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	userID := c.GetUserID(ctx)
	resp, err := c.service.GetUserCollectionList(userID, page, pageSize)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetUserSolutionList 获取用户方案列表
// @Summary 获取用户方案列表
// @Description 小程序用户查看自己的方案列表（基于收藏与审美数据关联）
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页条数，默认10" default(10)
// @Success 200 {object} Response{data=PageResponse{list=[]AestheticData}} "成功响应"
// @Router /api/aesthetic/solution/list [get]
func (c *Controller) GetUserSolutionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	userID := c.GetUserID(ctx)
	// Pro准入校验
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户信息失败: "+err.Error())
		return
	}
	if !c.service.IsProActive(user) {
		c.ResponseError(ctx, 403, "您未开通pro版，请联系PLAN客服开通！")
		return
	}
	resp, err := c.service.GetUserSolutionList(userID, page, pageSize)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// GetUserSolutionDetail 获取用户方案详情
// @Summary 获取用户方案详情
// @Description 通过solution_id查询user_solutions拿到aid，再返回审美报告详情
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param solution_id query int true "方案ID"
// @Success 200 {object} Response{data=AestheticDataRsp} "成功响应"
// @Router /api/aesthetic/solution/detail [get]
func (c *Controller) GetUserSolutionDetail(ctx *gin.Context) {
	sidStr := ctx.Query("solution_id")
	sid, err := strconv.Atoi(sidStr)
	if err != nil || sid <= 0 {
		c.ResponseError(ctx, 400, "无效的solution_id参数")
		return
	}

	userID := c.GetUserID(ctx)
	// Pro准入校验
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户信息失败: "+err.Error())
		return
	}
	if !c.service.IsProActive(user) {
		c.ResponseError(ctx, 403, "您未开通pro版，请联系PLAN客服开通！")
		return
	}
	data, err := c.service.GetUserSolutionDetail(userID, uint(sid))
	if err != nil {
		c.ResponseError(ctx, 500, "获取方案详情失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, data)
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
	data, err := c.service.GetAestheticDataDetail(uint(id), userID, false)
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

// OpenUserPro 管理端开通Pro
// @Summary 开通Pro
// @Description 管理员为指定用户开通Pro，支持30/90/180/365天
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param id path int true "用户ID"
// @Param body body OpenProRequest true "开通时长（天）"
// @Success 200 {object} Response "成功响应"
// @Router /admin/user/{id}/pro/open [post]
func (c *Controller) OpenUserPro(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.ResponseError(ctx, 400, "无效的ID参数")
		return
	}

	var req OpenProRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	if err := c.service.UpdateUserPro(uint(id), req.Days); err != nil {
		c.ResponseError(ctx, 500, "开通Pro失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, nil)
}

// BatchOpenUserPro 管理端批量开通Pro
// @Summary 批量开通Pro
// @Description 管理员按手机号批量为用户开通Pro，支持30/90/180/365天
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param body body BatchOpenProRequest true "批量开通请求"
// @Success 200 {object} Response{data=BatchOpenProResponse} "成功响应"
// @Router /admin/user/pro/open/batch [post]
func (c *Controller) BatchOpenUserPro(ctx *gin.Context) {
	var req BatchOpenProRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.BatchUpdateUserPro(req.Phones, req.Days)
	if err != nil {
		c.ResponseError(ctx, 500, "批量开通Pro失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
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

// GetDashboardOverview 获取管理后台看板数据
// @Summary 获取管理后台看板数据
// @Description 获取近7/15/30天用户新增量与审美数据新增量，按天聚合并区分来源
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param Authorization header string true "管理员令牌"
// @Param days query int false "日期范围，仅支持7/15/30" default(7)
// @Success 200 {object} Response{data=DashboardOverviewResponse} "成功响应"
// @Router /admin/dashboard/overview [get]
func (c *Controller) GetDashboardOverview(ctx *gin.Context) {
	var req DashboardOverviewRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.Days <= 0 {
		req.Days = 7
	}

	resp, err := c.service.GetDashboardOverview(&req)
	if err != nil {
		c.ResponseError(ctx, 500, "获取看板数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

// UpdateUserInfo 更新用户信息
func (c *Controller) UpdateUserInfo(ctx *gin.Context) {
	var req UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUserID(ctx)
	if userID == 0 {
		c.ResponseError(ctx, 401, "未授权")
		return
	}

	err := c.service.UpdateUserInfo(userID, req)
	if err != nil {
		c.ResponseError(ctx, 500, "更新用户信息失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, nil)
}

// GetUserInfo 获取用户信息
func (c *Controller) GetUserInfo(ctx *gin.Context) {
	userID := c.GetUserID(ctx)
	if userID == 0 {
		c.ResponseError(ctx, 401, "未授权")
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		c.ResponseError(ctx, 500, "获取用户信息失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, user)
}

// GetImageList 获取图片数据
// @Summary 获取图片数据
// @Description 获取所有类别的图片数据
// @Tags 小程序
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Success 200 {object} Response{data=map[string]interface{}} "成功响应"
// @Router /api/aesthetic/images [get]
func (c *Controller) GetImageList(ctx *gin.Context) {
	gender := ctx.Query("gender")
	indexImages := c.service.GetIndexImage(gender)
	imageData := make(map[string][]map[string]interface{})

	for _, category := range indexImages {
		images := make([]map[string]interface{}, 0)
		for _, subItem := range category.SubItems {
			images = append(images, map[string]interface{}{
				"id":           subItem.Name,
				"url":          kcfg.GetString("app.global.host") + subItem.URL,
				"categoryName": subItem.CategoryName,
			})
		}
		imageData[category.CategoryEnglishName] = images
	}

	c.ResponseSuccess(ctx, imageData)
}

func (c *Controller) GetAllImage(ctx *gin.Context) {
	images := c.service.GetAllImage()
	imagesRsp := make([]string, 0)
	for _, imgname := range images {
		imagesRsp = append(imagesRsp, kcfg.GetString("app.global.host")+"/imgv2/"+imgname)
	}
	c.ResponseSuccess(ctx, imagesRsp)
}

func (c *Controller) GetColorList(ctx *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	colors := c.service.GetIndexColor()
	// 将colors中元素的顺序随机打乱
	//rand.Shuffle(len(colors), func(i, j int) {
	//	colors[i], colors[j] = colors[j], colors[i]
	//})

	c.ResponseSuccess(ctx, colors)
}

func (c *Controller) GetWordList(ctx *gin.Context) {
	words := c.service.GetIndexWord()
	c.ResponseSuccess(ctx, words)
}

func (c *Controller) GetCityList(ctx *gin.Context) {
	citys := c.service.GetAllCity()
	c.ResponseSuccess(ctx, citys)
}

func (c *Controller) CertQuery(ctx *gin.Context) {
	var req CertQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "参数错误: "+err.Error())
		return
	}
	certlevel := "ncd2"
	if strings.Contains(req.CertType, "三级") {
		certlevel = "ncd3"
	} else if strings.Contains(req.CertType, "四级") {
		certlevel = "ncd4"
		req.IDLastFour = "0000"
	}

	userCertName := fmt.Sprintf("%s-%s-%s", req.Name, req.IDLastFour, certlevel)
	// 获取证书目录 ./ncdcert下的全部文件
	files, err := ioutil.ReadDir("./ncdcert")
	if err != nil {
		c.ResponseError(ctx, 500, "获取证书目录失败: "+err.Error())
		return
	}
	// 遍历文件列表
	certfile := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), userCertName) {
			certfile = file.Name()
			break
		}
	}
	if certfile == "" {
		c.ResponseSuccess(ctx, gin.H{
			"cert_url": "",
		})
	} else {
		c.ResponseSuccess(ctx, gin.H{
			"cert_url": "https://plan-living.com/ncdcert/" + certfile,
		})
	}
}
