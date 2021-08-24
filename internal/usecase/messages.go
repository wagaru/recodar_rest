package usecase

import (
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/repository"
)

type MessageBrokerUsecase interface {
	SendMessages(*domain.RabbitMQMeta, []byte) error
	ConsumeMessages(*domain.RabbitMQMeta) (repository.MessageBrokerMessage, error)
}

type messageBrokerUsecase struct {
	repo   repository.MessageBrokerRepository
	config *config.Config
}

func NewMessageBrokerUsecase(repo repository.MessageBrokerRepository, config *config.Config) MessageBrokerUsecase {
	return &messageBrokerUsecase{
		repo:   repo,
		config: config,
	}
}

func (usecase *messageBrokerUsecase) SendMessages(meta *domain.RabbitMQMeta, message []byte) error {
	return usecase.repo.SendMessages(meta, message)
}

func (usecase *messageBrokerUsecase) ConsumeMessages(meta *domain.RabbitMQMeta) (repository.MessageBrokerMessage, error) {
	return usecase.repo.ConsumeMessages(meta)
}
