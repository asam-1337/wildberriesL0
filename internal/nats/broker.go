package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/nats-io/stan.go"
	"log"
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
	}
)

func NewBroker(cfg config.StanConfig, cache CacheService, repo Repository) *Broker {
	return &Broker{
		cfg: cfg,
	}
}

func (b *Broker) Subscribe() error {
	sc, err := stan.Connect(b.cfg.ClusterID, b.cfg.ClientID)
	if err != nil {
		return fmt.Errorf("cant connect to nats server: %s", err)
	}
	defer func() {
		err := sc.Close()
		if err != nil {
			log.Printf("cant close conn: %s", err.Error())
		}
	}()

	dur, err := time.ParseDuration("1s")
	if err != nil {
		log.Println("cant parse time")
		return err
	}

	sub, err := sc.Subscribe(
		b.cfg.ChannelID,
		b.receiveHandler,
		stan.DurableName("durable"),
		stan.SetManualAckMode(),
		stan.AckWait(dur))
	if err != nil {
		return fmt.Errorf("cant subscribe: %s", err.Error())
	}

	defer func() {
		err = sub.Unsubscribe()
		if err != nil {
			log.Println(err.Error())
		}

		err = sub.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}

func (b *Broker) receiveHandler(m *stan.Msg) {
	err := m.Ack()
	if err != nil {
		log.Println(err.Error())
		return
	}

	order, err := b.validateMessage(m.Data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if exist := b.cache.Exist(order.OrderUID); exist {
		log.Println("order already exist")
		return
	}

	ctx := context.Background()
	b.cache.Store(order.OrderUID, order)
	err = b.repo.Insert(ctx, order)
	if err != nil {
		log.Println(err.Error())
	}
}

func (b *Broker) validateMessage(msg []byte) (entity.Order, error) {
	order := entity.Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		return entity.Order{}, fmt.Errorf("cant unmarshal json: %s", err.Error())
	}

	val := reflect.ValueOf(order).Elem()
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsZero() && len(order.OrderUID) != 19 {
			return entity.Order{}, err
		}
	}
	return order, nil
}

func (b *Broker) restoreCash() error {
	orders, err := b.repo.SelectAll(context.Background())
	if err != nil {
		return err
	}

	for _, order := range orders {
		b.cache.Store(order.OrderUID, order)
	}
	return nil
}
