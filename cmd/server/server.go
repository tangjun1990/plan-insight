package main

import (
	"github.com/tangjun1990/flygo"
	"github.com/tangjun1990/flygo/component/server/kin"
	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/plan-insight/internal/server"
	"github.com/tangjun1990/plan-insight/pkg/db"
)

func main() {
	if err := flygo.New().
		Invoker(db.InitBatch).
		Serve(func() *kin.Component {
			// 初始化 HTTP 服务
			httpServer := kin.Load("server.http").Build()

			// 获取数据库连接
			dbConn := db.GetDBMust("default")

			// 创建服务实例
			srv := server.NewServer(httpServer, dbConn)

			// 启动服务
			if err := srv.Start(nil); err != nil {
				klog.Panic("server start failed", klog.FieldErr(err))
			}

			// 返回HTTP服务组件
			return httpServer
		}()).Run(); err != nil {
		klog.Panic("startup", klog.FieldErr(err))
	}
}
