package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

func consumeMessages() {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.GroupID,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}

		log.Printf("Received message: %s", string(m.Value))

		// Парсим сообщение в модель
		order, err := validateOrder(m.Value)
		if err != nil {
			log.Printf("Invalid order data: %v", err)
			continue
		}
		cache[order.OrderUID] = *order
		// Сохраняем в базу данных
		err = saveToDB(*order)
		if err != nil {
			log.Printf("Failed to save order to DB: %v", err)
		}
	}
}
func validateOrder(data []byte) (*Order, error) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return nil, errors.New("failed to parse JSON")
	}

	// Проверка обязательных полей
	if order.OrderUID == "" {
		return nil, errors.New("missing field: order_uid")
	}
	if len(order.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	return &order, nil
}
func loadCacheFromDB() error {
	rows, err := db.Query(context.Background(), "SELECT order_uid, data FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		var jsonData []byte
		err := rows.Scan(&order.OrderUID, &jsonData)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		err = json.Unmarshal(jsonData, &order)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v", err)
			continue
		}

		// Записываем в кэш
		cache[order.OrderUID] = order
	}
	log.Println("Cache restored from database")
	return nil
}
