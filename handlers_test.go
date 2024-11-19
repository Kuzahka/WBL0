package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderHandler(t *testing.T) {
	cache := make(map[string]Order)
	order := Order{OrderUID: "b563feb7b2b84b6test"}
	cache["b563feb7b2b84b6test"] = order

	// Создаём запрос
	req, err := http.NewRequest("GET", "/b563feb7b2b84b6test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Создаём тестовый сервер
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getOrder)
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа
	var returnedOrder Order
	if err := json.NewDecoder(rr.Body).Decode(&returnedOrder); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if returnedOrder.OrderUID != order.OrderUID {
		t.Errorf("Handler returned wrong order: got %v, want %v", returnedOrder.OrderUID, order.OrderUID)
	}
}
