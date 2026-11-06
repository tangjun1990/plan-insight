package aesthetic

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOpenAestheticDataListV3 开放平台获取用户审美报告列表（V3）
func (c *Controller) GetOpenAestheticDataListV3(ctx *gin.Context) {
	var req OpenAestheticListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.GetOpenAestheticDataListV3(&req)
	if err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 500, "获取审美报告列表失败: "+err.Error())
		return
	}

	c.responseOpenAPISuccess(ctx, req.Phone, resp)
}

// GetOpenUserSolutionListV3 开放平台获取用户方案列表（V3）
func (c *Controller) GetOpenUserSolutionListV3(ctx *gin.Context) {
	var req OpenSolutionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.service.GetOpenUserSolutionListV3(&req)
	if err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 500, "获取方案列表失败: "+err.Error())
		return
	}

	c.responseOpenAPISuccess(ctx, req.Phone, resp)
}

// OpenCreateUserSolutionV3 开放平台生成用户方案（V3）
func (c *Controller) OpenCreateUserSolutionV3(ctx *gin.Context) {
	var req OpenCreateSolutionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 400, "请求参数错误: "+err.Error())
		return
	}

	us, extra, err := c.service.OpenCreateUserSolutionV3(&req)
	if err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 500, "生成方案失败: "+err.Error())
		return
	}

	c.responseOpenAPISuccess(ctx, req.Phone, gin.H{
		"solution_id":   us.ID,
		"refresh_count": extra.RefreshCount,
		"solution_list": buildSolutionItemsFromExtraV3(extra.Images),
	})
}

// GetOpenUserSolutionDetailV3 开放平台获取用户方案详情（V3）
func (c *Controller) GetOpenUserSolutionDetailV3(ctx *gin.Context) {
	var req OpenSolutionDetailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 400, "请求参数错误: "+err.Error())
		return
	}

	data, err := c.service.GetOpenUserSolutionDetailV3(&req)
	if err != nil {
		c.responseOpenAPIError(ctx, req.Phone, 500, "获取方案详情失败: "+err.Error())
		return
	}

	c.responseOpenAPISuccess(ctx, req.Phone, data)
}

func normalizeOpenListPageV3(page, pageSize int) (int, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 20 {
		return 0, 0, errors.New("page_size最大不超过20")
	}
	return page, pageSize, nil
}

func (s *Service) getActiveOpenUserByPhoneV3(phone string) (*User, error) {
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	var user User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	return &user, nil
}

// GetOpenAestheticDataListV3 开放平台获取用户审美报告列表（V3）
// 逻辑参考 GetUserAestheticDataListV3：返回用户全部版本报告列表，由版本字段区分。
func (s *Service) GetOpenAestheticDataListV3(req *OpenAestheticListRequest) (*PageResponse, error) {
	phone := req.Phone
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	page, pageSize, err := normalizeOpenListPageV3(req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var user User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &PageResponse{
				List:     []*AestheticDataRsp{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}
		return nil, err
	}
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	return s.GetUserAestheticDataListV3(user.ID, page, pageSize)
}

// GetOpenUserSolutionListV3 开放平台获取用户方案列表（V3）
// 逻辑参考现有方案列表：返回用户全部方案列表，报告版本由返回数据中的 version 区分。
func (s *Service) GetOpenUserSolutionListV3(req *OpenSolutionListRequest) (*PageResponse, error) {
	page, pageSize, err := normalizeOpenListPageV3(req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	user, err := s.getActiveOpenUserByPhoneV3(req.Phone)
	if err != nil {
		return nil, err
	}

	return s.GetUserSolutionList(user.ID, page, pageSize)
}

// OpenCreateUserSolutionV3 开放平台生成用户方案（V3）
func (s *Service) OpenCreateUserSolutionV3(req *OpenCreateSolutionRequest) (*UserSolution, *UserSolutionExtraV3, error) {
	user, err := s.getActiveOpenUserByPhoneV3(req.Phone)
	if err != nil {
		return nil, nil, err
	}

	if req.AID <= 0 {
		return nil, nil, errors.New("无效的报告ID参数")
	}

	var report AestheticData
	if err := s.db.Where("id = ? AND user_id = ? AND version = ?", req.AID, user.ID, 3).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("v3报告不存在或不属于当前用户")
		}
		return nil, nil, err
	}

	return s.CreateUserSolutionV3(user.ID, uint(req.AID), 1)
}

// GetOpenUserSolutionDetailV3 开放平台获取用户方案详情（V3）
func (s *Service) GetOpenUserSolutionDetailV3(req *OpenSolutionDetailRequest) (*AestheticDataRspV3, error) {
	user, err := s.getActiveOpenUserByPhoneV3(req.Phone)
	if err != nil {
		return nil, err
	}

	if req.SolutionID <= 0 {
		return nil, errors.New("无效的solution_id参数")
	}

	return s.GetUserSolutionDetailV3(user.ID, uint(req.SolutionID))
}
