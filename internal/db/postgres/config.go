package postgres

import (
	"fmt"

	"github.com/Naumovets/go-search/internal/db"
	"github.com/spf13/viper"
)

func NewConfig(path string) (db.Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return db.Config{}, fmt.Errorf("cannot read config from %s: %w", path, err)
	}
	var cfg db.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return db.Config{}, fmt.Errorf("cannot unmarshal config: %w", err)
	}
	return cfg, nil
}
