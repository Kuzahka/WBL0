package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func connectDB() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://user:password@localhost:5433/mydb")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to PostgreSQL")

	// Создаем таблицу, если её нет
	err = createTable()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
	    order_uid VARCHAR PRIMARY KEY,
	    data JSONB
	);
	`
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	log.Println("Table 'orders' ensured to exist")
	return nil
}
func saveToDB(order Order) error {
	data, err := json.Marshal(order) // Преобразуем данные в JSON
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (order_uid, data) VALUES ($1, $2) ON CONFLICT (order_uid) DO NOTHING`
	_, err = db.Exec(context.Background(), query, order.OrderUID, data)
	if err != nil {
		return err
	}
	log.Printf("Order %s saved to database", order.OrderUID)
	return nil
}
