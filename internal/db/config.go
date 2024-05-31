package db

type Config struct {
	DB_HOST  string `mapstructure:"DB_HOST"`
	DB_PORT  string `mapstructure:"DB_PORT"`
	DB_NAME  string `mapstructure:"DB_NAME"`
	PASSWORD string `mapstructure:"POSTGRES_PASSWORD"`
	USER     string `mapstructure:"POSTGRES_USER"`
}
