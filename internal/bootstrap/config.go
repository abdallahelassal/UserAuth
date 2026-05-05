package bootstrap

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
	AccessTokenSecret   string
	RefreshTokenSecret	string
	AccessExpiration   	int
	RefreshExpiration 	int
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
			AccessTokenSecret: getEnv("ACCESS_TOKEN_SECRET",""),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET",""),
			AccessExpiration: getEnvInt("ACCESS_TOKEN_EXPIRY_HOUR",60),
			RefreshExpiration: getEnvInt("REFRESH_TOKEN_EXPIRY_HOUR",9),
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