package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func RandomHandler(w http.ResponseWriter, r *http.Request) {
	number := getRandomNumber()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, number)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/rand", RandomHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("ServerPort: 8080")
	server.ListenAndServe()
}

func getRandomNumber() int {
	return rand.Intn(6) + 1
}
