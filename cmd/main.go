package main

import (
	"log"


	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/container"

	"github.com/abdallahelassal/UserAuth/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := bootstrap.LoadConfig()
	applogger := logger.NewLogger(cfg.Enviroment)

	conn := bootstrap.NewConnection(cfg,applogger)

	if err := conn.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}


	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	c := container.NewContainer(conn.DB,applogger,*cfg)
	c.Handler.SetupRoutes()
	
	

	if err := c.Router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
	
	
}