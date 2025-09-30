package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Env string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	App App `yaml:"app"`
	DB  DB `yaml:"db"`
	Migrations Migrations `yaml:"migrations"`
	Outbox Outbox `yaml:"outbox"`
}

type App struct{
	Name string `yaml:"name" env-default:"meeting-svc" env-required:"true"`
	HTTP
}

type HTTP struct{
	Address string `yaml:"address" env-default:"localhost:8081"`
	User string `yaml:"user"`
	Password string `yaml:"password" env:"HTTP_SERVER_PASSWORD"`
	HTTPTimeout `yaml:"timeouts"`
}

type HTTPTimeout struct{
	ReadTimeout time.Duration `yaml:"read" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"idle" env-default:"60s"`
	EventTimeout time.Duration `yaml:"event" env-default:"5s"`
}

type DB struct{
	DSN string `yaml:"dsn" env-default:"postgres://user:pass@db:5432/meetings?sslmode=disable" env-required:"true"`
	MaxOpenConns int `yaml:"max_open_conns" env-default:"20"`
	MaxIdleConns int `yaml:"max_idle_conns" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"30m"`
}

type Migrations struct{
	Dir string `yaml:"dir" env-default:"file://migrations"`
}

type Outbox struct{
	PollInterval time.Duration `yaml:"poll_interval" env-default:"3s"`
	BatchSize int `yaml:"batch_size" env-default:"100"`
}

func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath,&cfg); err != nil{
		log.Fatalf("cannot read config file: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %v", err)
	}

	return &cfg
}
