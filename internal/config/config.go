package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Kafka  KafkaConfig
	Server ServerConfig
	Cache  CacheConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type ServerConfig struct {
	Port string
}

type CacheConfig struct {
	Size int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	var err error

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port, err = getEnvAsInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "")
	cfg.DB.Name = getEnv("DB_NAME", "orders")

	brokers := getEnv("KAFKA_BROKERS", "localgost:9092")
	cfg.Kafka.Brokers = strings.Split(brokers, ",")
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC", "orders")
	cfg.Kafka.GroupID = getEnv("KAFKA_GROUP", "order_service")

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	cfg.Cache.Size, err = getEnvAsInt("CACHE_SIZE", 100)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	if val := os.Getenv(key); val != "" {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return valInt, nil
	}
	return defaultValue, nil
}
