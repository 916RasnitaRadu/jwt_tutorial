package main

import (
	"encoding/json"
	"net/http"
)

func HandleGreet(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello my friend. Looks like u r authorized.")
}
