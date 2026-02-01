package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Login(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, department_id, created_at 
		FROM users 
		WHERE email = $1
	`

	var password string
	var departmentID sql.NullInt64
	var createdAt sql.NullTime

	err := r.DB.QueryRow(query, email).Scan(
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

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email, password, role, department_id, created_at)
			VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	err := r.DB.QueryRow(query,
		user.Name,
		user.Email,
		user.Pass,
		user.Role,
		user.DepartmentID,
		time.Now(),
	).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("email sudah terdaftar")
		}
		return fmt.Errorf("gagal membuat user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, password, role, department_id, created_at FROM users WHERE email = $1`

	var password string
	var departmentID sql.NullInt64
	var createdAt sql.NullTime

	err := r.DB.QueryRow(query, email).Scan(
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
