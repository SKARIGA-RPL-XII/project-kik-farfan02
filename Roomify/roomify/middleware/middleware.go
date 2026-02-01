package middleware

import (
	"fmt"
	"strings"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware - Middleware untuk validasi token JWT
func AuthMiddleware(cfg *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token tidak ditemukan",
			})
		}

		// Hapus prefix "Bearer "
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
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

		// Validasi token dan ambil claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Simpan user info ke locals untuk digunakan di handler
			c.Locals("id", int(claims["id"].(float64)))
			c.Locals("role", claims["role"].(string))

			// Lanjut ke handler berikutnya
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

// AdminMiddleware - Middleware untuk cek role admin
func AdminMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Ambil role dari locals (sudah diset oleh AuthMiddleware)
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "User tidak ditemukan",
			})
		}

		// Cek apakah role admin
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

// UserMiddleware - Middleware untuk cek role user (optional)
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