package middlewares

import (
	"net/http"
	"strings"
	"wordCraft/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userId, err := utils.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("userId", userId)
	c.Next()
}
