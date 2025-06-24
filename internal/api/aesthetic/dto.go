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
	LikedImages     []string `json:"liked_images" binding:"required,min=1,max=10"`     // 喜欢的图片
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
