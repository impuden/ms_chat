package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUsername string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	JWTToken   string
	JWTSecret  string
	JWTSalt    string
}

var GlobalConfig Config

func LoadConfig() *Config {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}

	// Инициализируем конфигурацию из переменных окружения
	GlobalConfig = Config{
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTToken:   os.Getenv("JWT_TOKEN"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JWTSalt:    os.Getenv("JWT_SALT"),
	}

	// Проверка обязательных полей
	requiredFields := map[string]string{
		"DB_USERNAME": GlobalConfig.DBUsername,
		"DB_PASSWORD": GlobalConfig.DBPassword,
		"DB_NAME":     GlobalConfig.DBName,
		"DB_HOST":     GlobalConfig.DBHost,
		"DB_PORT":     GlobalConfig.DBPort,
		"JWT_TOKEN":   GlobalConfig.JWTToken,
		"JWT_SECRET":  GlobalConfig.JWTSecret,
		"JWT_SALT":    GlobalConfig.JWTSalt,
	}
	for key, value := range requiredFields {
		if value == "" {
			log.Fatalf("Ошибка: %s не установлено", key)
		}
	}

	return &GlobalConfig
}
