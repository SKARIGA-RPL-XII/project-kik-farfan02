package util

import (
	"fmt"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID int, role string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTKey))

	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}

	return tokenString, nil
}

func GetUserIDFromContext(c *fiber.Ctx) int {
	userID := c.Locals("id")
	if userID == nil {
		return 0
	}
	return userID.(int)
}

func GetUserRoleFromContext(c *fiber.Ctx) string {
	userRole := c.Locals("role")
	if userRole == nil {
		return ""
	}
	return userRole.(string)
}