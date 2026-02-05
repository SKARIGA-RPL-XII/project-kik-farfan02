package service

import (
	"errors"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
	"fmt"
)

type LokasiService struct {
	lokasiRepo repository.LokasiRepository
}

func NewLokasiService(lokasiRepo repository.LokasiRepository) *LokasiService {
	return &LokasiService{lokasiRepo: lokasiRepo}
}

func (s *LokasiService) CreateLokasi(lok *models.CreateLocationRequest) (int, error) {
	if lok.NamaLokasi == "" {
		return 0, errors.New("nama lokasi wajib diisi")
	}

	if lok.Capacity <= 0 {
		return 0, errors.New("kapasitas lokasi harus lebih dari 0")
	}

	if len(lok.Ruangan) == 0 {
		return 0, errors.New("setidaknya harus ada 1 ruangan")
	}

	for i, ruang := range lok.Ruangan {
		if ruang.NamaRuangan == "" {
			return 0, errors.New("nama ruangan pada posisi " + fmt.Sprint(i+1) + " tidak boleh kosong")
		}
		if ruang.Capacity <= 0 {
			return 0, errors.New("kapasitas ruangan pada posisi " + fmt.Sprint(i+1) + " harus lebih dari 0")
		}
	}

	return s.lokasiRepo.CreateLokasi(lok)
}

func (s *LokasiService) GetAllLocations() ([]models.Location, error) {
	return s.lokasiRepo.GetAllLocations()
}

func (s *LokasiService) GetLocationDetails(lokID int) ([]models.DetailLocation, error) {
	return s.lokasiRepo.GetLocationDetails(lokID)
}

func (s *LokasiService) GetAllLocationsWithDetails() ([]models.LocationWithDetails, error) {
	return s.lokasiRepo.GetAllLocationsWithDetails()
}

func (s *LokasiService) UpdateLokasi(id int, lok *models.CreateLocationRequest) error {
	if lok.NamaLokasi == "" {
		return errors.New("nama lokasi wajib diisi")
	}

	if lok.Capacity <= 0 {
		return errors.New("kapasitas lokasi harus lebih dari 0")
	}

	if len(lok.Ruangan) == 0 {
		return errors.New("setidaknya harus ada 1 ruangan")
	}

	for i, ruang := range lok.Ruangan {
		if ruang.NamaRuangan == "" {
			return errors.New("nama ruangan pada posisi " + fmt.Sprint(i+1) + " tidak boleh kosong")
		}
		if ruang.Capacity <= 0 {
			return errors.New("kapasitas ruangan pada posisi " + fmt.Sprint(i+1) + " harus lebih dari 0")
		}
	}

	return s.lokasiRepo.UpdateLokasi(id, lok)
}

func (s *LokasiService) DeleteLokasi(id int) error {
	return s.lokasiRepo.DeleteLokasi(id)
}

func (s *LokasiService) UpdateDetailLokasi(id int, detail *models.DetailLocation) error {
	if detail.NamaRuangan == "" {
		return errors.New("nama ruangan wajib diisi")
	}

	if detail.Capacity <= 0 {
		return errors.New("kapasitas ruangan harus lebih dari 0")
	}

	return s.lokasiRepo.UpdateDetailLokasi(id, detail)
}

func (s *LokasiService) DeleteDetailLokasi(id int) error {
	return s.lokasiRepo.DeleteDetailLokasi(id)
}
