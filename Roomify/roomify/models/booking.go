package models

import "time"

type Booking struct {
	ID           int       `json:"id"`
	Judul        string    `json:"judul"`
	Deskripsi    string    `json:"deskripsi"`
	IDLokasi     int       `json:"id_lokasi"`
	IDDetail     int       `json:"id_detail"`
	StartBook    time.Time `json:"start_book"`
	EndBook      time.Time `json:"end_book"`
	Tanggal      time.Time `json:"tanggal"`
	DepartmentID int       `json:"department_id"`
	CreatedBy    int       `json:"created_by"`
	Status       string    `json:"status"`
	ReminderTime time.Time `json:"reminder_time"`
	ReminderSent bool      `json:"reminder_sent"` 
	IsRecurring    bool      `json:"is_recurring"`
	RecurrenceRule string    `json:"recurrence_rule"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateBookingRequest struct {
	Judul          string `json:"judul"`
	Deskripsi      string `json:"deskripsi"`
	IDLokasi       int    `json:"id_lokasi"`
	IDDetail       int    `json:"id_detail"`
	StartBook      string `json:"start_book"`
	EndBook        string `json:"end_book"`
	IsRecurring    bool   `json:"is_recurring"`
	RecurrenceRule string `json:"recurrence_rule"`
}
