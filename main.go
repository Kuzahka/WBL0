package main

import (
	"context"
	"net/http"

	"log"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

var db *pgx.Conn
var cache = make(map[string]Order)

func main() {

	connectDB()
	defer db.Close(context.Background())
	err := loadCacheFromDB()
	if err != nil {
		log.Fatalf("Failed to load cache from DB: %v", err)
	}
	// Запуск Kafka-подписки

	go consumeMessages()
	// Запуск сервера
	r := mux.NewRouter()
	r.HandleFunc("/{id}", getOrder).Methods("GET")

	http.ListenAndServe(":8080", r)

}
