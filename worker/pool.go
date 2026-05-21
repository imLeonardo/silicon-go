package worker

import (
	"context"
	"silicon-go/config"
	"silicon-go/logger"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Pool struct {
	registry.NopHook

	AppConfig *config.Config
	*logger.Logger
	Pool *ants.Pool

	ShutdownCtx        context.Context    // 用于通知所有goroutine退出的上下文
	ShutdownCancelFunc context.CancelFunc // 取消函数
}

func NewPool(config *config.Config, logger *logger.Logger) (*Pool, error) {
	return &Pool{
		AppConfig: config,
		Logger:    logger,
	}, nil
}

func (p *Pool) OnStart() error {
	var err error
	p.Pool, err = ants.NewPool(
		p.AppConfig.Worker.Pool.Size,
		ants.WithExpiryDuration(time.Duration(p.AppConfig.Worker.Pool.ExpiryDuration)*time.Millisecond),
		ants.WithPanicHandler(func(val interface{}) {
			// 记录 panic 到日志
			p.PanicWithFields("ants pool panic", p.LogField("panic", val))
		}),
	)
	if err != nil {
		return err
	}

	// 创建用于通知所有goroutine退出的上下文
	p.ShutdownCtx, p.ShutdownCancelFunc = context.WithCancel(context.Background())

	return nil
}

func (p *Pool) OnStop() error {
	p.ShutdownCancelFunc()

	return p.Pool.ReleaseTimeout(time.Duration(p.AppConfig.Worker.Pool.ExpiryDuration) * time.Millisecond)
}

func (p *Pool) Submit(task func(context.Context) error) error {
	ctx, cancel := context.WithCancel(p.ShutdownCtx)

	return p.Pool.Submit(func() {
		defer cancel()
		err := task(ctx)
		if err != nil {
			p.ErrorWithFields("worker pool task error", p.LogField("error", err))
			return
		}
	})
}
