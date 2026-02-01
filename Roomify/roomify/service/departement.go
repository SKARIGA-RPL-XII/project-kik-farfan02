package service

import (
	"errors"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
)

type DeptService struct {
	deptRepo repository.DeptRepository
}

func NewDeptService(deptRepo repository.DeptRepository) *DeptService {
	return &DeptService{deptRepo: deptRepo}
}

func (s *DeptService) InputDepartment(dpt *models.Departement) error {
	exists, err := s.deptRepo.CheckDepartmentExists(dpt.Nama_dtm)
	if err != nil {
		return errors.New("terjadi kesalahan pada server") 
	}

	if exists {
		return errors.New("department sudah terdaftar")
	}

	if dpt.Code == "" {
		return errors.New("code departement harus terisi")
	}
	println(dpt.Code)

	return s.deptRepo.InputDepartment(dpt)
}

func (s *DeptService) GetAllDepartment() ([]models.Departement, error) {
	return s.deptRepo.GetAllDepartment()
}

func (s *DeptService) UpdateDepartment(id int, dpt *models.Departement) error {
	allDepts, err := s.deptRepo.GetAllDepartment()
	if err != nil {
		return errors.New("terjadi kesalahan pada server")
	}

	found := false
	for _, dept := range allDepts {
		if dept.ID == id {
			found = true
			break
		}
	}

	if !found {
		return errors.New("department tidak ditemukan")
	}

	return s.deptRepo.UpdateDepartment(id, dpt)
}

func (s *DeptService) DeleteDepartment(id int) error {
	if id <= 0 {
		return errors.New("id harap dimasukkan")
	}
	return s.deptRepo.DeleteDepartment(id)
}