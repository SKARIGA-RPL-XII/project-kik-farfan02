package middleware

import (
	"fmt"
	"strings"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token tidak ditemukan",
			})
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, fmt.Errorf("invalid signing method")
			}

			return []byte(cfg.JWTKey), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token tidak valid",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("id", int(claims["id"].(float64)))
			c.Locals("role", claims["role"].(string))

			c.Next()
			return nil
		}

		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Token tidak valid",
		})
	}
}

func AdminMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "User tidak ditemukan",
			})
		}

		if userRole.(string) != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusForbidden,
				Message: "Forbidden",
				Error:   "Hanya Admin yang bisa akses",
			})
		}

		c.Next()
		return nil
	}
}

func UserMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "User tidak ditemukan",
			})
		}

		// Cek apakah role user atau admin
		role := userRole.(string)
		if role != "user" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusForbidden,
				Message: "Forbidden",
				Error:   "Akses ditolak",
			})
		}

		c.Next()
		return nil
	}
}