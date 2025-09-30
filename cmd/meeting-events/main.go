package main

import (
	
	"github.com/joho/godotenv"

	"github.com/hihikaAAa/meeting-events/internal/config"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
)

func main(){
	_ = godotenv.Load("local.env")

	cfg := config.MustLoad()

	log := setup.SetupLogger(cfg.Env)

	log.Info("Запустил проект")
}