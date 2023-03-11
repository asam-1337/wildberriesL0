package nats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type (
	CacheService interface {
		Store(key string, value entity.Order)
		Load(key string) (value entity.Order, loaded bool)
		Exist(key string) (loaded bool)
	}

	Repository interface {
		Insert(ctx context.Context, order entity.Order) error
		SelectById(ctx context.Context, id string) (entity.Order, error)
		SelectAll(ctx context.Context) ([]entity.Order, error)
	}

	Broker struct {
		cfg   config.StanConfig
		repo  Repository
		cache CacheService
		sc    stan.Conn
		sub   stan.Subscription
	}
)

func NewBroker(cfg config.StanConfig, cache CacheService, repo Repository) *Broker {
	return &Broker{
		cfg:   cfg,
		cache: cache,
		repo:  repo,
	}
}

func (b *Broker) Subscribe() error {
	var err error
	b.sc, err = stan.Connect(b.cfg.ClusterID, b.cfg.ClientID)
	if err != nil {
		log.WithField("err", err.Error()).Error("cant connect to nats-streaming-server")
		return err
	}

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

func (b *Broker) Close() {
	var err error
	err = b.sub.Unsubscribe()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant unsubscribe")
	}

	err = b.sub.Close()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant close subscription")
	}

	err = b.sc.Close()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant close conn")
	}
}

func (b *Broker) receiveHandler(m *stan.Msg) {
	err := m.Ack()
	if err != nil {
		log.WithField("err", err.Error()).Warn("cant acknowledge")
		return
	}

	order, ok := b.validateMessage(m.Data)
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

func (b *Broker) validateMessage(msg []byte) (entity.Order, bool) {
	order := entity.Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		return entity.Order{}, false
	}

	val := reflect.ValueOf(order).Elem()
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsZero() && len(order.OrderUID) != 19 {
			return entity.Order{}, false
		}
	}
	return order, true
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
