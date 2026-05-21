package worker

import (
	"silicon-go/config"
	"silicon-go/logger"

	"go.uber.org/fx"
)

type ParamIn struct {
	fx.In

	AppConfig *config.Config
	Logger    *logger.Logger
}

type ParamOut struct {
	fx.Out

	Pool *Pool
}

func newFxModule(params ParamIn) (ParamOut, error) {
	pool, err := NewPool(params.AppConfig, params.Logger)
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{Pool: pool}, nil
}

var FxModule = fx.Module(
	"worker",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(pool *Pool) error {
				return pool.OnStart()
			}),
			fx.OnStop(func(pool *Pool) error {
				return pool.OnStop()
			}),
		),
	),
)
