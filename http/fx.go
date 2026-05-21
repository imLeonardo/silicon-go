package http

import (
	"silicon-go/config"
	"silicon-go/logger"
	"silicon-go/worker"

	"go.uber.org/fx"
)

type ParamIn struct {
	fx.In

	AppConfig *config.Config
	Logger    *logger.Logger
	Pool      *worker.Pool
}

type ParamOut struct {
	fx.Out

	HttpServer *Server
}

func newFxModule(params ParamIn) (ParamOut, error) {
	httpServer, err := NewServer(params.AppConfig, params.Logger, params.Pool)
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{HttpServer: httpServer}, nil
}

var FxModule = fx.Module(
	"http",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(http *Server) error {
				return http.OnStart()
			}),
			fx.OnStop(func(http *Server) error {
				return http.OnStop()
			}),
		),
	),
)
