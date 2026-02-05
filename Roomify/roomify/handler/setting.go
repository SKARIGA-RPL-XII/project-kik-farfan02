package handler

import (
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
	settingService service.SettingService
}

func NewSettingHandler(settingService service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

func (h *SettingHandler) GetSetting(c *fiber.Ctx) error {
	type Request struct {
		Key string `json:"key"`
	}

	req := new(Request)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if req.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "key is required",
			Error:   "validation failed",
		})
	}

	value, err := h.settingService.GetSetting(req.Key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusNotFound,
			Message: "setting not found",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil setting",
		Data: fiber.Map{
			"key":   req.Key,
			"value": value,
		},
	})
}

func (h *SettingHandler) GetAllSettings(c *fiber.Ctx) error {
	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil settings",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil semua settings",
		Data:    settings,
	})
}

func (h *SettingHandler) UpdateSetting(c *fiber.Ctx) error {
	req := new(models.SettingRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if req.Key == "" || req.Value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "key and value are required",
			Error:   "validation failed",
		})
	}

	err := h.settingService.UpdateSetting(req.Key, req.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal update setting",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update setting",
		Data: fiber.Map{
			"key":   req.Key,
			"value": req.Value,
		},
	})
}