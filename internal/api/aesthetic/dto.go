package aesthetic

// 响应基础结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ======== 微信小程序用户相关 ========

// 微信小程序鉴权响应
type WxAuthResponse struct {
	Token     string      `json:"token"`               // 用户token
	ExpiresIn int64       `json:"expires_in"`          // token有效期（秒）
	UserInfo  interface{} `json:"user_info,omitempty"` // 用户信息
}

// 开放平台获取用户token请求
type OpenUserAuthTokenRequest struct {
	Phone    string `form:"phone" binding:"required,len=11"` // 用户手机号
	WxOpenID string `form:"wx_open_id" binding:"required"`   // 微信小程序用户ID
}

// 开放平台用户token详情
type OpenUserTokenItem struct {
	Token    string `json:"token"`     // 用户token
	ExpireAt int64  `json:"expire_at"` // 过期时间戳（秒）
}

// 开放平台获取用户token响应
type OpenUserAuthTokenResponse struct {
	UserToken OpenUserTokenItem `json:"user_token"`
}

// 开放平台获取审美报告列表请求
type OpenAestheticListRequest struct {
	Phone    string `form:"phone" binding:"required,len=11"`            // 用户手机号
	Page     int    `form:"page" binding:"omitempty,min=1"`             // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=20"` // 每页条数
}

// 开放平台生成方案请求
type OpenCreateSolutionRequest struct {
	Phone string `json:"phone" binding:"required,len=11"` // 用户手机号
	AID   int    `json:"report_id" binding:"required"`    // 报告ID
}

// 开放平台获取方案详情请求
type OpenSolutionDetailRequest struct {
	Phone      string `form:"phone" binding:"required,len=11"`      // 用户手机号
	SolutionID int    `form:"solution_id" binding:"required,min=1"` // 方案ID
}

// 开放平台获取方案列表请求
type OpenSolutionListRequest struct {
	Phone    string `form:"phone" binding:"required,len=11"`            // 用户手机号
	Page     int    `form:"page" binding:"omitempty,min=1"`             // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=20"` // 每页条数
}

// ======== 审美数据相关 ========

// 审美数据请求
type AestheticDataRequest struct {
	Name            string   `json:"name" binding:"required"`                          // 姓名
	Gender          string   `json:"gender" binding:"required"`                        // 性别
	Age             int      `json:"age" binding:"required,min=1,max=120"`             // 年龄
	City            string   `json:"city" binding:"required"`                          // 所在城市
	Phone           string   `json:"phone" binding:"required"`                         // 手机号码
	LikedColors     []string `json:"liked_colors" binding:"required,min=1,max=10"`     // 喜欢的颜色
	DislikedColors  []string `json:"disliked_colors" binding:"required,min=1,max=5"`   // 讨厌的颜色
	LikedAdjectives []string `json:"liked_adjectives" binding:"required,min=1,max=10"` // 喜欢的形容词
	LikedImages     []string `json:"liked_images" binding:"required,min=1,max=12"`     // 喜欢的图片
}

type AddCollectionRequest struct {
	AID string `json:"aid" binding:"required"`
}

type CancelCollectionRequest struct {
	AID string `json:"aid" binding:"required"`
}

// 生成方案请求
type CreateSolutionRequest struct {
	AID string `json:"aid" binding:"required"`
}

// 审美数据列表查询请求
type AestheticDataListRequest struct {
	Page     int    `form:"page" binding:"min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"min=1,max=100"` // 每页条数
	Name     string `form:"name"`                              // 姓名
	Gender   string `form:"gender"`                            // 性别
	AgeMin   int    `form:"age_min"`                           // 最小年龄
	AgeMax   int    `form:"age_max"`                           // 最大年龄
	Province string `form:"province"`                          // 省份
	City     string `form:"city"`                              // 所在城市
	Phone    string `form:"phone"`                             // 手机号码
	Source   *int   `form:"source" binding:"omitempty,oneof=0 1"`
}

// 审美数据统计分析请求
type AestheticAnalysisRequest struct {
	AnalysisType string `form:"analysis_type" binding:"required"` // 分析类型: color, disliked_color, adjective, image, region
	Dimension    string `form:"dimension"`                        // 分析维度: count, top, percent, map
	Top          int    `form:"top" binding:"max=100"`            // 取前N条数据
	Gender       string `form:"gender"`                           // 按性别过滤
	AgeMin       int    `form:"age_min"`                          // 最小年龄
	AgeMax       int    `form:"age_max"`                          // 最大年龄
	City         string `form:"city"`                             // 按城市过滤
	Province     string `form:"province"`                         // 省份
	Source       *int   `form:"source" binding:"omitempty,oneof=0 1"`
}

// 审美数据分析结果项
type AnalysisItem struct {
	Name    string  `json:"name"`    // 分析项名称
	Count   int     `json:"count"`   // 数量
	Percent float64 `json:"percent"` // 百分比
}

// ======== 管理员相关 ========

// 管理员登录请求
type AdminLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`    // 手机号码
	Password string `json:"password" binding:"required"` // 密码
}

// 管理员登录响应
type AdminLoginResponse struct {
	Token     string `json:"token"`      // 管理员token
	ExpiresIn int64  `json:"expires_in"` // token有效期（秒）
}

// 用户列表查询请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"min=1,max=100"` // 每页条数
	Phone    string `form:"phone"`                             // 手机号码
	Status   *int   `form:"status"`                            // 状态
}

// 用户状态更新请求
type UserStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"` // 状态 1:正常 0:禁用
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Name   string `json:"name"`   // 姓名
	Gender string `json:"gender"` // 性别
	Age    int    `json:"age"`    // 年龄
	City   string `json:"city"`   // 城市
}

// OpenProRequest 管理端开通Pro请求
type OpenProRequest struct {
	Days int `json:"days" binding:"required,oneof=30 90 180 365"` // 开通时长（天）
}

// BatchOpenProRequest 管理端批量开通Pro请求
type BatchOpenProRequest struct {
	Phones []string `json:"phones" binding:"required"`                   // 手机号列表，最多200个
	Days   int      `json:"days" binding:"required,oneof=30 90 180 365"` // 开通时长（天）
}

// BatchOpenProResponse 管理端批量开通Pro响应
type BatchOpenProResponse struct {
	TotalInput     int      `json:"total_input"`      // 输入手机号数量
	UniquePhones   int      `json:"unique_phones"`    // 去重后手机号数量
	SuccessCount   int      `json:"success_count"`    // 成功开通数量
	UpdatedPhones  []string `json:"updated_phones"`   // 成功开通手机号
	NotFoundPhones []string `json:"not_found_phones"` // 未找到的手机号
}

// DashboardOverviewRequest 管理后台看板查询请求
type DashboardOverviewRequest struct {
	Days int `form:"days" binding:"omitempty,oneof=7 15 30"` // 日期范围（天）
}

// DashboardMetricSummary 看板指标汇总
type DashboardMetricSummary struct {
	Total   int `json:"total"`   // 总量
	Plan    int `json:"plan"`    // 来源=0，璞览
	Moxiang int `json:"moxiang"` // 来源=1，茉香
}

// DashboardTrendItem 看板按天趋势项
type DashboardTrendItem struct {
	Date    string `json:"date"`    // 日期，格式 YYYY-MM-DD
	Total   int    `json:"total"`   // 总量
	Plan    int    `json:"plan"`    // 来源=0，璞览
	Moxiang int    `json:"moxiang"` // 来源=1，茉香
}

// DashboardMetricBlock 看板单个指标块
type DashboardMetricBlock struct {
	Summary DashboardMetricSummary `json:"summary"`
	Trend   []DashboardTrendItem   `json:"trend"`
}

// DashboardOverviewResponse 管理后台看板响应
type DashboardOverviewResponse struct {
	Days          int                  `json:"days"`
	StartDate     string               `json:"start_date"`
	EndDate       string               `json:"end_date"`
	Users         DashboardMetricBlock `json:"users"`
	AestheticData DashboardMetricBlock `json:"aesthetic_data"`
}

// PartnerAdminItem 管理后台合作方列表项
type PartnerAdminItem struct {
	Name                   string `json:"name"`                     // 合作方名称
	Source                 int    `json:"source"`                   // 合作方来源标识
	UserCount              int64  `json:"user_count"`               // 用户量
	AestheticDataCount     int64  `json:"aesthetic_data_count"`     // 审美数据量
	SolutionCount          int64  `json:"solution_count"`           // 生成方案量
	Status                 string `json:"status"`                   // 合作状态
	OpenPubkey             string `json:"open_pubkey"`              // 开放平台公钥
	OpenPrikey             string `json:"open_prikey"`              // 开放平台私钥
	SolutionRemainingCount string `json:"solution_remaining_count"` // 方案生成剩余次数
	SolutionExceedStrategy string `json:"solution_exceed_strategy"` // 方案超限策略
}

type CertQueryRequest struct {
	CertType   string `json:"cert_type"`
	Name       string `json:"name"`
	IDLastFour string `json:"id_last_four"`
}
