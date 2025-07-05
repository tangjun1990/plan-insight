package kcron

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/kmetric"
)

type wrappedJob struct {
	NamedJob
	logger *klog.Component
}

func (wj wrappedJob) Run() {
	wj.run()
}

func (wj wrappedJob) run() {
	ctx := context.Background()

	kmetric.JobHandleCounter.Inc("cron", wj.Name(), "begin")
	var fields = []klog.Field{zap.String("name", wj.Name())}

	wj.logger.WithCtx(ctx).Info("cron start", fields...)
	var beg = time.Now()
	defer func() {
		var err error
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		fields = append(fields, klog.Duration("cost", time.Since(beg)))
		if err != nil {
			fields = append(fields, klog.FieldErr(err))
			wj.logger.WithCtx(ctx).Error("cron end", fields...)
		} else {
			wj.logger.WithCtx(ctx).Info("cron end", fields...)
		}
		kmetric.JobHandleHistogram.Observe(time.Since(beg).Seconds(), "cron", wj.Name())
	}()

	err := wj.NamedJob.Run(ctx)
	if err != nil {
		fields = append(fields, klog.FieldErr(err))
		wj.logger.WithCtx(ctx).Error("cron run failed", fields...)
	}
}
