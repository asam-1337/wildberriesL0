package main

import (
	"context"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/cache"
	"github.com/asam-1337/wildberriesL0/internal/delivery"
	"github.com/asam-1337/wildberriesL0/internal/nats"
	"github.com/asam-1337/wildberriesL0/internal/pgxbalancer"
	"github.com/asam-1337/wildberriesL0/internal/render"
	"github.com/asam-1337/wildberriesL0/internal/repository"
	"github.com/asam-1337/wildberriesL0/internal/service"
	"github.com/asam-1337/wildberriesL0/pkg/httpServer"
	"github.com/asam-1337/wildberriesL0/pkg/logger"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	logger.InitLogger()
	cfg, err := config.InitConfig()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on init config")
	}

	c := cache.NewCache(time.Hour, time.Minute)
	balancer, err := pgxbalancer.NewTransactionBalancer(context.Background(), cfg.Pg)
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on balancer init")
	}

	repo := repository.NewOrdersRepository(balancer)
	broker := nats.NewBroker(cfg.Stan, c, repo)
	go func() {
		err := broker.Subscribe()
		if err != nil {
			log.WithField("err", err.Error()).Fatal("error occurred on broker init")
		}
	}()

	rnd := render.NewRenderService()
	svc := service.NewService(c, repo)
	handler := delivery.NewHandler(rnd, svc)
	router := handler.InitRoutes()

	s := httpServer.NewHttpServer("8080", router)
	log.WithField("port", cfg.Port).Info("starting server")

	err = s.ListenAndServe()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("cant starting server")
		return
	}
}
