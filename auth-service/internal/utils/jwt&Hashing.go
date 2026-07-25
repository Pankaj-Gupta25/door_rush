package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTManager struct {
	SecretKey []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{SecretKey: []byte(secret)}
}

func (j *JWTManager) GenerateToken(userID string, email string, password string, username string) (string, string, error) {
	tokenClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"nbf":     jwt.NewNumericDate(time.Now()),
		"iat":     jwt.NewNumericDate(time.Now()),
		"iss":     "auth-service",
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"nbf":     jwt.NewNumericDate(time.Now()),
		"iat":     jwt.NewNumericDate(time.Now()),
		"iss":     "auth-service",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := accessToken.SignedString(j.SecretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(j.SecretKey)
	if err != nil {
		return "", "", err
	}
	return tokenString, refreshTokenString, nil
}

func (j *JWTManager) ValidateToken(token string) (bool, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (j *JWTManager) ValidateRefreshToken(token string) (bool, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (j *JWTManager) GetEmailFromToken(ctx *gin.Context) (string, error) {
	authToken := ctx.GetHeader("Authorization")
	if authToken == "" {
		return "", errors.New("unauthorized: no token provided")
	}
	if !strings.HasPrefix(authToken, "Bearer ") {
		return "", errors.New("unauthorized: invalid token format")
	}

	// Extract token string
	tokenString := strings.TrimPrefix(authToken, "Bearer ")

	// Parse token with claims
	userEmailClaims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, userEmailClaims, func(token *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", errors.New("unauthorized: invalid token")
	}

	// Extract email from claims
	email, ok := userEmailClaims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("unauthorized: email not found in token")
	}

	return email, nil
}
func (j *JWTManager) Authorize(ctx *gin.Context) (bool, error) {
	clientToken := ctx.GetHeader("Authorization")
	if clientToken == "" {
		return false, errors.New("Unauthorized")
	}
	if clientToken == "" {
		return false, errors.New("Unauthorized")
	}
	if !strings.HasPrefix(clientToken, "Bearer ") {
		return false, errors.New("Unauthorized")
	}
	clientToken = strings.TrimPrefix(clientToken, "Bearer ")
	isValid, err := j.ValidateToken(clientToken)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, errors.New("Unauthorized")
	}
	return true, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func CheckIfThePasswordIsCorrect(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
