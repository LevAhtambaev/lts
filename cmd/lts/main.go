package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"log"
	"lts/internal/pkg/app"
	"os"

	"lts/internal/app/config"
)

var (
	//cfgPath = flag.String("config", "", "path to config file")
	cfgPath = "./configs/config.yaml"
)

// @title LTS (Leo`s Travel Stories)
// @version 1.0
// @description A collection of travels and visited places

// @contact.name API Support
// @contact.url https://vk.com/hopply_time
// @contact.email ahtambaev.lev@wb.ru

// @license.name AS IS (NO WARRANTY)

// @host localhost:8080
// @schemes  http https
// @BasePath /api

func main() {
	ctx := context.Background()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error at zap.NewProduction: %s", err.Error())
	}

	logger := zapLogger.Sugar()

	// считали конфиг
	flag.Parse()
	cfg, err := config.Read(ctx, cfgPath)
	if err != nil {
		log.Print("[config.Read]")

		os.Exit(2)
	}

	// Создание приложения
	application := app.New(ctx, cfg, logger)

	// Запуск приложения
	err = application.Run()
	if err != nil {
		log.Print("[application.Run]")

		os.Exit(2)
	}
}
