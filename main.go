package main

import (
	"fmt"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/delivery/http"
	logger "github.com/wagaru/recodar-rest/internal/logger"
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

	// Init usecase
	usecase := usecase.NewUsecase(repo, config)
	logger.Logger.Println("Init usecase completed.")

	// Init delivery
	delivery := http.NewHttpDelivery(usecase, config)
	logger.Logger.Println("Init delivery completed.")

	delivery.Run(config.ServerPort)
}
