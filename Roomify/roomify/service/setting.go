package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
)

type SettingService struct {
	settingRepo repository.SettingRepository
}

func NewSettingService(settingRepo repository.SettingRepository) *SettingService {
	return &SettingService{settingRepo: settingRepo}
}


func (s *SettingService) GetSetting(key string) (string, error) {
	return s.settingRepo.GetSetting(key)
}


func (s *SettingService) GetAllSettings() ([]models.Setting, error) {
	return s.settingRepo.GetAllSettings()
}

func (s *SettingService) UpdateSetting(key, value string) error {
	return s.settingRepo.UpdateSetting(key, value)
}

func (s *SettingService) GetWorkingHours() (startTime, endTime string, err error) {
	start, err := s.settingRepo.GetSetting("working_hours_start")
	if err != nil {
		return "", "", errors.New("working hours not configured")
	}
	end, err := s.settingRepo.GetSetting("working_hours_end")
	if err != nil {
		return "", "", errors.New("working hours not configured")
	}
	return start, end, nil
}

func (s *SettingService) GetHolidays() ([]string, error) {
	holidaysJSON, err := s.settingRepo.GetSetting("holidays")
	if err != nil {
		return nil, errors.New("holidays not configured")
	}
	
	var holidays []string
	if err := json.Unmarshal([]byte(holidaysJSON), &holidays); err != nil {
		return nil, err
	}
	return holidays, nil
}

func (s *SettingService) IsHoliday(date time.Time) (bool, error) {
	holidays, err := s.GetHolidays()
	if err != nil {
		return false, err
	}
	
	dateStr := date.Format("2006-01-02")
	for _, h := range holidays {
		if h == dateStr {
			return true, nil
		}
	}
	return false, nil
}