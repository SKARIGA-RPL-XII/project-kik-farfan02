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

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, department_id, created_at 
		FROM users 
		WHERE id = $1
	`

	var password string
	var departmentID sql.NullInt64
	var createdAt sql.NullTime

	err := r.DB.QueryRow(query, id).Scan(
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

// GetUsers - Ambil semua user dengan fitur search
func (r *UserRepository) GetUsers(filter models.UserFilter) ([]models.User, error) {
	query := `
		SELECT id, name, email, password, role, department_id, created_at 
		FROM users 
	`

	var args []interface{}
	var conditions []string
	argCount := 1

	// Search by name or email
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name LIKE $%d OR email LIKE $%d", argCount, argCount))
		args = append(args, "%"+filter.Search+"%")
		argCount++
	}

	// Filter by department_id
	if filter.DepartmentID > 0 {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argCount))
		args = append(args, filter.DepartmentID)
		argCount++
	}

	// Combine conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY id ASC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var password string
		var departmentID sql.NullInt64
		var createdAt sql.NullTime

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&password,
			&user.Role,
			&departmentID,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan users: %w", err)
		}

		user.Pass = password

		if departmentID.Valid {
			user.DepartmentID = int(departmentID.Int64)
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(id int, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, role = $3, department_id = $4 
		WHERE id = $5
	`

	deptID := sql.NullInt64{}
	if user.DepartmentID > 0 {
		deptID.Valid = true
		deptID.Int64 = int64(user.DepartmentID)
	}

	result, err := r.DB.Exec(query,
		user.Name,
		user.Email,
		user.Role,
		deptID,
		id,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("email sudah digunakan")
		}
		return fmt.Errorf("gagal update user: %w", err)
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

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal hapus user: %w", err)
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

func (r *UserRepository) CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gagal cek email: %w", err)
	}

	return exists, nil
}


