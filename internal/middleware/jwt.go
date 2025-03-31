package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lutestringamend/perwebbe/internal/config"
)

func JWTAuthMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func GenerateJWT(userID uint, username string, role string, jwtConfig config.JWTConfig) (string, error) {
	claims := jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),
		"username": username,
		"role":     role,
		"iss":      jwtConfig.Issuer,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(jwtConfig.TokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID uint, jwtConfig config.JWTConfig) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"iss": jwtConfig.Issuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(jwtConfig.RefreshExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractUserIDFromToken(c *gin.Context) (uint, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return 0, errors.New("token claims not found")
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims format")
	}

	sub, ok := mapClaims["sub"].(string)
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	var userID uint
	_, err := fmt.Sscanf(sub, "%d", &userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
