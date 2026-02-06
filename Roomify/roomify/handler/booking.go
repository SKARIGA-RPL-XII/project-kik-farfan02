package handler

import (
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("id").(int)

	booking := new(models.CreateBookingRequest)
	if err := c.BodyParser(booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	// Validasi input
	if booking.Judul == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "judul wajib diisi",
			Error:   "validation failed",
		})
	}

	// Buat booking
	newBooking, err := h.bookingService.CreateBooking(userID, booking)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "gagal membuat booking",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusCreated,
		Message: "berhasil membuat booking",
		Data:    newBooking,
	})
}


func (h *BookingHandler) GetBookingByID(c *fiber.Ctx) error {
	type Request struct {
		ID int `json:"id"`
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

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID booking tidak valid",
			Error:   "validation failed",
		})
	}

	booking, err := h.bookingService.GetBookingByID(req.ID)
	if err != nil {
		if err.Error() == "booking tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "booking tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil booking",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil booking",
		Data:    booking,
	})
}

func (h *BookingHandler) GetBookingsByUser(c *fiber.Ctx) error {
	userID := c.Locals("id").(int)

	bookings, err := h.bookingService.GetBookingsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusInternalServerError,
			Message: "gagal mengambil booking",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil mengambil booking",
		Data:    bookings,
	})
}

func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
    bookings, err := h.bookingService.GetAllBookings()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Success: false,
            Code:    fiber.StatusInternalServerError,
            Message: "gagal mengambil data booking",
            Error:   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
        Success: true,
        Code:    fiber.StatusOK,
        Message: "berhasil mengambil semua booking",
        Data:    bookings,
    })
}


func (h *BookingHandler) UpdateBooking(c *fiber.Ctx) error {
	userID := c.Locals("id").(int)

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID booking tidak valid",
			Error:   err.Error(),
		})
	}

	booking := new(models.CreateBookingRequest)
	if err := c.BodyParser(booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "invalid format JSON",
			Error:   err.Error(),
		})
	}

	err = h.bookingService.UpdateBooking(userID, id, booking)
	if err != nil {
		if err.Error() == "booking tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "booking tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "gagal update booking",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil update booking",
		Data: fiber.Map{
			"id": id,
		},
	})
}

func (h *BookingHandler) DeleteBooking(c *fiber.Ctx) error {
	userID := c.Locals("id").(int)

	type Request struct {
		ID int `json:"id"`
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

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "ID booking tidak valid",
			Error:   "validation failed",
		})
	}

	err := h.bookingService.DeleteBooking(userID, req.ID)
	if err != nil {
		if err.Error() == "booking tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Success: false,
				Code:    fiber.StatusNotFound,
				Message: "booking tidak ditemukan",
				Error:   err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Code:    fiber.StatusBadRequest,
			Message: "gagal hapus booking",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: "berhasil hapus booking",
		Data: fiber.Map{
			"id": req.ID,
		},
	})
}