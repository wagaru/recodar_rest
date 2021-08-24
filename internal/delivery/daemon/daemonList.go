package daemon

import (
	"encoding/json"

	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/logger"
)

func (d *Daemon) messageBrokerAuthUserDaemon() {
	go func() {
		meta := &domain.RabbitMQMeta{
			ExchangeType: "fanout",
			ExchangeName: "fanout",
		}
		msgs, err := d.messageBrokerUsecase.ConsumeMessages(meta)
		if err != nil {
			logger.Logger.Printf("consume messages got error:%v", err)
			return
		}
		for msg := range msgs {
			var message domain.MessageUserLogin
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				logger.Logger.Printf("Decode message failed:%v", err)
				continue
			}
			logger.Logger.Printf("I got %v", message.ID)
			msg.Ack(false)
		}
	}()
}

func (d *Daemon) messageBrokerAuthUserGoogleDaemon() {
	go func() {
		meta := &domain.RabbitMQMeta{
			ExchangeType: "direct",
			ExchangeName: "userAuth",
			BindingKey:   "google",
		}
		msgs, err := d.messageBrokerUsecase.ConsumeMessages(meta)
		if err != nil {
			logger.Logger.Printf("consume messages got error:%v", err)
			return
		}
		for msg := range msgs {
			var message domain.MessageUserLogin
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				logger.Logger.Printf("Decode message failed:%v", err)
				continue
			}
			logger.Logger.Printf("I got %v from google exchange", message.ID)
			msg.Ack(false)
		}
	}()
}

func (d *Daemon) messageBrokerAuthUserLineDaemon() {
	go func() {
		meta := &domain.RabbitMQMeta{
			ExchangeType: "direct",
			ExchangeName: "userAuth",
			BindingKey:   "line",
		}
		msgs, err := d.messageBrokerUsecase.ConsumeMessages(meta)
		if err != nil {
			logger.Logger.Printf("consume messages got error:%v", err)
			return
		}
		for msg := range msgs {
			var message domain.MessageUserLogin
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				logger.Logger.Printf("Decode message failed:%v", err)
				continue
			}
			logger.Logger.Printf("I got %v from line exchange", message.ID)
			msg.Ack(false)
		}
	}()
}
