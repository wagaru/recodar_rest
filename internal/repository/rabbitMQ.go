package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
)

type MessageBrokerRepository interface {
	Disconnect()
	SendMessages(*domain.RabbitMQMeta, []byte) error
	ConsumeMessages(*domain.RabbitMQMeta) (MessageBrokerMessage, error)
}

type RabbitMQRepo struct {
	connection  *amqp.Connection
	channelPool map[string]*amqp.Channel
}

type MessageBrokerMessage <-chan amqp.Delivery

func NewRabbitMQRepo(config *config.Config) (MessageBrokerRepository, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%v",
		config.RabbitMQUserName,
		config.RabbitMQPassword,
		config.RabbitMQHost,
		config.RabbitMQPort,
	)
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("create rabbitMQ connection failed:%w", err)
	}
	defaultChannel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("create rabbitMQ channel failed:%w", err)
	}
	return &RabbitMQRepo{
		connection: connection,
		channelPool: map[string]*amqp.Channel{
			"default": defaultChannel,
		},
	}, nil
}

func (r *RabbitMQRepo) Disconnect() {
	for _, channel := range r.channelPool {
		err := channel.Close()
		if err != nil {
			log.Printf("close rabbitmq channel failed: %v", err)
		}
	}
	err := r.connection.Close()
	if err != nil {
		log.Printf("close rabbitmq connection failed: %v", err)
	}
}

func (r *RabbitMQRepo) SendMessages(meta *domain.RabbitMQMeta, message []byte) error {
	channelName := meta.ChannelName
	if channelName == "" {
		channelName = "default"
	}
	var ch *amqp.Channel
	if v, ok := r.channelPool[channelName]; ok {
		ch = v
	} else {
		ch, err := r.connection.Channel()
		if err != nil {

		}
		r.channelPool[channelName] = ch
	}
	if meta.ExchangeType == "" {
		return errors.New(fmt.Sprintf("Invalid exchange type provided: %v", meta.ExchangeType))
	}
	err := ch.ExchangeDeclare(meta.ExchangeName, meta.ExchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Declare rabbitMQ exchange failed.%w", err)
	}
	err = ch.Publish(meta.ExchangeName, meta.RoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		return fmt.Errorf("Publish rabbitMQ messages failed.%w", err)
	}
	return nil
}

func (r *RabbitMQRepo) ConsumeMessages(meta *domain.RabbitMQMeta) (MessageBrokerMessage, error) {
	channelName := meta.ChannelName
	if channelName == "" {
		channelName = "default"
	}
	var ch *amqp.Channel
	if v, ok := r.channelPool[channelName]; ok {
		ch = v
	} else {
		ch, err := r.connection.Channel()
		if err != nil {

		}
		r.channelPool[channelName] = ch
	}
	if meta.ExchangeType == "" {
		return nil, errors.New(fmt.Sprintf("Invalid exchange type provided: %v", meta.ExchangeType))
	}
	err := ch.ExchangeDeclare(meta.ExchangeName, meta.ExchangeType, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Declare rabbitMQ exchange failed.%w", err)
	}
	queue, err := ch.QueueDeclare(meta.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Declare rabbitMQ queue failed.%w", err)
	}
	err = ch.QueueBind(queue.Name, meta.BindingKey, meta.ExchangeName, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Bind rabbitMQ queue failed.%w", err)
	}
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Consume rabbitMQ queue failed.%w", err)
	}
	return msgs, nil
}
