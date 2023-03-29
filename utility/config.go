package utility

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig(path string) (*viper.Viper, error) {
	cfg := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	cfg.SetConfigFile(path)
	cfg.SetConfigType("ini")
	cfg.AutomaticEnv()
	if err := cfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading configuration: %v", err)
	}

	return cfg, nil
}
