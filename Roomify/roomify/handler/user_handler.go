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


func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	type SearchRequest struct {
		Search       string `json:"search"`
		DepartmentID int    `json:"department_id"`
	}

	load := new(SearchRequest)

	if err := c.BodyParser(load); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	filter := models.UserFilter{
		Search:       load.Search,
		DepartmentID: load.DepartmentID,
	}

	users, err := h.userService.GetUsers(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil data user",
			Error:   err.Error(),
		})
	}

	for i := range users {
		users[i].Pass = ""
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil data user",
		Data:    users,
	})
}

func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
	}

	type UpdateRequest struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		DepartmentID int    `json:"department_id"`
	}

	load := new(UpdateRequest)

	if err := c.BodyParser(load); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if load.Name == "" || load.Email == "" || load.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "semua field wajib diisi",
			Error:   "validation failed",
		})
	}

	userData := &models.User{
		Name:         load.Name,
		Email:        load.Email,
		Role:         load.Role,
		DepartmentID: load.DepartmentID,
	}

	err = h.userService.UpdateUser(id, userData)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "user tidak ditemukan",
				Error:   err.Error(),
			})
		}
		if err.Error() == "email sudah digunakan" {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusConflict,
				Message: "email sudah digunakan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal update user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update user",
		Data: fiber.Map{
			"id":             id,
			"name":           userData.Name,
			"email":          userData.Email,
			"role":           userData.Role,
			"department_id":  userData.DepartmentID,
		},
	})
}

func (h *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	id := &models.RequestDelete{} 
	err := c.BodyParser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
	}

	err = h.userService.DeleteUser(id.ID)
	if err != nil {
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
			Message: "gagal hapus user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil hapus user",
		Data: fiber.Map{
			"id": id,
		},
	})
}