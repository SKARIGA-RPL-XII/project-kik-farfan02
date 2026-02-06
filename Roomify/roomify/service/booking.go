package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
)

type BookingService struct {
	bookingRepo    repository.BookingRepository
	settingService SettingService
	userRepo       repository.UserRepository
	cfg            *config.Config
}

func NewBookingService(bookingrepo repository.BookingRepository, settingservice SettingService, userrepo repository.UserRepository, cfg *config.Config) *BookingService {
	return &BookingService{
		bookingRepo:    bookingrepo,
		settingService: settingservice,
		userRepo:       userrepo,
		cfg:            cfg,
	}
}

func (s *BookingService) CreateBooking(userID int, booking *models.CreateBookingRequest) (*models.Booking, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	startBook, err := time.ParseInLocation("2006-01-02 15:04", booking.StartBook, loc)
	if err != nil {
		return nil, errors.New("format waktu tidak valid")
	}

	endBook, err := time.ParseInLocation("2006-01-02 15:04", booking.EndBook, loc)
	if err != nil {
		return nil, errors.New("format waktu tidak valid")
	}


	startWork, endWork, err := s.settingService.GetWorkingHours()
	if err != nil {
		return nil, errors.New("jam kerja belum dikonfigurasi")
	}

	startHour, _ := time.Parse("15:04", startWork)
	endHour, _ := time.Parse("15:04", endWork)

	if startBook.Hour() < startHour.Hour() || startBook.Hour() > endHour.Hour() {
		return nil, errors.New("booking hanya bisa dibuat antara jam kerja")
	}

	isHoliday, err := s.settingService.IsHoliday(startBook)
	if err != nil {
		return nil, errors.New("gagal cek hari libur")
	}
	if isHoliday {
		return nil, errors.New("booking tidak bisa dibuat di hari libur")
	}

	conflict, err := s.bookingRepo.CheckBookingConflict(booking.IDDetail, startBook, endBook)
	if err != nil {
		return nil, errors.New("gagal cek bentrok booking")
	}
	if conflict {
		return nil, errors.New("ruangan sudah dibooking pada waktu tersebut")
	}

	newBooking := &models.Booking{
		Judul:          booking.Judul,
		Deskripsi:      booking.Deskripsi,
		IDLokasi:       booking.IDLokasi,
		IDDetail:       booking.IDDetail,
		StartBook:      startBook,
		EndBook:        endBook,
		Tanggal:        time.Date(startBook.Year(), startBook.Month(), startBook.Day(), 0, 0, 0, 0, startBook.Location()),
		DepartmentID:   user.DepartmentID,
		CreatedBy:      userID,
		Status:         "pending",
		IsRecurring:    booking.IsRecurring,
		RecurrenceRule: booking.RecurrenceRule,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if s.isAutoReminderEnabled() {
		newBooking.ReminderTime = startBook.Add(-1 * time.Hour)
	}

	err = s.bookingRepo.CreateBooking(newBooking)
	if err != nil {
		fmt.Printf("Error create booking di repo: %v\n", err)
		return nil, fmt.Errorf("gagal membuat booking: %w", err)
	}

	return newBooking, nil
}

func (s *BookingService) isAutoReminderEnabled() bool {
	value, err := s.settingService.GetSetting("auto_reminder")
	if err != nil {
		return false
	}
	return value == "true"
}

func (s *BookingService) GetBookingByID(id int) (*models.Booking, error) {
	return s.bookingRepo.GetBookingByID(id)
}

func (s *BookingService) GetBookingsByUser(userID int) ([]models.Booking, error) {
	return s.bookingRepo.GetBookingsByUser(userID)
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	return s.bookingRepo.GetAllBookings()
}

func (s *BookingService) UpdateBooking(userID, bookingID int, booking *models.CreateBookingRequest) error {
	existingBooking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	if existingBooking.CreatedBy != userID {
		return errors.New("anda tidak memiliki akses untuk update booking ini")
	}

	startBook, err := time.Parse("2006-01-02 15:04", booking.StartBook)
	if err != nil {
		return errors.New("format waktu tidak valid")
	}

	endBook, err := time.Parse("2006-01-02 15:04", booking.EndBook)
	if err != nil {
		return errors.New("format waktu tidak valid")
	}

	startWork, endWork, err := s.settingService.GetWorkingHours()
	if err != nil {
		return errors.New("jam kerja belum dikonfigurasi")
	}

	startHour, _ := time.Parse("15:04", startWork)
	endHour, _ := time.Parse("15:04", endWork)

	if startBook.Hour() < startHour.Hour() || startBook.Hour() > endHour.Hour() {
		return errors.New("booking hanya bisa dibuat antara jam kerja")
	}

	isHoliday, err := s.settingService.IsHoliday(startBook)
	if err != nil {
		return errors.New("gagal cek hari libur")
	}
	if isHoliday {
		return errors.New("booking tidak bisa dibuat di hari libur")
	}

	conflict, err := s.bookingRepo.CheckBookingConflict(booking.IDDetail, startBook, endBook)
	if err != nil {
		return errors.New("gagal cek bentrok booking")
	}
	if conflict && existingBooking.IDDetail != booking.IDDetail {
		return errors.New("ruangan sudah dibooking pada waktu tersebut")
	}

	updatedBooking := &models.Booking{
		Judul:          booking.Judul,
		IDLokasi:       booking.IDLokasi,
		IDDetail:       booking.IDDetail,
		StartBook:      startBook,
		EndBook:        endBook,
		IsRecurring:    booking.IsRecurring,
		RecurrenceRule: booking.RecurrenceRule,
	}

	err = s.bookingRepo.UpdateBooking(bookingID, updatedBooking)
	if err != nil {
		return errors.New("gagal update booking")
	}

	return nil
}

func (s *BookingService) DeleteBooking(userID, bookingID int) error {
	existingBooking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	if existingBooking.CreatedBy != userID {
		return errors.New("anda tidak memiliki akses untuk hapus booking ini")
	}

	err = s.bookingRepo.DeleteBooking(bookingID)
	if err != nil {
		return errors.New("gagal hapus booking")
	}

	return nil
}
