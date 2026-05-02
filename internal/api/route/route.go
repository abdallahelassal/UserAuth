package route

import (
	"time"

	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func SetupRoutes(cfg *bootstrap.Config, timeout *time.Duration, db *gorm.DB , gin *gin.Engine) {
	puplicRouter := gin.Group("")
	// puplic router 
	

}