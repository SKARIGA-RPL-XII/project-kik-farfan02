package handler

import (

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
)

type LokasiHandler struct {
	lokasiService service.LokasiService
}

func NewLokasiHandler(lokasiService service.LokasiService) *LokasiHandler {
	return &LokasiHandler{
		lokasiService: lokasiService,
	}
}

func (h *LokasiHandler) CreateLokasi(c *fiber.Ctx) error {
	lok := new(models.CreateLocationRequest)

	if err := c.BodyParser(lok); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
	Error:   err.Error(),
		})
	}

	if lok.NamaLokasi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "nama lokasi wajib diisi",
			Error:   "validation failed",
		})
	}

	if lok.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "kapasitas lokasi harus lebih dari 0",
			Error:   "validation failed",
		})
	}

	if len(lok.Ruangan) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "setidaknya harus ada 1 ruangan",
			Error:   "validation failed",
		})
	}

	lokID, err := h.lokasiService.CreateLokasi(lok)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal membuat lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusCreated,
		Message: "berhasil membuat lokasi",
		Data: fiber.Map{
			"id":         lokID,
			"nama_lokasi": lok.NamaLokasi,
			"capacity":   lok.Capacity,
			"ruangan":    lok.Ruangan,
		},
	})
}

func (h *LokasiHandler) GetAllLocations(c *fiber.Ctx) error {
	locations, err := h.lokasiService.GetAllLocations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil data lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil data lokasi",
		Data:    locations,
	})
}

func (h *LokasiHandler) GetLocationDetails(c *fiber.Ctx) error {
	lokID, err := c.ParamsInt("id")
	if err != nil || lokID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID lokasi tidak valid",
			Error:   err.Error(),
		})
	}


	details, err := h.lokasiService.GetLocationDetails(lokID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil detail lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil detail lokasi",
		Data:    details,
	})
}

func (h *LokasiHandler) GetAllLocationsWithDetails(c *fiber.Ctx) error {
	locations, err := h.lokasiService.GetAllLocationsWithDetails()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil data lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil data lokasi",
		Data:    locations,
	})
}

func (h *LokasiHandler) UpdateLokasi(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID lokasi tidak valid",
			Error:   err.Error(),
		})
	}

	lok := new(models.CreateLocationRequest)

	if err := c.BodyParser(lok); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if lok.NamaLokasi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "nama lokasi wajib diisi",
			Error:   "validation failed",
		})
	}

	if lok.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "kapasitas lokasi harus lebih dari 0",
			Error:   "validation failed",
		})
	}

	if len(lok.Ruangan) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "setidaknya harus ada 1 ruangan",
			Error:   "validation failed",
		})
	}

	err = h.lokasiService.UpdateLokasi(id, lok)
	if err != nil {
		if err.Error() == "lokasi tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "lokasi tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal update lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update lokasi",
		Data: fiber.Map{
			"id":          id,
			"nama_lokasi": lok.NamaLokasi,
			"capacity":    lok.Capacity,
			"ruangan":     lok.Ruangan,
		},
	})
}

func (h *LokasiHandler) DeleteLokasi(c *fiber.Ctx) error {
	type DeleteRequest struct {
		ID int `json:"id"`
	}

	req := new(DeleteRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID lokasi tidak valid",
			Error:   "validation failed",
		})
	}

	err := h.lokasiService.DeleteLokasi(req.ID)
	if err != nil {
		if err.Error() == "lokasi tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "lokasi tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal hapus lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil hapus lokasi",
		Data: fiber.Map{
			"id": req.ID,
		},
	})
}

func (h *LokasiHandler) UpdateDetailLokasi(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID detail tidak valid",
			Error:   err.Error(),
		})
	}

	detail := new(models.DetailLocation)

	if err := c.BodyParser(detail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if detail.NamaRuangan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "nama ruangan wajib diisi",
			Error:   "validation failed",
		})
	}

	if detail.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "kapasitas ruangan harus lebih dari 0",
			Error:   "validation failed",
		})
	}

	err = h.lokasiService.UpdateDetailLokasi(id, detail)
	if err != nil {
		if err.Error() == "detail lokasi tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "detail lokasi tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal update detail lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update detail lokasi",
		Data: fiber.Map{
			"id":           id,
			"nama_ruangan": detail.NamaRuangan,
			"capacity":     detail.Capacity,
		},
	})
}

// DeleteDetailLokasi - Delete detail lokasi
func (h *LokasiHandler) DeleteDetailLokasi(c *fiber.Ctx) error {
	type DeleteRequest struct {
		ID int `json:"id"`
	}

	req := new(DeleteRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID detail tidak valid",
			Error:   "validation failed",
		})
	}

	err := h.lokasiService.DeleteDetailLokasi(req.ID)
	if err != nil {
		if err.Error() == "detail lokasi tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "detail lokasi tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal hapus detail lokasi",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil hapus detail lokasi",
		Data: fiber.Map{
			"id": req.ID,
		},
	})
}