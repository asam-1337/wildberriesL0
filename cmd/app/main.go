package main

import (
	"context"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/cache"
	"github.com/asam-1337/wildberriesL0/internal/delivery"
	"github.com/asam-1337/wildberriesL0/internal/nats"
	"github.com/asam-1337/wildberriesL0/internal/render"
	"github.com/asam-1337/wildberriesL0/internal/repository"
	"github.com/asam-1337/wildberriesL0/internal/service"
	"github.com/asam-1337/wildberriesL0/pkg/httpServer"
	_ "github.com/asam-1337/wildberriesL0/pkg/logger"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on init config")
	}

	c := cache.NewCache(time.Hour, 30*time.Minute)
	pool, err := repository.NewPgxPool(context.Background(), cfg.Pg)
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on db init")
	}

	repo := repository.NewOrdersRepository(pool)
	broker, err := nats.NewBroker(cfg.Stan, c, repo)
	if err != nil {
		log.Fatal(err)
	}

	err = broker.Subscribe()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on broker init")
	}

	rnd, err := render.NewRenderService()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("error occurred on render init")
	}

	svc := service.NewService(c, repo)
	handler := delivery.NewHandler(rnd, svc)
	router := handler.InitRoutes()

	s := httpServer.NewHttpServer(cfg.Port, router)
	log.WithField("port", cfg.Port).Info("starting server")

	err = s.ListenAndServe()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("cant starting server")
		return
	}
}
