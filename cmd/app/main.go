package main

import (
	"github.com/asam-1337/wildberriesL0/internal/handler"
	"github.com/asam-1337/wildberriesL0/pkg/httpServer"
)

func main() {
	handler := handler.NewHandler()
	s := httpServer.NewHttpServer("8080", handler)
	s.ListenAndServe()
}
