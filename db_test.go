package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

const testDBConnString = "host=localhost port=5433 user=user password=password dbname=mydb sslmode=disable"

func TestDBConnection(t *testing.T) {
	var err error

	// Подключаемся к базе данных
	testDB, err = sql.Open("postgres", testDBConnString)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Пингуем базу данных, чтобы убедиться в доступности
	err = testDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping the database: %v", err)
	}

	// Проверяем успешное подключение
	t.Log("Database connection successful")
}

func TestDBCreateTable(t *testing.T) {
	// Проверяем, что соединение установлено
	if testDB == nil {
		t.Fatal("Database connection is not initialized")
	}

	// Создаём таблицу для тестов
	query := `
		CREATE TABLE IF NOT EXISTS test_orders (
			order_uid VARCHAR PRIMARY KEY,
			data JSONB
		)`
	_, err := testDB.Exec(query)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Проверка успешного выполнения запроса
	t.Log("Test table created successfully")
}

func TestDBInsertData(t *testing.T) {
	// Проверяем, что соединение установлено
	if testDB == nil {
		t.Fatal("Database connection is not initialized")
	}

	// Тестовые данные
	orderUID := "b563feb7b2b84b6test"
	jsonData := `{"order_uid":"b563feb7b2b84b6test"}`

	// Вставляем данные
	_, err := testDB.Exec("INSERT INTO test_orders (order_uid, data) VALUES ($1, $2)", orderUID, jsonData)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Проверка успешной вставки
	t.Log("Data inserted successfully")
}

func TestDBFetchData(t *testing.T) {
	// Проверяем, что соединение установлено
	if testDB == nil {
		t.Fatal("Database connection is not initialized")
	}

	// Извлекаем данные
	var orderUID string
	var jsonData string
	err := testDB.QueryRow("SELECT order_uid, data FROM test_orders WHERE order_uid = $1", "b563feb7b2b84b6test").Scan(&orderUID, &jsonData)
	if err != nil {
		t.Fatalf("Failed to fetch data: %v", err)
	}

	// Проверка данных
	if orderUID != "b563feb7b2b84b6test" {
		t.Errorf("Data mismatch: got orderUID=%s, data=%s", orderUID, jsonData)
	} else {
		t.Log("Data fetched successfully")
	}
}

func TestDBCleanup(t *testing.T) {
	// Проверяем, что соединение установлено
	if testDB == nil {
		t.Fatal("Database connection is not initialized")
	}

	// Удаляем тестовые данные и таблицу
	_, err := testDB.Exec("DROP TABLE IF EXISTS test_orders")
	if err != nil {
		t.Fatalf("Failed to drop test table: %v", err)
	}

	// Проверка успешного выполнения
	t.Log("Test table dropped successfully")
}
