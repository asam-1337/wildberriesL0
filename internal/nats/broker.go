package nats

import (
	"encoding/json"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/nats-io/stan.go"
)

type Broker struct {
	cfg config.StanConfig
}

func NewBroker(cfg config.StanConfig) *Broker {
	return &Broker{
		cfg: cfg,
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
		return err
	}

	err = sub.Unsubscribe()
	err = sub.Close()

	return err
}

func (b *Broker) receiveOrder(m *stan.Msg) {
	order := entity.Order{}
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		//b.log.Error(fmt.Errorf("cant unmarshal json: %s", err.Error()))
		return
	}

	return
}
