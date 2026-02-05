package repository

import (
	"database/sql"
	"errors"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
)

type SettingRepository struct {
	DB *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{DB: db}
}

func (r *SettingRepository) GetSetting(key string) (string, error) {
	var value string
	err := r.DB.QueryRow("SELECT value FROM app_settings WHERE key = $1", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("setting not found")
		}
		return "", err
	}
	return value, nil
}

func (r *SettingRepository) GetAllSettings() ([]models.Setting, error) {
	rows, err := r.DB.Query("SELECT key, value FROM app_settings ORDER BY key ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var s models.Setting
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}

	return settings, nil
}

func (r *SettingRepository) UpdateSetting(key, value string) error {
	_, err := r.DB.Exec(`
		INSERT INTO app_settings (key, value) 
		VALUES ($1, $2) 
		ON CONFLICT (key) DO UPDATE SET value = $2
	`, key, value)
	return err
}