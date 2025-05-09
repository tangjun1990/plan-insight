package aesthetic

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Name           string         `gorm:"size:50" json:"name"`              // 用户姓名
	Phone          string         `gorm:"size:20;uniqueIndex" json:"phone"` // 手机号码
	Gender         string         `gorm:"size:10" json:"gender"`            // 性别
	Age            int            `json:"age"`                              // 年龄
	City           string         `gorm:"size:50" json:"city"`              // 所在城市
	Status         int            `gorm:"default:1" json:"status"`          // 状态 1:正常 0:禁用
	FirstLoginTime *time.Time     `json:"first_login_time"`                 // 首次登录时间
	LastLoginTime  *time.Time     `json:"last_login_time"`                  // 上次登录时间
	Token          string         `gorm:"size:255" json:"-"`                // 用户token
	TokenExpireAt  *time.Time     `json:"-"`                                // token过期时间
	WxOpenID       string         `gorm:"size:50;index" json:"-"`           // 微信小程序OpenID
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// AestheticData 审美数据模型
type AestheticData struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	UserID          uint           `gorm:"index" json:"user_id"`              // 用户ID
	User            User           `gorm:"foreignKey:UserID" json:"user"`     // 关联用户
	Name            string         `gorm:"size:50" json:"name"`               // 用户姓名
	Gender          string         `gorm:"size:10" json:"gender"`             // 性别
	Age             int            `json:"age"`                               // 年龄
	City            string         `gorm:"size:50" json:"city"`               // 所在城市
	Phone           string         `gorm:"size:20" json:"phone"`              // 手机号码
	LikedColors     string         `gorm:"type:text" json:"liked_colors"`     // 喜欢的10个颜色，JSON格式
	DislikedColors  string         `gorm:"type:text" json:"disliked_colors"`  // 讨厌的5个颜色，JSON格式
	LikedAdjectives string         `gorm:"type:text" json:"liked_adjectives"` // 喜欢的10个形容词，JSON格式
	LikedImages     string         `gorm:"type:text" json:"liked_images"`     // 喜欢的10张图片，JSON格式
	ColorImageURL   string         `gorm:"size:255" json:"color_image_url"`   // 色彩分析图片URL
	BoxImageURL     string         `gorm:"size:255" json:"box_image_url"`     // 坐标分析图片URL
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type AestheticDataRsp struct {
	AestheticData
	ColorExplorImage string   `json:"color_explor_image"`
	BoxExplorImage   string   `json:"box_explor_image"`
	Comment          []string `json:"comment"`
}

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Phone     string         `gorm:"size:20;uniqueIndex" json:"phone"` // 手机号码
	Password  string         `gorm:"size:255" json:"-"`                // 密码
	Token     string         `gorm:"size:255" json:"-"`                // 管理员token
	ExpireAt  *time.Time     `json:"-"`                                // token过期时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
