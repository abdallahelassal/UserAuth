package middelware

import (
	"net/http"
	"strings"

	"github.com/abdallahelassal/UserAuth/pkg/jwt"
	"github.com/gin-gonic/gin"
)



func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		authorized, err := jwt.IsAuthorized(token, secret)
		if err != nil || !authorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid token"})
			c.Abort()
			return
		}

		userID, err := jwt.ExtractIDFromToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid token payload"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}


