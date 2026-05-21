package http

import (
	"context"
	"silicon-go/config"
	"silicon-go/logger"
	"silicon-go/worker"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Ctx struct {
	fiber.Ctx
}

type Server struct {
	registry.NopHook

	AppConfig *config.Config
	*logger.Logger
	Ticker *time.Ticker

	Pool *worker.Pool
	*fiber.App
}

func NewServer(appConfig *config.Config, logger *logger.Logger, pool *worker.Pool) (*Server, error) {
	return &Server{
		AppConfig: appConfig,
		Logger:    logger,
		Pool:      pool,
		App:       fiber.New(appConfig.Http.Fiber),
	}, nil
}

func (s *Server) OnStart() error {
	// 注册路由
	s.initRouter()

	// 启动服务器
	err := s.Pool.Submit(func(ctx context.Context) error {
		err := s.Listen(s.AppConfig.Http.Addr)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) OnStop() error {
	err := s.Shutdown()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) initRouter() {
	opGroup := s.Group("/op")
	opGroup.Post("/ready", requestHandlerReady)
	opGroup.Post("/health", requestHandlerHealth)

	apiGroup := s.Group("/api")
	apiGroup.Post("/", requestHandlerApi)
}

func requestHandlerRoot(ctx *fiber.Ctx) error {
	return ctx.SendString("root")
}

func requestHandlerReady(ctx *fiber.Ctx) error {
	return ctx.SendString("ready")
}

func requestHandlerHealth(ctx *fiber.Ctx) error {
	return ctx.SendString("health")
}

func requestHandlerApi(ctx *fiber.Ctx) error {
	return ctx.SendString("api")
}
