package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if order, found := cache[id]; found {
		json.NewEncoder(w).Encode(order)
	} else {
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}
