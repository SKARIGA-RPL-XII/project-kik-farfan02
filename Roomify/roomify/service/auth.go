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

type AuthService struct {
	authRepo repository.AuthRepository
	cfg      *config.Config
}

func NewAuthService(authRepo repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		cfg:      cfg,
	}
}

func (a *AuthService) Login(creds *models.Credentials) (*models.LoginResponse, error) {
	user, err := a.authRepo.Login(creds.Email)
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

	token, err := util.GenerateToken(user.ID, user.Role, a.cfg)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &models.LoginResponse{
		Token:   token,
		Message: "anda berhasil login",
		Name:    user.Name,
	}, nil
}

func (a *AuthService) ChangePassword(userID int, oldPassword, newPassword, confirmPassword string) error {
	if newPassword != confirmPassword {
		return errors.New("password baru dan konfirmasi tidak sama")
	}

	if len(newPassword) < 6 {
		return errors.New("password minimal 6 karakter")
	}

	user, err := a.authRepo.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("terjadi kesalahan pada server")
	}

	cleanHash := strings.Trim(user.Pass, `"`)

	err = bcrypt.CompareHashAndPassword([]byte(cleanHash), []byte(oldPassword))
	if err != nil {
		return errors.New("password lama salah")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}

	return a.authRepo.ChangePassword(userID, string(hashedPassword))
}

