package kredis

import (
	"context"
	"errors"
	"fmt"
	"git.4321.sh/feige/flygo/core/klog"
	"github.com/go-redis/redis/v8"
	"time"
)

type accessPlugin struct {
	addr   string
	config *config
	logger *klog.Component
}

func NewAccessLogPlugin(addr string, config *config, logger *klog.Component) *accessPlugin {
	p := &accessPlugin{
		addr:   addr,
		config: config,
		logger: logger,
	}
	return p
}

func (th *accessPlugin) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, ctxBegKey, time.Now()), nil
}

func (th *accessPlugin) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var fields = make([]klog.Field, 0, 15)
	var err = cmd.Err()
	cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

	fields = append(fields,
		klog.FieldClientKind(),
		klog.FieldDbSystem("redis"),
		klog.FieldDbOperation(cmd.Name()),
		klog.FieldDuration(cost),
		klog.FieldDBStatement(cmd.FullName()),
	)

	bodyStr := cmd.String()
	size := len(bodyStr)
	if size > th.config.MaxResContentSize {
		bodyStr = bodyStr[:th.config.MaxResContentSize] + " ..."
	}

	fields = append(fields, klog.FieldDBResult(bodyStr))

	fields = append(fields, klog.FieldNetPeerIp(th.addr))

	// redis 组件 debug 模式
	if th.config.Debug {
		th.config.HookReq = true
		th.config.HookRsp = true
	}

	key := ""

	if len(cmd.Args()) > 1 {
		key = fmt.Sprintf("%v", cmd.Args()[1])
	}

	// 慢执行
	if th.config.SlowLogThreshold > time.Duration(0) && cost > th.config.SlowLogThreshold {
		th.logger.WithCtx(ctx).Warn(fmt.Sprintf("[redis][slow][%s][%s] cost : %dms", cmd.Name(), key, cost.Milliseconds()), fields...)
	}

	// error metric
	if err != nil {
		fields = append(fields, klog.FieldErr(err))
		if errors.Is(err, redis.Nil) {
			th.logger.WithCtx(ctx).Warn(fmt.Sprintf("[redis][access][%s][%s] error : %s", cmd.Name(), key, err.Error()), fields...)
			return err
		}
		th.logger.WithCtx(ctx).Error(fmt.Sprintf("[redis][access][%s][%s] error : %s", cmd.Name(), key, err.Error()), fields...)
		return err
	}

	if th.config.HookLog {
		th.logger.WithCtx(ctx).Info(fmt.Sprintf("[redis][access][%s][%s]", cmd.Name(), key), fields...)
	}
	return err
}

func (th *accessPlugin) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (th *accessPlugin) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
