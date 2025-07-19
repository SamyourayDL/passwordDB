package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	DB         `yaml:"db"`
	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type DB struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Address  string `yaml:"address" env-default:"localhost:8082"`
	Name     string `yaml:"name" env-required:"true"`
	Key      []byte `env:"POSTGRES_KEY" env-default:"0123456789abcdef0123456789abcdef"`
}

const (
	configPathEnvVar = "CONFIG_PATH"
	dbPasswordEnvVar = "POSTGRES_PASSWORD"
	dbKeyEnvVar      = "POSTGRES_KEY"
)

func MustLoad() *Config { //die if something went wrong
	cfgPath := os.Getenv(configPathEnvVar)
	if cfgPath == "" {
		log.Fatal("failed to get config path")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatal("config file is not exist")
	}

	var config Config
	if err := cleanenv.ReadConfig(cfgPath, &config); err != nil {
		fmt.Println(err)
		log.Fatal("failed to read config")
	}

	return &config
}
