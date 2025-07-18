package server

import (
	"context"

	"github.com/tangjun1990/flygo/component/server/kin"
	"github.com/tangjun1990/plan-insight/internal/api/aesthetic"
	"github.com/tangjun1990/plan-insight/internal/server/router"
	"gorm.io/gorm"
)

// Server 服务结构体
type Server struct {
	server *kin.Component
	db     *gorm.DB
}

// NewServer 创建服务实例
func NewServer(server *kin.Component, db *gorm.DB) *Server {
	return &Server{
		server: server,
		db:     db,
	}
}

// Start 启动服务
func (s *Server) Start(ctx context.Context) error {
	// 初始化数据表结构
	if err := s.initDatabase(); err != nil {
		return err
	}

	// 注册路由
	router.AddRouter(s.server, s.db)

	return nil
}

// initDatabase 初始化数据库
func (s *Server) initDatabase() error {
	// 迁移审美感知应用相关的表结构
	if err := aesthetic.AutoMigrate(s.db); err != nil {
		return err
	}

	return nil
}
