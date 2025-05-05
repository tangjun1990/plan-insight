package kredis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.4321.sh/feige/flygo/core/klog"
	"git.4321.sh/feige/flygo/core/utils/xdebug"

	"github.com/go-redis/redis/v8"
)

const ctxBegKey = "_cmdResBegin_"

type hook struct {
	beforeProcess         func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess          func(ctx context.Context, cmd redis.Cmder) error
	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func (i *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return i.beforeProcess(ctx, cmd)
}

func (i *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return i.afterProcess(ctx, cmd)
}

func (i *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return i.beforeProcessPipeline(ctx, cmds)
}

func (i *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return i.afterProcessPipeline(ctx, cmds)
}

func newHook(compName string, config *config, logger *klog.Component) *hook {
	return &hook{
		beforeProcess: func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcess: func(ctx context.Context, cmd redis.Cmder) error {
			return nil
		},
		beforeProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) error {
			return nil
		},
	}
}

func (i *hook) setBeforeProcess(p func(ctx context.Context, cmd redis.Cmder) (context.Context, error)) *hook {
	i.beforeProcess = p
	return i
}

func (i *hook) setAfterProcess(p func(ctx context.Context, cmd redis.Cmder) error) *hook {
	i.afterProcess = p
	return i
}

func (i *hook) setBeforeProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)) *hook {
	i.beforeProcessPipeline = p
	return i
}

func (i *hook) setAfterProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) error) *hook {
	i.afterProcessPipeline = p
	return i
}

func fixedHook(compName string, config *config, logger *klog.Component) *hook {
	return newHook(compName, config, logger).
		setBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, ctxBegKey, time.Now()), nil
		}).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			var err = cmd.Err()
			// go-redis script的error做了prefix处理
			// https://github.com/go-redis/redis/blob/master/script.go#L61
			if err != nil && !strings.HasPrefix(err.Error(), "NOSCRIPT ") {
				err = fmt.Errorf("kredis exec command %s fail, %w", cmd.Name(), err)
			}
			return err
		})
}

// debug 模式下 命令行输出
func debugHook(compName string, config *config, logger *klog.Component) *hook {
	return newHook(compName, config, logger).setAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			if err != nil {
				log.Println("[kredis.response]",
					xdebug.MakeReqResError(compName, fmt.Sprintf("%v", config.Addrs), cost, fmt.Sprintf("%v", cmd.Args()), err.Error()),
				)
			} else {
				log.Println("[kredis.response]",
					xdebug.MakeReqResInfo(compName, fmt.Sprintf("%v", config.Addrs), cost, fmt.Sprintf("%v", cmd.Args()), response(cmd)),
				)
			}
			return err
		},
	)
}

// 处理解析 redis返回值
func response(cmd redis.Cmder) string {
	switch cmd.(type) {
	case *redis.Cmd:
		return fmt.Sprintf("%v", cmd.(*redis.Cmd).Val())
	case *redis.StringCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StringCmd).Val())
	case *redis.StatusCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StatusCmd).Val())
	case *redis.IntCmd:
		return fmt.Sprintf("%v", cmd.(*redis.IntCmd).Val())
	case *redis.DurationCmd:
		return fmt.Sprintf("%v", cmd.(*redis.DurationCmd).Val())
	case *redis.BoolCmd:
		return fmt.Sprintf("%v", cmd.(*redis.BoolCmd).Val())
	case *redis.CommandsInfoCmd:
		return fmt.Sprintf("%v", cmd.(*redis.CommandsInfoCmd).Val())
	case *redis.StringSliceCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StringSliceCmd).Val())
	default:
		return ""
	}
}
