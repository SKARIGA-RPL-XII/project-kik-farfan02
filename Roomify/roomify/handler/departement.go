package handler

import (
	
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
)

type DeptHandler struct {
	deptService service.DeptService
}

func NewDeptHandler(deptService service.DeptService) *DeptHandler {
	return &DeptHandler{deptService: deptService}
}

func (h *DeptHandler) InputDepartment(c *fiber.Ctx) error {
	dpt := &models.Departement{}

	if err := c.BodyParser(dpt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if dpt.Nama_dtm == "" || dpt.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "semua field wajib diisi",
			Error:   "validation failed",
		})
	}

	err := h.deptService.InputDepartment(dpt)
	if err != nil {
		if err.Error() == "department sudah terdaftar" {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusConflict,
				Message: "department sudah terdaftar",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal membuat department",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusCreated,
		Message: "berhasil menambahkan department",
		Data: fiber.Map{
			"id":       dpt.ID,
			"nama_dtm": dpt.Nama_dtm,
			"code":     dpt.Code,
		},
	})
}

func (h *DeptHandler) GetAllDepartment(c *fiber.Ctx) error {
	departments, err := h.deptService.GetAllDepartment()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil data department",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil data department",
		Data:    departments,
	})
}

func (h *DeptHandler) UpdateDepartment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0{
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID tidak valid",
			Error:   err.Error(),
		})
	}

	type UpdateRequest struct {
		Nama_dtm string `json:"nama_dtm"`
		Code     string `json:"code"`
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

		if load.Nama_dtm == "" || load.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "semua field wajib diisi",
			Error:   "validation failed",
		})
	}

	
	
	dpt := &models.Departement{
		ID: id,
		Nama_dtm: load.Nama_dtm,
		Code: load.Code,
	}

	err = h.deptService.UpdateDepartment(id, dpt)
	if err != nil {
		if err.Error() == "department tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "department tidak ditemukan",
				Error:   err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal update department",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update department",
		Data: fiber.Map{
			"id":       id,
			"nama_dtm": dpt.Nama_dtm,
			"code":     dpt.Code,
		},
	})
}


func (h *DeptHandler) DeleteDepartment(c *fiber.Ctx) error {
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

	err = h.deptService.DeleteDepartment(id.ID)
	if err != nil {
		if err.Error() == "department tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "department tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal hapus department",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil hapus department",
		Data: "data berhasil terhapus",
	})
}