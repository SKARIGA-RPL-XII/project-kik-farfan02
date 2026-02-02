package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Pass         string    `json:"password"`
	Role         string    `json:"role"`
	DepartmentID int       `json:"department_id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type UserFilter struct {
	Search       string `json:"search"`
	DepartmentID int    `json:"department_id"`
}
