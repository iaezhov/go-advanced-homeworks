package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db      DbConfig
	AppPort string
}
type DbConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			DBHost:     getEnvAsStr("DB_HOST", "localhost"),
			DBPort:     getEnvAsInt("DB_PORT", 5432),
			DBUser:     getEnvAsStr("DB_USER", ""),
			DBPassword: getEnvAsStr("DB_PASSWORD", ""),
			DBName:     getEnvAsStr("DB_NAME", ""),
			DBSSLMode:  getEnvAsStr("DB_SSLMODE", "disable"),
		},
		AppPort: getEnvAsStr("APP_PORT", "8081"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Db.DBHost,
		c.Db.DBUser,
		c.Db.DBPassword,
		c.Db.DBName,
		c.Db.DBPort,
		c.Db.DBSSLMode,
	)
}

func getEnvAsStr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: invalid int for %s: %s, using default %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
