package logger

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("cannot load env from %s", path)
	}

	return &Config{
		Env: os.Getenv("ENV"),
	}, nil

}

var Log *slog.Logger

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
