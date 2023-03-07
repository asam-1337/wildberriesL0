package main

import (
	"github.com/asam-1337/wildberriesL0/internal/handler"
	"github.com/asam-1337/wildberriesL0/pkg/httpServer"
	"html/template"
	"log"
)

func main() {
	tmp, err := template.ParseFiles("./templates/index.html", "./templates/find.html")
	if err != nil {
		log.Fatalf(err.Error())
	}

	handler := handler.NewHandler(tmp)
	router := handler.InitRoutes()

	s := httpServer.NewHttpServer("8080", router)
	log.Println("starting server at :8080")
	s.ListenAndServe()
}
