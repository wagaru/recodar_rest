package main

import (
	"fmt"
	"log"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/delivery/http"
	"github.com/wagaru/recodar-rest/internal/repository"
	"github.com/wagaru/recodar-rest/internal/usecase"
)

func main() {
	log.Printf("Init...")

	// Load config
	config, err := config.LoadConfig("./", "app", "env")
	if err != nil {
		fmt.Printf("Load config failed: %v", err)
	}
	log.Println("Load config completed.")

	// Init MongoDB
	repo, err := repository.NewMongoRepo(config)
	if err != nil {
		fmt.Printf("Init MongoDB failed: %v", err)
	}
	defer repo.Disconnect()
	log.Println("Init Mongo DB completed.")

	// Init usecase
	usecase := usecase.NewUsecase(repo, config)
	log.Println("Init usecase completed.")

	// Init delivery
	delivery := http.NewHttpDelivery(usecase, config)
	log.Println("Init delivery completed.")

	delivery.Run(config.ServerPort)
}
