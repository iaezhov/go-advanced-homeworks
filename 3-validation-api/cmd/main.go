package main

import (
	"3-validation-api/3-validation-api/configs"
	"3-validation-api/3-validation-api/internal/verify"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{Config: config})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
