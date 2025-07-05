package kin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/tangjun1990/flygo/component/server"
	"github.com/tangjun1990/flygo/core/kapp"

	"github.com/gin-gonic/gin"
	"github.com/tangjun1990/flygo/core/klog"
)

const PackageName = "server.kin"

type Component struct {
	mu     sync.Mutex
	name   string
	config *Config
	logger *klog.Component
	*gin.Engine
	Server           *http.Server
	listener         net.Listener
	routerCommentMap map[string]string
}

func newComponent(name string, config *Config, logger *klog.Component) *Component {
	gin.SetMode(config.Mode)
	return &Component{
		name:             name,
		config:           config,
		logger:           logger,
		Engine:           gin.New(),
		listener:         nil,
		routerCommentMap: make(map[string]string),
	}
}

func (c *Component) ConfigKey() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) RegisterRouteComment(method, path, comment string) {
	c.routerCommentMap[commentUniqKey(method, path)] = comment
}

func (c *Component) routerLog() {
	if c.config.Mode == gin.ReleaseMode {
		return
	}
	for _, route := range c.Engine.Routes() {
		info, flag := c.routerCommentMap[commentUniqKey(route.Method, route.Path)]
		if flag {
			c.logger.Info("add route", klog.FieldHttpMethod(route.Method), klog.String("path", route.Path), klog.Any("info", info))
		} else {
			c.logger.Info("add route", klog.FieldHttpMethod(route.Method), klog.String("path", route.Path))
		}
	}
}

func (c *Component) Start() error {
	c.routerLog()

	var err error
	c.mu.Lock()
	// windows不支持热更新，暂注释
	if c.config.Grace {
		/*s := endless.NewServer(c.config.Address(), c.Engine)
		err = s.ListenAndServe()
		c.mu.Unlock()
		if err != nil {
			log.Fatal("server err:" + err.Error())
		}*/
	} else {
		listener, err := net.Listen("tcp", c.config.Address())
		if err != nil {
			c.logger.Panic("new kin server err", klog.FieldErrKind("listen err"), klog.FieldErr(err))
		}
		c.config.Port = listener.Addr().(*net.TCPAddr).Port
		c.listener = listener

		c.Server = &http.Server{
			Addr:    c.config.Address(),
			Handler: c,
		}
		c.mu.Unlock()
		err = c.Server.Serve(c.listener)
		if err == http.ErrServerClosed {
			return nil
		}
	}

	return err
}

func (c *Component) Stop() error {
	c.mu.Lock()
	err := c.Server.Close()
	c.mu.Unlock()
	return err
}

func (c *Component) GraceShutdown(ctx context.Context) error {
	c.mu.Lock()
	err := c.Server.Shutdown(ctx)
	c.mu.Unlock()
	return err
}

func (c *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(c.config.Address()),
		server.WithKind(kapp.ServiceProvider),
	)
	return &info
}

func commentUniqKey(method, path string) string {
	return fmt.Sprintf("%s@%s", strings.ToLower(method), path)
}
