package main

import (
	"log"

	"github.com/abdallahelassal/UserAuth.git/config"
	"github.com/abdallahelassal/UserAuth.git/internal/database"
	"github.com/abdallahelassal/UserAuth.git/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	applogger := logger.NewLogger(cfg.Enviroment)

	conn := database.NewConnection(cfg,applogger)

	if err := conn.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := database.RunMigrations(cfg.DatabaseConfig.DatabaseURL,"./migration",applogger); err != nil{
		if err := conn.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
	
	
}