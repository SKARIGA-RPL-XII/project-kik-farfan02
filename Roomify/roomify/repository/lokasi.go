package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type LokasiRepository struct {
	DB *sql.DB
}

func NewLokasiRepository(db *sql.DB) *LokasiRepository {
	return &LokasiRepository{DB: db}
}

func (l *LokasiRepository) CreateLokasi(lok *models.CreateLocationRequest) (int, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var lokID int
	err = tx.QueryRow(`
		INSERT INTO lokasi (nama_lokasi, capacity)
		VALUES ($1, $2) RETURNING id
	`, lok.NamaLokasi, lok.Capacity).Scan(&lokID)
	if err != nil {
		return 0, fmt.Errorf("failed to create lokasi: %w", err)
	}

	for _, ruang := range lok.Ruangan {
		_, err := tx.Exec(`
			INSERT INTO detail_lokasi (id_lokasi, nama_ruangan, capacity)
			VALUES ($1, $2, $3)
		`, lokID, ruang.NamaRuangan, ruang.Capacity)
		if err != nil {
			return 0, fmt.Errorf("failed to create detail lokasi: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return lokID, nil
}

func (l *LokasiRepository) GetAllLocations() ([]models.Location, error) {
	query := `
		SELECT l.id, l.nama_lokasi, l.capacity, l.created_at
		FROM lokasi l
		ORDER BY l.id ASC
	`

	rows, err := l.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query locations: %w", err)
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		err := rows.Scan(
			&loc.ID,
			&loc.NamaLokasi,
			&loc.Capacity,
			&loc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return locations, nil
}

func (l *LokasiRepository) GetLocationDetails(lokID int) ([]models.DetailLocation, error) {
	query := `
		SELECT id, id_lokasi, nama_ruangan, capacity
		FROM detail_lokasi
		WHERE id_lokasi = $1
		ORDER BY id ASC
	`

	rows, err := l.DB.Query(query, lokID)
	if err != nil {
		return nil, fmt.Errorf("failed to query details: %w", err)
	}
	defer rows.Close()

	var details []models.DetailLocation
	for rows.Next() {
		var detail models.DetailLocation
		err := rows.Scan(
			&detail.ID,
			&detail.IDLokasi,
			&detail.NamaRuangan,
			&detail.Capacity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan detail: %w", err)
		}
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return details, nil
}

func (l *LokasiRepository) GetAllLocationsWithDetails() ([]models.LocationWithDetails, error) {
	locations, err := l.GetAllLocations()
	if err != nil {
		return nil, err
	}

	var result []models.LocationWithDetails
	for _, loc := range locations {
		details, err := l.GetLocationDetails(loc.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.LocationWithDetails{
			Location: loc,
			Details:  details,
		})
	}

	return result, nil
}

func (l *LokasiRepository) DeleteLokasi(id int) error {
	tx, err := l.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`DELETE FROM detail_lokasi WHERE id_lokasi = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete details: %w", err)
	}

	result, err := tx.Exec(`DELETE FROM lokasi WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete lokasi: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("lokasi tidak ditemukan")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (l *LokasiRepository) UpdateLokasi(id int, lok *models.CreateLocationRequest) error {
	tx, err := l.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		UPDATE lokasi 
		SET nama_lokasi = $1, capacity = $2 
		WHERE id = $3
	`, lok.NamaLokasi, lok.Capacity, id)
	if err != nil {
		return fmt.Errorf("failed to update lokasi: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM detail_lokasi WHERE id_lokasi = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete old details: %w", err)
	}

	for _, ruang := range lok.Ruangan {
		_, err := tx.Exec(`
			INSERT INTO detail_lokasi (id_lokasi, nama_ruangan, capacity)
			VALUES ($1, $2, $3)
		`, id, ruang.NamaRuangan, ruang.Capacity)
		if err != nil {
			return fmt.Errorf("failed to create detail lokasi: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (l *LokasiRepository) UpdateDetailLokasi(id int, detail *models.DetailLocation) error {
	query := `
		UPDATE detail_lokasi 
		SET nama_ruangan = $1, capacity = $2 
		WHERE id = $3
	`

	result, err := l.DB.Exec(query, detail.NamaRuangan, detail.Capacity, id)
	if err != nil {
		return fmt.Errorf("failed to update detail lokasi: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("detail lokasi tidak ditemukan")
	}

	return nil
}

func (l *LokasiRepository) DeleteDetailLokasi(id int) error {
	query := `DELETE FROM detail_lokasi WHERE id = $1`

	result, err := l.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete detail lokasi: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("detail lokasi tidak ditemukan")
	}

	return nil
}

