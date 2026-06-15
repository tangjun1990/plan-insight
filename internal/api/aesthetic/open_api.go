package aesthetic

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tangjun1990/flygo/core/kcfg"
	"gorm.io/gorm"
)

const (
	openAPIAccessTokenTTL = 100 * time.Minute // 测试模式下100分钟过期，生产环境10分钟
	openUserTokenTTL      = 30 * 24 * time.Hour
)

func buildOpenAccessSign(pubkey, prikey string, t time.Time) string {
	fmt.Println(strconv.FormatInt(t.Unix(), 10))
	return hashPasswordMD5(strings.TrimSpace(pubkey) + "|" + strings.TrimSpace(prikey) + "|" + strconv.FormatInt(t.Unix(), 10))
}

func BuildOpenAccessToken(pubkey, prikey string, t time.Time) string {
	pubkey = strings.TrimSpace(pubkey)
	timestamp := strconv.FormatInt(t.Unix(), 10)
	sign := buildOpenAccessSign(pubkey, prikey, t)
	raw := pubkey + "." + sign + "." + timestamp
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func BuildOpenAuthorizationHeader(pubkey, prikey string, t time.Time) string {
	return "Bearer " + BuildOpenAccessToken(pubkey, prikey, t)
}

func getOpenAPIKeys() (string, string) {
	return strings.TrimSpace(kcfg.GetString("app.global.open_pubkey")),
		strings.TrimSpace(kcfg.GetString("app.global.open_prikey"))
}

func isValidPhoneNumber(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func generateOpenUserToken(phone string) string {
	return "open_" + hashPasswordMD5(fmt.Sprintf("%s|%d|%d", phone, time.Now().UnixNano(), time.Now().Unix()))
}

func (s *Service) ValidateOpenAccessToken(token string) error {
	pubkey, prikey := getOpenAPIKeys()
	if pubkey == "" || prikey == "" {
		return errors.New("开放平台密钥未配置")
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return errors.New("开放平台token格式错误")
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		return errors.New("开放平台token内容错误")
	}

	tokenPubKey := strings.TrimSpace(parts[0])
	sign := strings.TrimSpace(parts[1])
	timestampStr := strings.TrimSpace(parts[2])

	if tokenPubKey != pubkey {
		return errors.New("开放平台公钥无效")
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return errors.New("开放平台时间戳错误")
	}

	expectedSign := buildOpenAccessSign(tokenPubKey, prikey, time.Unix(timestamp, 0))
	if sign != expectedSign {
		return errors.New("开放平台签名无效")
	}

	tokenTime := time.Unix(timestamp, 0)
	now := time.Now()
	if tokenTime.Before(now.Add(-openAPIAccessTokenTTL)) || tokenTime.After(now.Add(openAPIAccessTokenTTL)) {
		return errors.New("开放平台token已过期")
	}

	return nil
}

func (s *Service) IssueOpenUserToken(phone string, wxOpenID string) (*OpenUserAuthTokenResponse, error) {
	phone = strings.TrimSpace(phone)
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}
	wxOpenID = strings.TrimSpace(wxOpenID)
	if wxOpenID == "" {
		return nil, errors.New("微信小程序用户ID格式错误")
	}

	now := time.Now()
	expireAt := now.Add(openUserTokenTTL)

	var user User
	result := s.db.Where("phone = ?", phone).First(&user)
	switch {
	case result.Error == nil:
		if user.Status == 0 {
			return nil, errors.New("用户已被禁用")
		}
		if user.OpenToken == "" || user.OpenTokenExpireAt == nil || user.OpenTokenExpireAt.Before(now) || !strings.HasPrefix(user.OpenToken, "open_") {
			user.OpenToken = generateOpenUserToken(phone)
			user.OpenTokenExpireAt = &expireAt
		} else {
			expireAt = *user.OpenTokenExpireAt
		}
		user.LastLoginTime = &now
		if err := s.db.Save(&user).Error; err != nil {
			return nil, err
		}
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		user = User{
			Name:              "PLAN用户" + getRandomString(6),
			Phone:             phone,
			Status:            1,
			OpenToken:         generateOpenUserToken(phone),
			OpenTokenExpireAt: &expireAt,
			Token:             "",
			FirstLoginTime:    &now,
			LastLoginTime:     &now,
			WxOpenID:          wxOpenID,
			Source:            1, // 用户注册来源是华筑会
		}
		if err := s.db.Create(&user).Error; err != nil {
			return nil, err
		}
	default:
		return nil, result.Error
	}

	return &OpenUserAuthTokenResponse{
		UserToken: OpenUserTokenItem{
			Token:    user.OpenToken,
			ExpireAt: expireAt.Unix(),
		},
	}, nil
}

func (s *Service) GetOpenAestheticDataList(req *OpenAestheticListRequest) (*PageResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 20 {
		return nil, errors.New("page_size最大不超过20")
	}

	var user User
	userID := uint(0)
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err == nil {
		if user.Status == 0 {
			return nil, errors.New("用户已被禁用")
		}
		userID = user.ID
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var total int64
	if err := s.db.Model(&AestheticData{}).Where("phone = ?", phone).Count(&total).Error; err != nil {
		return nil, err
	}

	var list []AestheticData
	if err := s.db.Where("phone = ?", phone).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := make([]*AestheticDataRsp, 0, len(list))
	for _, item := range list {
		detail, err := s.GetAestheticDataDetail(item.ID, userID, false)
		if err != nil {
			return nil, err
		}
		result = append(result, detail)
	}

	return &PageResponse{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) GetOpenUserSolutionList(req *OpenSolutionListRequest) (*PageResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 20 {
		return nil, errors.New("page_size最大不超过20")
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

	return s.GetUserSolutionList(user.ID, page, pageSize)
}

func (s *Service) OpenCreateUserSolution(req *OpenCreateSolutionRequest) (*UserSolution, error) {
	phone := strings.TrimSpace(req.Phone)
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	aid := req.AID
	if aid <= 0 {
		return nil, errors.New("无效的报告ID参数")
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

	var report AestheticData
	if err := s.db.Where("id = ? AND user_id = ?", aid, user.ID).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("报告不存在或不属于当前用户")
		}
		return nil, err
	}

	return s.CreateUserSolution(user.ID, uint(aid))
}

func (s *Service) GetOpenUserSolutionDetail(req *OpenSolutionDetailRequest) (*AestheticDataRsp, error) {
	phone := strings.TrimSpace(req.Phone)
	if !isValidPhoneNumber(phone) {
		return nil, errors.New("手机号格式错误")
	}

	if req.SolutionID <= 0 {
		return nil, errors.New("无效的solution_id参数")
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

	return s.GetUserSolutionDetail(user.ID, uint(req.SolutionID))
}
