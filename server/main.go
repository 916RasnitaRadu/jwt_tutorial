package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var JwtKey []byte

func main() {
	JwtKey = []byte("supersecretkey")

	r := mux.NewRouter()
	r.Use(jwtAuthMiddleware)
	r.HandleFunc("/hello", HandleGreet).Methods("GET")

	port := "8081"
	log.Printf("auth service listening on: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
