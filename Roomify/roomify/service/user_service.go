package service

import (
	"errors"
	"strings"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/models"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *UserService) Login(creds *models.Credentials) (*models.LoginResponse, error) {
	user, err := s.userRepo.Login(creds.Email)
	if err != nil {
		if err.Error() == "email tidak ditemukan" {
			return nil, errors.New("email tidak ditemukan")
		}
		return nil, errors.New("terjadi kesalahan pada server")
	}

	cleanHash := strings.Trim(user.Pass, `"`)

	err = bcrypt.CompareHashAndPassword([]byte(cleanHash), []byte(creds.Pass))
	if err != nil {
		return nil, errors.New("password salah")
	}

	token, err := util.GenerateToken(user.ID, user.Role, s.cfg)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &models.LoginResponse{
		Token:   token,
		Message: "anda berhasil login",
		Name:    user.Name,
	}, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	if !strings.Contains(user.Email, "@") {
		return errors.New("email tidak valid")
	}

	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}

	user.Pass = string(hashedPassword)
	
	if user.Role == "" {
		user.Role = "user"
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetUsers(filter models.UserFilter) ([]models.User, error) {
	return s.userRepo.GetUsers(filter)
}

func (s *UserService) UpdateUser(id int, userData *models.User) error {
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("terjadi kesalahan pada server")
	}


	exists, err := s.userRepo.CheckEmailExists(userData.Email)
	if err != nil {
		return errors.New("terjadi kesalahan pada server")
	}

	if exists {
		currentUser, _ := s.userRepo.GetUserByID(id)
		if currentUser.Email != userData.Email {
			return errors.New("email sudah digunakan")
		}
	}

	return s.userRepo.UpdateUser(id, userData)
}

func (s *UserService) DeleteUser(id int) error {
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("terjadi kesalahan pada server")
	}

	return s.userRepo.DeleteUser(id)
}



