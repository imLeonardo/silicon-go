package config

import (
	"silicon-go/flag"

	"go.uber.org/fx"
)

type ParamIn struct {
	fx.In

	Flag *flag.Flag
}

type ParamOut struct {
	fx.Out

	AppConfig *Config
}

func newFxModule(params ParamIn) (ParamOut, error) {
	config, err := NewConfig(params.Flag)
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{AppConfig: config}, nil
}

var FxModule = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(appConfig *Config) error {
				return appConfig.OnStart()
			}),
			fx.OnStop(func(appConfig *Config) error {
				return appConfig.OnStop()
			}),
		),
	),
)
