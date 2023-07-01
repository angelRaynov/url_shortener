package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"url_shortener/infrastructure/config"
)

func AuthMiddleware(cfg *config.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := bearerToken[1]

		// Verify and parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract the username claim from the token's payload
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Store the username in the Gin context
			username := claims["username"].(string)
			c.Set("username", username)

			uid := claims["uid"].(string)
			c.Set("owner_uid", uid)

			email := claims["email"].(string)
			c.Set("email", email)
		}

		// Token is valid, proceed with the next handler
		c.Next()
	}
}
