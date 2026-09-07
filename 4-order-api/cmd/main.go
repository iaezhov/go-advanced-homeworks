package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	database := db.NewDb(config)

	// Repositories
	productRepository := product.NewProductRepository(database)

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{ProductRepository: productRepository})

	// Server
	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}
	fmt.Printf("Server is listening on port %v\n", config.AppPort)
	server.ListenAndServe()

}
