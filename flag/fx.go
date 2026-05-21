package flag

import (
	"go.uber.org/fx"
)

type ParamIn struct {
	fx.In
}

type ParamOut struct {
	fx.Out

	Flag *Flag
}

func newFxModule(params ParamIn) (ParamOut, error) {
	flag, err := NewFlag()
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{Flag: flag}, nil
}

var FxModule = fx.Module(
	"flag",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(flag *Flag) error {
				return flag.OnStart()
			}),
			fx.OnStop(func(flag *Flag) error {
				return flag.OnStop()
			}),
		),
	),
)
