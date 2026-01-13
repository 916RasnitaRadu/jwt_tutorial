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

	r.HandleFunc("/login", HandleLogin).Methods("POST")

	port := "8080"
	log.Printf("the service is listening on: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
