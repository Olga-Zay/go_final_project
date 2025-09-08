package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port string
	DB   string
	Pass string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env, используем переменные окружения")
	}

	return &Config{
		Port: getEnv("TODO_PORT", "8080"),
		DB:   getEnv("TODO_DBFILE", "scheduler.db"),
		Pass: getEnv("TODO_PASSWORD", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
