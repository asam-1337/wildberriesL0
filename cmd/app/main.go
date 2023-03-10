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
	"log"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}

	c := cache.NewCache(time.Hour, time.Minute)
	balancer, err := pgxbalancer.NewTransactionBalancer(context.Background(), cfg.Pg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	repo := repository.NewOrdersRepository(balancer)
	broker := nats.NewBroker(cfg.Stan, c, repo)
	go func() {
		err := broker.Subscribe()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	rnd := render.NewRenderService()
	svc := service.NewService(c, repo)
	handler := delivery.NewHandler(rnd, svc)
	router := handler.InitRoutes()

	s := httpServer.NewHttpServer("8080", router)
	log.Println("starting server at :8080")
	s.ListenAndServe()
}
