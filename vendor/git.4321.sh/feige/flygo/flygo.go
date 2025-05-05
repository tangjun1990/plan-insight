package flygo

import (
	"context"
	"errors"
	"fmt"
	"git.4321.sh/feige/flygo/component/server"
	"git.4321.sh/feige/flygo/component/task/kcron"
	"git.4321.sh/feige/flygo/component/task/kjob"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.4321.sh/feige/flygo/core/kapp"
	"git.4321.sh/feige/flygo/core/kcfg"
	_ "git.4321.sh/feige/flygo/core/kcfg/file"
	"git.4321.sh/feige/flygo/core/kcfg/manager"
	"git.4321.sh/feige/flygo/core/kflag"
	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/utils/xcolor"
	"git.4321.sh/feige/flygo/core/utils/xcycle"
	"git.4321.sh/feige/flygo/core/utils/xtime"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM}

type Flygo struct {
	cycle    *xcycle.Cycle
	smu      *sync.RWMutex
	logger   *klog.Component
	err      error
	ctx      context.Context
	cancel   func()
	inits    []func() error
	invokers []func() error
	servers  []server.Server
	crons    []kcron.Kcron
	jobs     map[string]kjob.Kjob
	opts     opts
}

type opts struct {
	ctx               context.Context
	configPrefix      string
	hang              bool
	disableBanner     bool
	disableFlagConfig bool
	beforeStopClean   []func() error
	afterStopClean    []func() error
	stopTimeout       time.Duration
	shutdownSignals   []os.Signal
}

func New(options ...Option) *Flygo {
	e := &Flygo{
		cycle:    xcycle.NewCycle(),
		smu:      &sync.RWMutex{},
		logger:   klog.FlygoLogger,
		err:      nil,
		inits:    make([]func() error, 0),
		invokers: make([]func() error, 0),
		servers:  make([]server.Server, 0),
		crons:    make([]kcron.Kcron, 0),
		jobs:     make(map[string]kjob.Kjob),
		opts: opts{
			ctx:             context.Background(),
			hang:            false,
			configPrefix:    "",
			beforeStopClean: make([]func() error, 0),
			afterStopClean:  make([]func() error, 0),
			stopTimeout:     xtime.Duration("5s"),
			shutdownSignals: shutdownSignals,
		},
	}

	ctx, cancel := context.WithCancel(e.opts.ctx)
	e.ctx = ctx
	e.cancel = cancel

	e.inits = []func() error{
		e.parseFlags,
		loadConfig,
		loadGlobalConfig,
		e.setupLogger,
		setupMaxProcs,
	}

	e.err = runSerialFuncReturnError(e.inits)

	// 重新加载 logger
	e.logger = klog.FlygoLogger
	options = append(options, WithAfterStopClean(klog.DefaultLogger.Flush, klog.FlygoLogger.Flush))
	for _, option := range options {
		option(e)
	}

	return e
}

func (o *Flygo) Invoker(fns ...func() error) *Flygo {
	o.smu.Lock()
	defer o.smu.Unlock()

	o.invokers = append(o.invokers, fns...)

	o.err = runSerialFuncReturnError(o.invokers)
	return o
}

func (o *Flygo) Serve(s ...server.Server) *Flygo {
	o.smu.Lock()
	defer o.smu.Unlock()
	o.servers = append(o.servers, s...)
	return o
}

func (o *Flygo) Cron(w ...kcron.Kcron) *Flygo {
	o.crons = append(o.crons, w...)
	return o
}

func (o *Flygo) Job(runners ...kjob.Kjob) *Flygo {
	jobFlag := kflag.String("job")
	if jobFlag == "" {
		o.logger.Info("flag jobs name empty", klog.FieldComponent(kjob.PackageName))
		return o
	}

	jobMap := make(map[string]struct{})
	if strings.Contains(jobFlag, ",") {
		jobArr := strings.Split(jobFlag, ",")
		for _, value := range jobArr {
			jobMap[value] = struct{}{}
		}
	} else {
		jobMap[jobFlag] = struct{}{}
	}

	for _, runner := range runners {
		jobName := runner.ConfigKey()
		if jobName == "" {
			o.logger.Error("runner job name empty", klog.FieldComponent(runner.PackageName()))
			return o
		}
		if kflag.Bool("disable-job") {
			o.logger.Info("runner disable job", klog.FieldComponent(runner.PackageName()), klog.FieldName(jobName))
			return o
		}

		_, flag := jobMap[jobName]
		if flag {
			o.logger.Info("init register job", klog.FieldComponent(runner.PackageName()), klog.FieldName(jobName))
			o.jobs[jobName] = runner
		}
	}
	return o
}

func (o *Flygo) Run() error {
	if o.err != nil {
		runSerialFuncLogError(o.opts.afterStopClean)
		return o.err
	}

	if len(o.jobs) > 0 {
		return o.startJobs()
	}

	o.handleSignals()

	_ = o.startServers(o.ctx)

	_ = o.startCrons()

	if err := <-o.cycle.Wait(o.opts.hang); err != nil {
		o.logger.Error("Flygo shutdown with error", klog.FieldComponent("app"), klog.FieldErr(err))
		runSerialFuncLogError(o.opts.afterStopClean)
		return err
	}
	o.logger.Info("stop Flygo.", klog.FieldComponent("app"))
	runSerialFuncLogError(o.opts.afterStopClean)
	return nil
}

func (o *Flygo) Stop(ctx context.Context, isGraceful bool) (err error) {
	runSerialFuncLogError(o.opts.beforeStopClean)

	o.smu.RLock()
	if isGraceful {
		for _, s := range o.servers {
			func(s server.Server) {
				o.cycle.Run(func() error {
					return s.GraceShutdown(ctx)
				})
			}(s)
		}
	} else {
		for _, s := range o.servers {
			func(s server.Server) {
				o.cycle.Run(s.Stop)
			}(s)
		}
	}
	o.smu.RLock()

	for _, w := range o.crons {
		func(w kcron.Kcron) {
			o.cycle.Run(w.Stop)
		}(w)
	}
	<-o.cycle.Done()
	o.cancel()
	o.cycle.Close()
	return err
}

func (o *Flygo) handleSignals() {
	sig := make(chan os.Signal, 2)
	signal.Notify(
		sig,
		o.opts.shutdownSignals...,
	)

	go func() {
		s := <-sig
		grace := s != syscall.SIGQUIT
		go func() {
			stopCtx, cancel := context.WithTimeout(context.Background(), o.opts.stopTimeout)
			defer func() {
				signal.Stop(sig)
				cancel()
			}()
			_ = o.Stop(stopCtx, grace)
			<-stopCtx.Done()
			if errors.Is(stopCtx.Err(), context.DeadlineExceeded) {
				klog.Error("handleSignals stop context err", klog.FieldErr(stopCtx.Err()))
			}
		}()
		<-sig
		klog.Error("handleSignals quit")
		os.Exit(128 + int(s.(syscall.Signal)))
	}()
}

// 启动服务
func (o *Flygo) startServers(ctx context.Context) error {
	if len(o.servers) == 0 {
		return nil
	}

	err := runSerialFuncReturnError([]func() error{
		o.showBanner,
	})
	if err != nil {
		return err
	}

	// 启动服务
	for _, s := range o.servers {
		s := s
		o.cycle.Run(func() (err error) {
			o.logger.Info("start server", klog.FieldComponent(s.PackageName()), klog.FieldComponentName(s.ConfigKey()), klog.Any("addr", s.Info().Label()))
			defer o.logger.WithCtx(ctx).Info("stop server", klog.FieldComponent(s.PackageName()), klog.FieldComponentName(s.ConfigKey()), klog.FieldErr(err), klog.Any("addr", s.Info().Label()))
			err = s.Start()
			return
		})
	}
	return nil
}

func (o *Flygo) startCrons() error {
	for _, w := range o.crons {
		w := w
		o.cycle.Run(func() error {
			return w.Start()
		})
	}
	return nil
}

func (o *Flygo) startJobs() error {
	if len(o.jobs) == 0 {
		return nil
	}
	var jobs = make([]func() error, 0)
	for _, runner := range o.jobs {
		runner := runner
		jobs = append(jobs, func() error {
			return runner.Start()
		})
	}

	eg := errgroup.Group{}
	for _, fn := range jobs {
		eg.Go(fn)
	}
	return eg.Wait()
}

func (o *Flygo) parseFlags() error {
	if !o.opts.disableFlagConfig {
		kflag.Register(&kflag.StringFlag{
			Name:    "config",
			Usage:   "--config",
			EnvVar:  kapp.AppConfigPath,
			Default: "./config.toml",
			Action:  func(name string, fs *kflag.FlagSet) {},
		})
	}

	kflag.Register(&kflag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: true,
		EnvVar:  "CONFIG_WATCH",
	})

	kflag.Register(&kflag.BoolFlag{
		Name:    "version",
		Usage:   "--version, print version",
		Default: false,
		Action: func(string, *kflag.FlagSet) {
			kapp.PrintVersion()
			os.Exit(0)
		},
	})

	kflag.Register(&kflag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		EnvVar:  kapp.EnvAppHost,
		Default: "0.0.0.0",
		Action:  func(string, *kflag.FlagSet) {},
	})
	return kflag.Parse()
}

// 加载配置
func loadConfig() error {
	var configAddr = kflag.String("config")
	provider, parser, tag, err := manager.NewDataSource(configAddr, kflag.Bool("watch"))

	if err == manager.ErrDefaultConfigNotExist {
		panic(fmt.Sprintf("no config, addr: %s, err: %s", configAddr, err.Error()))
	}

	if err != nil {
		panic("data source: provider error: " + err.Error())
	}

	if err := kcfg.LoadFromDataSource(provider, parser, kcfg.WithTagName(tag)); err != nil {
		panic("data source: load config - unmarshal config err: " + err.Error())
	}
	return nil
}

// 设置全局配置
func loadGlobalConfig() error {
	// 加载配置
	if err := kcfg.UnmarshalKey(kapp.GlobalKey, &kapp.GlobalConfig); err != nil {
		panic(fmt.Sprintf("[flygo.loadGlobalConfig] %s is required, error: %s", kapp.GlobalKey, err.Error()))
	}
	// 验证配置
	if err := kapp.GlobalConfig.CheckAndInit(); err != nil {
		panic("[flygo.loadGlobalConfig] error:" + err.Error())
	}

	return nil
}

// 设置日志组件
func (o *Flygo) setupLogger() error {
	// 系统日志
	if kcfg.Get(o.opts.configPrefix+"logger.flygo") != nil {
		klog.FlygoLogger = klog.Load(o.opts.configPrefix + "logger.flygo").Build()
		klog.FlygoLogger.Info("setup flygo logger", klog.FieldComponent(klog.PackageName))
	}

	// 应用日志
	if kcfg.Get(o.opts.configPrefix+"logger.default") != nil {
		klog.DefaultLogger = klog.Load(o.opts.configPrefix + "logger.default").Build()
		klog.FlygoLogger.Info("setup default logger", klog.FieldComponent(klog.PackageName))
	}
	return nil
}

func setupMaxProcs() error {
	if maxProcs := kcfg.GetInt("flygo.maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			klog.FlygoLogger.Panic("setup max procs", klog.FieldComponent("app"), klog.FieldErr(err))
		}
	}
	klog.FlygoLogger.Info("setup max procs", klog.FieldComponent("app"), klog.FieldValueAny(runtime.GOMAXPROCS(-1)))
	return nil
}

func (o *Flygo) showBanner() error {
	if o.opts.disableBanner {
		return nil
	}
	const banner = `
   ▄████████  ▄█       ▄██   ▄           ▄██████▄   ▄██████▄  
  ███    ███ ███       ███   ██▄        ███    ███ ███    ███ 
  ███    █▀  ███       ███▄▄▄███        ███    █▀  ███    ███ 
 ▄███▄▄▄     ███       ▀▀▀▀▀▀███       ▄███        ███    ███ 
▀▀███▀▀▀     ███       ▄██   ███      ▀▀███ ████▄  ███    ███ 
  ███        ███       ███   ███        ███    ███ ███    ███ 
  ███        ███▌    ▄ ███   ███        ███    ███ ███    ███ 
  ███        █████▄▄██  ▀█████▀         ████████▀   ▀██████▀  
             ▀                                                

flygo-v1.0.0启动中...
使用中任何问题可联系：tangjun@feige.cn
`
	fmt.Println(xcolor.Blue(banner))
	return nil
}

func runSerialFuncReturnError(fns []func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func runSerialFuncLogError(fns []func() error) {
	for _, clean := range fns {
		err := clean()
		if err != nil {
			klog.FlygoLogger.Error("beforeStopClean err", klog.FieldComponent("app"), klog.FieldErr(err))
		}
	}
}

type Option func(a *Flygo)

func WithHang(flag bool) Option {
	return func(a *Flygo) {
		a.opts.hang = flag
	}
}

func WithDisableBanner(disableBanner bool) Option {
	return func(a *Flygo) {
		a.opts.disableBanner = disableBanner
	}
}

func WithDisableFlagConfig(disableFlagConfig bool) Option {
	return func(a *Flygo) {
		a.opts.disableFlagConfig = disableFlagConfig
	}
}

func WithConfigPrefix(configPrefix string) Option {
	return func(a *Flygo) {
		a.opts.configPrefix = configPrefix
	}
}

func WithBeforeStopClean(fns ...func() error) Option {
	return func(a *Flygo) {
		a.opts.beforeStopClean = append(a.opts.beforeStopClean, fns...)
	}
}

func WithAfterStopClean(fns ...func() error) Option {
	return func(a *Flygo) {
		a.opts.afterStopClean = append(a.opts.afterStopClean, fns...)
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(e *Flygo) {
		e.opts.stopTimeout = timeout
	}
}

func WithShutdownSignal(signals ...os.Signal) Option {
	return func(e *Flygo) {
		e.opts.shutdownSignals = append(e.opts.shutdownSignals, signals...)
	}
}
