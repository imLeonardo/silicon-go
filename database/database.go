package database

import (
	"context"
	"silicon-go/config"
	"silicon-go/logger"
)

type Database struct {
	registry.NopHook

	AppConfig *config.Config
	*logger.Logger

	shutdownCtx        context.Context    // 用于通知所有goroutine退出的上下文
	shutdownCancelFunc context.CancelFunc // 取消函数
}

func NewDatabase(appConfig *config.Config, logger *logger.Logger) (*Database, error) {
	return &Database{
		AppConfig: appConfig,
		Logger:    logger,
	}, nil
}

func (d *Database) OnStart() error {
	// 创建用于通知所有goroutine退出的上下文
	d.shutdownCtx, d.shutdownCancelFunc = context.WithCancel(context.Background())

	return nil
}

func (d *Database) OnStop() error {
	d.shutdownCancelFunc()

	return nil
}
