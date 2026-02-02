package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type DeptRepository struct {
	DB *sql.DB
}

func NewDeptRepository(db *sql.DB) *DeptRepository {
	return &DeptRepository{DB: db}
}

func (v *DeptRepository) CheckDepartmentExists(namaDtm string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM department WHERE nama_dtm = $1)`

	err := v.DB.QueryRow(query, namaDtm).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gagal cek department: %w", err)
	}

	return exists, nil
}


func (v *DeptRepository) InputDepartment(dpt *models.Departement) error {
	query := `INSERT INTO department (nama_dtm, code) VALUES($1, $2) RETURNING id`	
	err := v.DB.QueryRow(query,
		dpt.Nama_dtm,
		dpt.Code,
		).Scan(&dpt.ID)

	if err != nil {
		if err.Error() == "failed to create department" {
			return errors.New("gagal membuat department")
		}
		return fmt.Errorf("gagal membuat department: %w", err)
	}
	return nil
}

func (v *DeptRepository) GetAllDepartment() ([]models.Departement, error) {
	query := `SELECT id, nama_dtm, code FROM department ORDER BY id ASC`
	rows, err := v.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data department: %w", err)
	}
	defer rows.Close()

	var departments []models.Departement
	for rows.Next() {
		var dept models.Departement
		err := rows.Scan(&dept.ID, &dept.Nama_dtm, &dept.Code)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data: %w", err)
		}
		departments = append(departments, dept)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi rows: %w", err)
	}

	return departments, nil
}


func (v *DeptRepository) UpdateDepartment( dpt *models.Departement) error {
	query := `UPDATE department SET nama_dtm = $1, code = $2 WHERE id = $3`

	result, err := v.DB.Exec(query, dpt.Nama_dtm, dpt.Code, dpt.ID)
	if err != nil {
		return fmt.Errorf("gagal update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal cek rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("department tidak ditemukan")
	}

	return nil
}

func (v *DeptRepository) DeleteDepartment(id int) error {
	query := `DELETE FROM department WHERE id = $1`

	result, err := v.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal hapus department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal cek rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("department tidak ditemukan")
	}

	return nil
}