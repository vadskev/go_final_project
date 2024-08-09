package config

import (
	"flag"

	"github.com/joho/godotenv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "env/.env", "path to config file; example: -config-path .env")
	flag.Parse()
}

func Load() error {
	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

type HTTPConfig interface {
	Address() string
}

type LogConfig interface {
	Level() string
}

type DBConfig interface {
	Path() string
}

type PasswordConfig interface {
	GetPass() string
	CreateHash(key string) string
}
