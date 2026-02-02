package handler

import (
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"


	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authsService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authsService,
	}
}



func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	response, err := h.authService.Login(creds)
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

func (a *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		PasswordLama      string `json:"password_lama"`
		PasswordBaru      string `json:"password_baru"`
		PasswordBaruUlang string `json:"password_baru_ulang"`
	}

	load := new(ChangePasswordRequest)

	if err := c.BodyParser(load); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if load.PasswordLama == "" || load.PasswordBaru == "" || load.PasswordBaruUlang == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "semua field wajib diisi",
			Error:   "validation failed",
		})
	}

	if load.PasswordBaru != load.PasswordBaruUlang {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "password baru dan ulang tidak sama",
			Error:   "validation failed",
		})
	}

	userID := c.Locals("id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "user tidak ditemukan",
		})
	}

	err := a.authService.ChangePassword(userID.(int), load.PasswordLama, load.PasswordBaru, load.PasswordBaruUlang)
	if err != nil {
		if err.Error() == "password lama salah" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusUnauthorized,
				Message: "password lama salah",
				Error:   err.Error(),
			})
		}
		if err.Error() == "password baru dan konfirmasi tidak sama" || err.Error() == "password minimal 6 karakter" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusBadRequest,
				Message: "validasi gagal",
				Error:   err.Error(),
			})
		}
		if err.Error() == "user tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "user tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal ganti password",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil ganti password",
		Data:    nil,
	})
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "user tidak ditemukan",
		})
	}

	if err := a.authService.Logout(userID.(int)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal logout",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil logout",
		Data:    nil,
	})
}

