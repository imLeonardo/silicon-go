package database

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

	Database *Database
}

func newFxModule(params ParamIn) (ParamOut, error) {
	database, err := NewDatabase(params.AppConfig, params.Logger)
	if err != nil {
		return ParamOut{}, err
	}
	return ParamOut{Database: database}, nil
}

var FxModule = fx.Module(
	"database",
	fx.Provide(
		fx.Annotate(
			newFxModule,
			fx.OnStart(func(database *Database) error {
				return database.OnStart()
			}),
			fx.OnStop(func(database *Database) error {
				return database.OnStop()
			}),
		),
	),
)
