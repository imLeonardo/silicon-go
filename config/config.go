package config

import (
	"fmt"
	"silicon-go/flag"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Config struct {
	registry.NopHook

	flag *flag.Flag

	App struct {
		MaxProcs  int      `mapstructure:"maxProcs"`  // 最大CPU数 0: 不设置
		AppModule []string `mapstructure:"appModule"` // 应用模块
	} `mapstructure:"app"`
	Worker struct {
		Pool struct {
			Size           int `mapstructure:"size"`
			ExpiryDuration int `mapstructure:"expiryDuration"`
		} `mapstructure:"pool"`
	} `mapstructure:"worker"`
	Logger struct {
		FilePattern  string `mapstructure:"filePattern"`
		MaxAge       int32  `mapstructure:"maxAge"`
		Level        int32  `mapstructure:"level"`
		SyncDuration int64  `mapstructure:"syncDuration"`
		IsDocker     bool   `mapstructure:"isDocker"`
		CallerSkip   int    `mapstructure:"callerSkip"`
	}
	Cluster struct {
		Type       string `mapstructure:"type"`
		Mode       int32  `mapstructure:"mode"`
		IsFrontend bool   `mapstructure:"isFrontend"`
		Port       int32  `mapstructure:"port"`
	} `mapstructure:"cluster"`
	Http struct {
		Addr  string       `mapstructure:"addr"`
		Fiber fiber.Config `mapstructure:"fiber"`
	} `mapstructure:"http"`
	Pitaya *pitayaConfig.PitayaConfig `mapstructure:"pitaya"`
}

func NewConfig(flag *flag.Flag) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(flag.ConfigFilePath)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{
		flag:   flag,
		Pitaya: pitayaConfig.NewDefaultPitayaConfig(),
	}

	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return c, nil
}

func (c *Config) OnStart() error {
	return nil
}

func (c *Config) Validate() error {
	if c.Cluster.IsFrontend {
		if c.Cluster.Port <= 0 || c.Cluster.Port > 65535 {
			return fmt.Errorf("invalid server port: %d", c.Cluster.Port)
		}
	}

	if c.Logger.Level < -1 || c.Logger.Level > 5 {
		return fmt.Errorf("invalid logger level: %d", c.Logger.Level)
	}

	if c.Worker.Pool.Size <= 0 {
		return fmt.Errorf("invalid worker pool size: %d", c.Worker.Pool.Size)
	}

	return nil
}
