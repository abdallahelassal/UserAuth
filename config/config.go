package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig DatabaseConfig
	ServerConfig   ServerConfig
	JWTConfig      JWTConfig
	Enviroment     string
}

type DatabaseConfig struct {
	Name     	string
	User     	string
	Password 	string
	Port     	string
	Host     	string
	DatabaseURL string
	SSLMode  	string
}

type ServerConfig struct {
	Port string
	Host string
}

type JWTConfig struct {
	Secret            string
	AccessExpiration   int
	RefreshExpiration int
}

func LoadConfig() *Config {
	
	var loadedFrom string 
	condidates := []string{".env","../.env","../../.env","../../../.env"}

	for _, v := range condidates{
		if err := godotenv.Load(v); err != nil {
			loadedFrom = v
			break
		}
	}

	if loadedFrom != "" {
		log.Println("loaded .env from",loadedFrom)
	}else{
		log.Printf("no .env file found (looked in .env and parent directories)")
	}

	return &Config{
		ServerConfig: ServerConfig{
			Port: getEnv("PORT","8000"),
			Host: getEnv("HOST","0.0.0.0"),
		},
		DatabaseConfig: DatabaseConfig{
			Name: getEnv("DB_NAME","users"),
			User: getEnv("DB_USER","postgres"),
			Password: getEnv("DB_PASSWORD",""),
			Port: getEnv("DB_PORT","5432"),
			SSLMode: getEnv("DB_SSL_MODE","disable"),
			Host: getEnv("DB_HOST","localhost"),
			DatabaseURL: getEnv("DATABASE_URL", ""),
		},
		JWTConfig: JWTConfig{
			Secret: getEnv("SECRET_KEY",""),
			AccessExpiration: getEnvInt("JWT_ACCESS_EXPIRATION",60),
			RefreshExpiration: getEnvInt("JWT_REFRESH_EXPIRATION", 9),
		},
		Enviroment: getEnv("ENVIROMENT", "DEVELOPMENT"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue
}
 
func getEnvInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value , err := strconv.Atoi(valueStr); err == nil{
		return value
	}
	return defaultValue
}