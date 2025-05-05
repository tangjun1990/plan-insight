package main

import (
	"git.4321.sh/feige/commonapi/internal/server"
	"git.4321.sh/feige/commonapi/pkg/db"
	"git.4321.sh/feige/flygo"
	"git.4321.sh/feige/flygo/component/server/kin"
	"git.4321.sh/feige/flygo/core/klog"
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
