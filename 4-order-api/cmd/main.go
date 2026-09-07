package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	database := db.NewDb(config)

	product.NewProductHandler(router, product.ProductHandlerDeps{Db: database})

	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}
	fmt.Printf("Server is listening on port %v\n", config.AppPort)
	server.ListenAndServe()

}

func startServer(config *configs.Config, router *http.ServeMux) {
	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	fmt.Printf("Server is listening on port %v/n", config.AppPort)
}
