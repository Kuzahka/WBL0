package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)

// Пример заказа
var testOrder = Order{
	OrderUID: "b563feb7b2b84b6test",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "b563feb7b2b84b6test",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
}

func TestKafkaConsumer(t *testing.T) {
	// Создаём мок Kafka Consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	mockConsumer := mocks.NewConsumer(t, config)

	// Настраиваем ожидание подписки на топик "orders"
	mockPartitionConsumer := mockConsumer.ExpectConsumePartition("orders", 0, sarama.OffsetOldest)

	// Генерируем тестовые данные (JSON)
	data, err := json.Marshal(testOrder)
	if err != nil {
		t.Fatalf("Failed to marshal test order: %v", err)
	}

	// Отправляем сообщение в партицию
	mockPartitionConsumer.YieldMessage(&sarama.ConsumerMessage{
		Topic:     "orders",
		Partition: 0,
		Offset:    0,
		Value:     data,
	})

	// Вызов тестируемой функции
	err = consumeKafkaWithMock(mockConsumer)
	if err != nil {
		t.Errorf("Kafka consumer test failed: %v", err)
	}

	// Проверяем, что все ожидания выполнены
	mockConsumer.Close()
}

// Функция с mock-реализацией Kafka Consumer
func consumeKafkaWithMock(consumer sarama.Consumer) error {
	partitionConsumer, err := consumer.ConsumePartition("orders", 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var order Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshalling Kafka message: %v", err)
			continue
		}

		// В тесте вы можете заменить реальные функции на моки
		log.Printf("Received order: %v", order.OrderUID)
	}
	return nil
}
