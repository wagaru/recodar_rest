package daemon

import (
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/logger"
	"github.com/wagaru/recodar-rest/internal/usecase"
)

type Daemon struct {
	messageBrokerUsecase usecase.MessageBrokerUsecase
	config               *config.Config
	daemonLists          map[string]func()
}

func NewDaemon(messageBrokerUsecase usecase.MessageBrokerUsecase, config *config.Config) *Daemon {
	daemon := &Daemon{
		messageBrokerUsecase: messageBrokerUsecase,
		config:               config,
		daemonLists:          map[string]func(){},
	}
	daemon.RegisterDaemon()
	return daemon
}

func (d *Daemon) RegisterDaemon() {
	d.Register("messageBrokerAuthUserDaemon", d.messageBrokerAuthUserDaemon)
	d.Register("messageBrokerAuthUserGoogleDaemon", d.messageBrokerAuthUserGoogleDaemon)
	d.Register("messageBrokerAuthUserLineDaemon", d.messageBrokerAuthUserLineDaemon)
}

func (d *Daemon) Register(name string, process func()) {
	if _, ok := d.daemonLists[name]; !ok {
		d.daemonLists[name] = process
	}
}

func (d *Daemon) Run() {
	for name, process := range d.daemonLists {
		logger.Logger.Printf("Run %s daemon\n", name)
		process()
	}
}
