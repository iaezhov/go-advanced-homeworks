package main

import (
	"4-order-api/configs"
	"4-order-api/internal/auth"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
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
	jwtPackage := jwt.NewJWT(config.Auth.Secret)

	// Middlewares
	mds := middleware.Chain(middleware.CORS, middleware.Logging)

	// Repositories
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)

	// Services
	authService := auth.NewAuthService(&auth.AuthServiceDeps{
		UserRepository: userRepository,
	})

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		JWT:         jwtPackage,
	})

	// Server
	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: mds(router),
	}
	fmt.Printf("Server is listening on port %v\n", config.AppPort)
	server.ListenAndServe()

}
