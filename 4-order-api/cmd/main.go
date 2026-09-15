package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"4-order-api/pkg/middleware"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	database := db.NewDb(config)

	// Middlewares
	mds := middleware.Chain(middleware.CORS, middleware.Logging)

	// Repositories
	productRepository := product.NewProductRepository(database)

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{ProductRepository: productRepository})

	// Server
	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: mds(router),
	}
	fmt.Printf("Server is listening on port %v\n", config.AppPort)
	server.ListenAndServe()

}
