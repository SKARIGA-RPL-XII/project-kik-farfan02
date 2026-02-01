package handler

import (
	"strings"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if user.Name == "" || user.Email == "" || user.Pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "semua field wajib diisi",
			Error:   "validation failed",
		})
	}

	if !strings.Contains(user.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "email tidak valid",
			Error:   "email harus mengandung @",
		})
	}

	
	if err := h.userService.CreateUser(user); err != nil {
		if err.Error() == "email sudah terdaftar" {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusConflict,
				Message: "email sudah terdaftar",
				Error:   err.Error(),
			})
		}
		if err.Error() == "email tidak valid" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusBadRequest,
				Message: "email tidak valid",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal membuat user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusCreated,
		Message: "berhasil menambahkan user",
		Data: fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}