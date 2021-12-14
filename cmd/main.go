package main

import (
	"flag"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/server"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"log"
)

// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("(EventSourcingService)")
	appLogger.Fatal(server.NewServer(cfg, appLogger).Run())
}
