package nats

import (
	"encoding/json"
	"fmt"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/pkg/logger"
	"github.com/nats-io/stan.go"
)

type Broker struct {
	cfg config.Stan
	log logger.Logger
}

func NewBroker(cfg config.Stan, log logger.Logger) *Broker {
	return &Broker{
		cfg: cfg,
		log: log,
	}
}

func (b *Broker) ReceiveMessage() error {
	sc, err := stan.Connect(b.cfg.ClusterID, b.cfg.ClientID)
	if err != nil {
		return err
	}
	defer sc.Close()

	sub, err := sc.Subscribe(b.cfg.ChannelID, b.receiveOrder)

	if err != nil {
		b.log.Error(err)
		return err
	}

	err = sub.Unsubscribe()
	err = sub.Close()

	return err
}

func (b *Broker) receiveOrder(m *stan.Msg) {
	order := &entity.Order{}
	err := json.Unmarshal(m.Data, order)
	if err != nil {
		b.log.Error(fmt.Errorf("cant unmarshal json: %s", err.Error()))
		return
	}

	return
}
