package main

import (
	"4-order-api/configs"
	"4-order-api/internal/auth"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initConfig() *configs.Config {
	return &configs.Config{
		Db: configs.DbConfig{
			DBHost:     "localhost",
			DBPort:     5433,
			DBUser:     "postgres",
			DBPassword: "my_pass",
			DBName:     "orders_test",
			DBSSLMode:  "disable",
		},
		Auth: configs.AuthConfig{
			Secret: "/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w=",
		},
	}
}

func initServer(config *configs.Config) *httptest.Server {
	ts := httptest.NewServer(App(config))
	return ts
}

func initDb(config *configs.Config) *gorm.DB {
	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&product.Product{}, &user.User{}, &order.Order{}, &order.OrderItem{})
	if err != nil {
		panic(err)
	}
	return db
}

func initTestEnviroment() (*httptest.Server, *gorm.DB) {
	config := initConfig()
	ts := initServer(config)
	db := initDb(config)
	return ts, db
}

func getToken(t *testing.T, ts *httptest.Server, db *gorm.DB) string {
	phone := "79001112233"

	// Auth Login: отправляем телефон, получаем sessionId
	loginPayload, _ := json.Marshal(&auth.LoginPayload{
		Phone: phone,
	})
	loginRes, err := http.Post(ts.URL+"/login", "application/json", bytes.NewReader(loginPayload))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	if loginRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(loginRes.Body)
		t.Fatalf("Expected 200 for login, got %d: %s", loginRes.StatusCode, string(body))
	}

	var loginData auth.LoginResponse
	if err := json.NewDecoder(loginRes.Body).Decode(&loginData); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	sessionId := loginData.SessionId

	// Извлекаем код из БД (имитируем получение SMS)
	var u user.User
	if err := db.Where("phone = ?", phone).First(&u).Error; err != nil {
		t.Fatalf("Could not find user in DB to get code: %v", err)
	}
	code := u.Code

	// Auth Verify: отправляем код и sessionId, получаем JWT токен
	verifyPayload, _ := json.Marshal(&auth.VerifyPayload{
		SessionId: sessionId,
		Code:      code,
	})
	verifyRes, err := http.Post(ts.URL+"/verify", "application/json", bytes.NewReader(verifyPayload))
	if err != nil {
		t.Fatalf("Verify request failed: %v", err)
	}
	if verifyRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(verifyRes.Body)
		t.Fatalf("Expected 200 for verify, got %d: %s", verifyRes.StatusCode, string(body))
	}

	var verifyData auth.VerifyResponse
	if err := json.NewDecoder(verifyRes.Body).Decode(&verifyData); err != nil {
		t.Fatalf("Failed to decode verify response: %v", err)
	}
	return verifyData.Token
}

func TestOrderCreateSuccess(t *testing.T) {
	ts, db := initTestEnviroment()
	defer ts.Close()
	dbTestData := NewDbTestData(db)
	dbTestData.Init(t)
	defer dbTestData.Clear()
	token := getToken(t, ts, db)

	// Order Create: создаем заказ с полученным токеном
	var p product.Product
	if err := db.First(&p).Error; err != nil {
		t.Fatalf("Could not find product in DB: %v", err)
	}

	orderPayload, _ := json.Marshal(&order.OrderCreatePayload{
		Items: []order.OrderItemRequest{
			{ProductID: p.ID, Quantity: 2},
		},
	})

	req, _ := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(orderPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	orderRes, err := client.Do(req)
	if err != nil {
		t.Fatalf("Order request failed: %v", err)
	}
	defer orderRes.Body.Close()

	if orderRes.StatusCode != http.StatusCreated && orderRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(orderRes.Body)
		t.Fatalf("Expected 201/200 for order creation, got %d: %s", orderRes.StatusCode, string(body))
	}

	var createdOrder order.Order
	if err := json.NewDecoder(orderRes.Body).Decode(&createdOrder); err != nil {
		t.Fatalf("Failed to decode order response: %v", err)
	}

	expectedPrice := p.Price * 2
	if createdOrder.Price != expectedPrice {
		t.Errorf("Expected order price %v, got %v", expectedPrice, createdOrder.Price)
	}

	if len(createdOrder.Items) != 1 {
		t.Errorf("Expected 1 order item, got %d", len(createdOrder.Items))
	}
}

// INIT DB DATA
type DbTestData struct {
	DB *gorm.DB
}

func NewDbTestData(db *gorm.DB) *DbTestData {
	return &DbTestData{
		DB: db,
	}
}

func (td *DbTestData) Init(t *testing.T) {
	// Create User explicitly using NewUser to ensure session and code are generated
	u := user.NewUser("79001112233")
	if err := td.DB.Create(u).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	products := []product.Product{
		{Name: "Product1", Description: "Description of Product1", Price: 11},
		{Name: "Product2", Description: "Description of Product2", Price: 22},
		{Name: "Product3", Description: "Description of Product3", Price: 33},
		{Name: "Product4", Description: "Description of Product4", Price: 44},
	}

	if err := td.DB.Create(&products).Error; err != nil {
		t.Fatalf("Failed to create test products: %v", err)
	}
}

func (td *DbTestData) Clear() {
	models := []any{
		&order.OrderItem{},
		&order.Order{},
		&user.User{},
		&product.Product{},
	}

	for _, m := range models {
		if err := td.DB.Unscoped().Where("1 = 1").Delete(m).Error; err != nil {
			fmt.Printf("Error clearing table: %v\n", err)
		}
	}
}
