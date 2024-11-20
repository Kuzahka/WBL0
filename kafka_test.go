package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestConsumeMessages(t *testing.T) {
	// Настраиваем тестовые параметры Kafka
	const (
		topic     = "test-orders"
		groupID   = "test-order-service"
		broker    = "localhost:9092"
		testOrder = `{"order_uid":"test123","items":[{"name":"item1","price":100}]}`
	)

	// Создаём тестовую тему
	err := createTestTopic(broker, topic)
	if err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}
	defer deleteTestTopic(broker, topic)

	// Запускаем обработчик в отдельной горутине
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		})

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("could not read message: %v", err)
				return
			}

			log.Printf("Received message: %s", string(m.Value))

			// Симуляция обработки
			_, err = validateOrder(m.Value)
			if err != nil {
				log.Printf("Invalid order data: %v", err)
			} else {
				t.Log("Message processed successfully!")
				return
			}
		}
	}()

	// Отправляем тестовое сообщение
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer w.Close()

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(testOrder),
	})
	if err != nil {
		t.Fatalf("Failed to send test message: %v", err)
	}

	// Ждём, чтобы обработчик успел прочитать сообщение
	time.Sleep(2 * time.Second)
}
func createTestTopic(broker, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	return err
}
func deleteTestTopic(broker, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.DeleteTopics(topic)
}
