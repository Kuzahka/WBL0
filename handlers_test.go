package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetOrderHandler(t *testing.T) {
	orderJSON :=
		`{
	   "order_uid": "b563feb7b2b84b6test",
	   "track_number": "WBILMTESTTRACK",
	   "entry": "WBIL",
	   "delivery": {
		  "name": "Test Testov",
		  "phone": "+9720000000",
		  "zip": "2639809",
		  "city": "Kiryat Mozkin",
		  "address": "Ploshad Mira 15",
		  "region": "Kraiot",
		  "email": "test@gmail.com"
	   },
	   "payment": {
		  "transaction": "b563feb7b2b84b6test",
		  "request_id": "",
		  "currency": "USD",
		  "provider": "wbpay",
		  "amount": 1817,
		  "payment_dt": 1637907727,
		  "bank": "alpha",
		  "delivery_cost": 1500,
		  "goods_total": 317,
		  "custom_fee": 0
	   },
	   "items": [
		  {
			 "chrt_id": 9934930,
			 "track_number": "WBILMTESTTRACK",
			 "price": 453,
			 "rid": "ab4219087a764ae0btest",
			 "name": "Mascaras",
			 "sale": 30,
			 "size": "0",
			 "total_price": 317,
			 "nm_id": 2389212,
			 "brand": "Vivienne Sabo",
			 "status": 202
		  }
	   ],
	   "locale": "en",
	   "internal_signature": "",
	   "customer_id": "test",
	   "delivery_service": "meest",
	   "shardkey": "9",
	   "sm_id": 99,
	   "date_created": "2021-11-26T06:22:19Z",
	   "oof_shard": "1"
	}`

	var order Order

	err := json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {

		return
	}
	cache[order.OrderUID] = order
	// Создаём запрос
	req, err := http.NewRequest("GET", "/b563feb7b2b84b6test", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Создаём тестовый сервер
	rr := httptest.NewRecorder()

	// Создаем роутер и регистрируем маршрут
	r := mux.NewRouter()
	r.HandleFunc("/{id}", getOrder).Methods("GET")

	// Передаём запрос в роутер
	r.ServeHTTP(rr, req)

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
