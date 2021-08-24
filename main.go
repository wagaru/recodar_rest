package main

import (
	"fmt"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/delivery/daemon"
	"github.com/wagaru/recodar-rest/internal/delivery/http"
	"github.com/wagaru/recodar-rest/internal/logger"
	"github.com/wagaru/recodar-rest/internal/repository"
	"github.com/wagaru/recodar-rest/internal/usecase"
)

func main() {
	logger.Logger.Printf("Init...")

	// Load config
	config, err := config.LoadConfig("./", "app", "env")
	if err != nil {
		fmt.Printf("Load config failed: %v", err)
	}
	logger.Logger.Println("Load config completed.")

	// Init MongoDB
	repo, err := repository.NewMongoRepo(config)
	if err != nil {
		fmt.Printf("Init MongoDB failed: %v", err)
	}
	defer repo.Disconnect()
	logger.Logger.Println("Init Mongo DB completed.")

	// Init rabbitMQ
	rabbitMQRepo, err := repository.NewRabbitMQRepo(config)
	if err != nil {
		fmt.Printf("Init RabbitMQ failed: %v", err)
	}
	defer rabbitMQRepo.Disconnect()
	logger.Logger.Println("Init RabbitMQ completed.")

	// Init usecase
	_usecase := usecase.NewUsecase(repo, config)
	logger.Logger.Println("Init usecase completed.")

	// Init message broker usecase
	messageBrokerUsecase := usecase.NewMessageBrokerUsecase(rabbitMQRepo, config)
	logger.Logger.Println("Init message broker usecase completed.")

	// Init daemon
	daemon := daemon.NewDaemon(messageBrokerUsecase, config)
	logger.Logger.Println("Init daemon completed.")
	daemon.Run()

	// Init http
	delivery := http.NewHttpDelivery(_usecase, messageBrokerUsecase, config)
	logger.Logger.Println("Init delivery completed.")
	delivery.Run(config.ServerPort)
}
