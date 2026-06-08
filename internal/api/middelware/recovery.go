package middelware

import (
	
	"net/http"
	"runtime/debug"

	"github.com/abdallahelassal/UserAuth/pkg/logger"
	"github.com/gin-gonic/gin"
	
)


func Recovery()gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func ()  {
			
			if err := recover();err != nil {
				logger.Error("critical panic recover",
					logger.Any("error", err),
					logger.String("path", c.Request.URL.RawPath),
					logger.String("method",c.Request.Method),						
				)
			logger.Error("stack track:"+ string(debug.Stack()))
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error in server"})
			c.Abort()
			return 
		}	
		}()
		c.Next()
	}
}