package nats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/cache"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"github.com/asam-1337/wildberriesL0/internal/repository"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

type Broker struct {
	cfg   config.StanConfig
	repo  repository.Repository
	cache cache.Service
	sc    stan.Conn
	sub   stan.Subscription
}

func NewBroker(cfg config.StanConfig, cache cache.Service, repo repository.Repository) (*Broker, error) {
	var err error
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID)
	if err != nil {
		log.WithField("err", err.Error()).Error("cant connect to nats-streaming-server")
		return nil, err
	}

	return &Broker{
		cfg:   cfg,
		cache: cache,
		repo:  repo,
		sc:    sc,
	}, nil
}

func (b *Broker) Subscribe() error {
	dur, err := time.ParseDuration("1s")
	if err != nil {
		log.WithField("err", err.Error()).Error("cant parse time")
		return err
	}

	b.sub, err = b.sc.Subscribe(
		b.cfg.ChannelID,
		b.receiveHandler,
		stan.DurableName("durable"),
		stan.SetManualAckMode(),
		stan.AckWait(dur))
	if err != nil {
		log.WithField("err", err.Error()).Error("cant subscribe")
		return err
	}

	log.Infof("succesfully subscribe")
	return nil
}

func (b *Broker) Publish(order entity.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return b.sc.Publish(b.cfg.ChannelID, data)
}

func (b *Broker) Close() {
	err := b.sc.Close()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant close conn")
	}
}

func (b *Broker) Unsubscribe() {
	var err error
	err = b.sub.Unsubscribe()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant unsubscribe")
	}
}

func (b *Broker) receiveHandler(m *stan.Msg) {
	err := m.Ack()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant acknowledge")
		return
	}

	order, ok := validate(m.Data)
	if !ok {
		log.Info("message is invalid")
		return
	}

	if exist := b.cache.Exist(order.OrderUID); exist {
		log.Info("order has already existed")
		return
	}

	ctx := context.Background()
	b.cache.Store(order.OrderUID, order)
	err = b.repo.Insert(ctx, order)
	if err != nil {
		if errors.Is(err, localErrors.ErrAlreadyExists) {
			log.WithField("info", err.Error()).Info("cant insert in repository")
			return
		}
		log.WithField("err", err.Error()).Error("cant insert in repository")
	}
}

func (b *Broker) restoreCash() error {
	orders, err := b.repo.SelectAll(context.Background())
	if err != nil {
		log.WithField("err", err).Error("cant select from db")
		return err
	}

	for _, order := range orders {
		b.cache.Store(order.OrderUID, order)
	}
	return nil
}

func validate(msg []byte) (entity.Order, bool) {
	order := entity.Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.WithField("err", err).Warn("json error")
		return entity.Order{}, false
	}

	val := reflect.ValueOf(&order).Elem()
	for i := 0; i < val.NumField(); i++ {
		if len(order.OrderUID) != 19 || val.Field(i).IsZero() && !strings.Contains(val.Type().Field(i).Tag.Get("json"), "notrequired") {
			return entity.Order{}, false
		}
	}
	return order, true
}
