package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
)

type JWTClaims struct {
	Email  string `json:"email"`
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func Authenticate(secretKey string) gin.HandlerFunc {
	key := []byte(secretKey)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.Log.Error("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
			return key, nil
		})
		if err != nil {
			utils.Log.Error("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !token.Valid {
			utils.Log.Error("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			utils.Log.Error("Unauthorized")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token from claims"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserId)

		c.Next()
	}
}
func HasUserId(c *gin.Context) (string, bool) {
	userID, exist := c.Get("user_id")
	if !exist {
		utils.Log.Warn("Unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return "", false
	}
	id, ok := userID.(string)
	if !ok {
		utils.Log.Warn("Unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return "", false
	}
	return id, ok && id != ""
}
