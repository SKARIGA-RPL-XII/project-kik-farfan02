package util

import (
	"fmt"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

func ValidateTime(start, end string) error{
	startTime, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil{
		startTime, err = time.Parse("2006-01-02 15:04", start)	
		if err != nil{
			return fmt.Errorf("format jam awal tidak valid ")
		}
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", end)
	if err != nil{
		endTime, err = time.Parse("2006-01-02 15:04", end)
		if err != nil{
			return fmt.Errorf("format jam akhir tidak valid")
		}
	}

	if !endTime.After(startTime){
		return fmt.Errorf("jam akhir harus lebih besar dari jam awal ")
	}

	if startTime.Before(time.Now()){
		return fmt.Errorf("tidak bisa booking di waktu yang sudah lewat ")
	}
	return nil
}

func  CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}