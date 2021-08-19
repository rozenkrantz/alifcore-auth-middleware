package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewConfig)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int64
	IsSet(key string) bool
}

type config struct {
	cfg *viper.Viper
}

func NewConfig() Config {

	cfg := viper.New()
	cfg.SetConfigName(".env")
	cfg.SetConfigType("env")
	cfg.AddConfigPath("./")
	cfg.AddConfigPath("../../")

	if err := cfg.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &config{cfg: cfg}
}

func (c *config) Get(key string) interface{} {
	return c.cfg.Get(key)
}

func (c *config) GetString(key string) string {
	return c.cfg.GetString(key)
}

func (c *config) IsSet(key string) bool {
	return c.cfg.IsSet(key)
}

func (c *config) GetInt(key string) int64 {
	return c.cfg.GetInt64(key)
}
