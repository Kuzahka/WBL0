package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSaveToDB(t *testing.T) {
	// Инициализируем мок
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize mock DB: %v", err)
	}
	defer db.Close()

	// Настраиваем ожидания
	order := Order{OrderUID: "testUID"}
	mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID).WillReturnResult(sqlmock.NewResult(1, 1))

	// Тестируем функцию
	if err := saveToDB(order); err != nil {
		t.Errorf("saveToDB failed: %v", err)
	}

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}
