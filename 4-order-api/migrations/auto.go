package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
)

func main() {
	config := configs.LoadConfig()
	database := db.NewDb(config)
	database.AutoMigrate(&product.Product{})
}
