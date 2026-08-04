package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
)

type JWTCLAIMS struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func Authorize(secretKey string) gin.HandlerFunc {
	key := []byte(secretKey)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			utils.Log.Warn("The Token is empty ")
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &JWTCLAIMS{}, func(token *jwt.Token) (any, error) {
			return key, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			utils.Log.Warn("The Token is invalid")
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			utils.Log.Warn("The Token is invalid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*JWTCLAIMS)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			utils.Log.Warn("The Token is invalid")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func GetTheUserId(c *gin.Context) (string, bool) {
	userId, ok := c.Get("user_id")
	if !ok {
		utils.Log.Warn("User ID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return "", false
	}
	id, ok := userId.(string)
	if !ok {
		utils.Log.Warn("User ID is not a string")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return "", false
	}
	return id, true
}
