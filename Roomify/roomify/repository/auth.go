package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (a *AuthRepository) Login(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, department_id, created_at 
		FROM users 
		WHERE email = $1
	`

	var password string
	var departmentID sql.NullInt64
	var createdAt sql.NullTime

	err := a.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&password,
		&user.Role,
		&departmentID,
		&createdAt,
	)

	user.Pass = strings.Trim(strings.TrimSpace(password), `"`)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("email tidak ditemukan")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	user.Pass = password
	
	if departmentID.Valid {
		user.DepartmentID = int(departmentID.Int64)
	}
	
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}

	return user, nil
}

func (a *AuthRepository) ChangePassword(userID int, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`

	result, err := a.DB.Exec(query, newPassword, userID)
	if err != nil {
		return fmt.Errorf("gagal update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal cek rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

func (a *AuthRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, department_id, created_at 
		FROM users 
		WHERE id = $1
	`

	var password string
	var departmentID sql.NullInt64
	var createdAt sql.NullTime

	err := a.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&password,
		&user.Role,
		&departmentID,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	user.Pass = password

	if departmentID.Valid {
		user.DepartmentID = int(departmentID.Int64)
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}

	return user, nil
}

func (a *AuthRepository) Logout(userID int) error {
	fmt.Printf("User ID %d logged out\n", userID)
	return nil
}
