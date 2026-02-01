package handler

import (
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"

	"github.com/gofiber/fiber/v2"
)




func (h *UserHandler) Login(c *fiber.Ctx) error {
	creds := &models.Credentials{}

	if err := c.BodyParser(creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if creds.Email == "" || creds.Pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "email dan password wajib diisi",
			Error:   "validation failed",
		})
	}

	response, err := h.userService.Login(creds)
	if err != nil {
		if err.Error() == "email tidak ditemukan" || err.Error() == "password salah" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "kredensial tidak valid",
				Error:   err.Error(),
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "terjadi kesalahan pada server",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil login",
		Data:    response,
	})
}

