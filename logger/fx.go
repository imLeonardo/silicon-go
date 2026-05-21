package logger

import (
	"silicon-go/config"

	"go.uber.org/fx"
)

type ParamIn struct {
	fx.In

	Config *config.Config
}
type ParamOut struct {
	fx.Out

	Logger *Logger
}

func newFxModule(params ParamIn) (ParamOut, error) {
	logger, err := NewLogger(params.Config)
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{Logger: logger}, nil
}

var FxModule = fx.Module(
	"logger",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(logger *Logger) error {
				return logger.OnStart()
			}),
			fx.OnStop(func(logger *Logger) error {
				return logger.OnStop()
			}),
		),
	),
)

type Provider struct {
	registry.NopHook
}

func (p Provider) Register() any {
	return newFxModule
}
