package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type BookingRepository struct {
	DB *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}


func (r *BookingRepository) CreateBooking(booking *models.Booking) error {
	query := `
		INSERT INTO booking (
			judul, deskripsi, id_lokasi, id_detail, start_book, end_book, tanggal,
			department_id, created_by, status,
			reminder_time, reminder_sent, is_recurring, recurrence_rule
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.DB.Exec(query,
		booking.Judul,
		booking.Deskripsi,
		booking.IDLokasi,
		booking.IDDetail,
		booking.StartBook,
		booking.EndBook,
		booking.Tanggal,
		booking.DepartmentID,
		booking.CreatedBy,
		booking.Status,
		booking.ReminderTime,
		false,
		booking.IsRecurring,
		booking.RecurrenceRule,
	)
	return err
}

func (r *BookingRepository) CheckBookingConflict(idDetail int, startBook, endBook time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM booking 
			WHERE id_detail = $1 
			AND status IN ('pending', 'approved')  
			AND start_book < $3 
			AND end_book > $2
		)
	`
	var exists bool
	err := r.DB.QueryRow(query, idDetail, startBook, endBook).Scan(&exists)
	return exists, err
}

func (r *BookingRepository) GetBookingByID(id int) (*models.Booking, error) {
	query := `
		SELECT id, judul, deskripsi, id_lokasi, id_detail, start_book, end_book,
		       department_id, created_by, status, reminder_time,
		       is_recurring, recurrence_rule, created_at, updated_at
		FROM booking WHERE id = $1
	`

	booking := &models.Booking{}
	err := r.DB.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.Judul,
		&booking.Deskripsi,
		&booking.IDLokasi,
		&booking.IDDetail,
		&booking.StartBook,
		&booking.EndBook,
		&booking.DepartmentID,
		&booking.CreatedBy,
		&booking.Status,
		&booking.ReminderTime,
		&booking.IsRecurring,
		&booking.RecurrenceRule,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("booking tidak ditemukan")
		}
		return nil, err
	}

	return booking, nil
}

func (r *BookingRepository) GetBookingsByUser(userID int) ([]models.Booking, error) {
	query := `
		SELECT id, judul, deskripsi, id_lokasi, id_detail, start_book, end_book,
		       department_id, created_by, status, reminder_time,
		       is_recurring, recurrence_rule, created_at, updated_at
		FROM booking 
		WHERE created_by = $1 
		ORDER BY start_book DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		booking := models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.Judul,
			&booking.Deskripsi,
			&booking.IDLokasi,
			&booking.IDDetail,
			&booking.StartBook,
			&booking.EndBook,
			&booking.DepartmentID,
			&booking.CreatedBy,
			&booking.Status,
			&booking.ReminderTime,
			&booking.IsRecurring,
			&booking.RecurrenceRule,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// GetAllBookings - Get all bookings (for admin)
func (r *BookingRepository) GetAllBookings() ([]models.Booking, error) {
    query := `
        SELECT id, judul, deskripsi, id_lokasi, id_detail, start_book, end_book,
               department_id, created_by, status, reminder_time,
               is_recurring, recurrence_rule, created_at, updated_at
        FROM booking 
        ORDER BY start_book DESC
    `

    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bookings []models.Booking
    for rows.Next() {
        booking := models.Booking{}
        err := rows.Scan(
            &booking.ID,
            &booking.Judul,
            &booking.Deskripsi,
            &booking.IDLokasi,
            &booking.IDDetail,
            &booking.StartBook,
            &booking.EndBook,
            &booking.DepartmentID,
            &booking.CreatedBy,
            &booking.Status,
            &booking.ReminderTime,
            &booking.IsRecurring,
            &booking.RecurrenceRule,
            &booking.CreatedAt,
            &booking.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        bookings = append(bookings, booking)
    }

    return bookings, nil
}

func (r *BookingRepository) UpdateBooking(id int, booking *models.Booking) error {
	query := `
		UPDATE booking 
		SET judul = $1, deskripsi = $2, id_lokasi = $3, id_detail = $4, 
		    start_book = $5, end_book = $6, status = $7,
		    is_recurring = $8, recurrence_rule = $9,
		    updated_at = NOW()
		WHERE id = $10
	`

	result, err := r.DB.Exec(query,
		booking.Judul,
		booking.Deskripsi,
		booking.IDLokasi,
		booking.IDDetail,
		booking.StartBook,
		booking.EndBook,
		booking.Status,
		booking.IsRecurring,
		booking.RecurrenceRule,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("booking tidak ditemukan")
	}

	return nil
}

func (r *BookingRepository) DeleteBooking(id int) error {
	query := `DELETE FROM booking WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("booking tidak ditemukan")
	}

	return nil
}