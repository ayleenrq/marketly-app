package userservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	userrequest "marketly-app/internal/dto/request/user_request"
	"marketly-app/internal/models"
	userrepo "marketly-app/internal/repositories/user_repository"
	errorresponse "marketly-app/pkg/constant/error_response"
	"marketly-app/pkg/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	userRepo   userrepo.IUserRepository
	cloudinary *cloudinary.Cloudinary
}

func NewUserServiceImpl(userRepo userrepo.IUserRepository, cld *cloudinary.Cloudinary) IUserService {
	return &UserServiceImpl{userRepo: userRepo, cloudinary: cld}
}

func (s *UserServiceImpl) Register(ctx context.Context, req userrequest.RegisterUserRequest) error {
	if strings.TrimSpace(req.NIK) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "NIK wajib diisi", 400)
	}
	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nama wajib diisi", 400)
	}
	if strings.TrimSpace(req.Email) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Email wajib diisi", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password wajib diisi", 400)
	}
	if strings.TrimSpace(req.TempatLahir) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Tempat Lahir wajib diisi", 400)
	}
	if strings.TrimSpace(req.Agama) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Agama wajib diisi", 400)
	}
	if strings.TrimSpace(req.Address) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Alamat wajib diisi", 400)
	}
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nomor Handphone wajib diisi", 400)
	}
	if strings.TrimSpace(req.Status) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Status Perkawinan wajib diisi", 400)
	}
	if strings.TrimSpace(req.ReasonRegister) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Alasan wajib diisi", 400)
	}

	if !utils.IsValidEmail(req.Email) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
	}

	if !utils.IsValidNIK(req.NIK) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "NIK harus terdiri dari 16 digit angka", 400)
	}

	if !utils.IsNumeric(req.PhoneNumber) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nomor Handphone harus berupa angka", 400)
	}

	if req.PhotoFile == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Foto wajib diunggah", 400)
	}

	existsNIK, err := s.userRepo.FindByNIK(ctx, req.NIK)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan NIK", 500)
	}
	if existsNIK != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "NIK sudah digunakan", 409)
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Gagal meng-hash password", 400)
	}

	role, err := s.userRepo.FindRoleUser(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Role 'user' not found.", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to get role user", 500)
	}

	birth, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format tanggal lahir harus YYYY-MM-DD", 400)
	}

	photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "civi-id/users")
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
	}

	genderML, err := utils.DetectGenderML(req.PhotoFile)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendeteksi jenis kelamin", 500)
	}

	jenisKelamin := utils.MLToIndo(genderML)

	user := &models.User{
		NIK:            &req.NIK,
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPass,
		JenisKelamin:   &jenisKelamin,
		TempatLahir:    &req.TempatLahir,
		BirthDate:      &birth,
		Agama:          &req.Agama,
		Address:        &req.Address,
		PhoneNumber:    &req.PhoneNumber,
		Status:         &req.Status,
		ReasonRegister: &req.ReasonRegister,
		PhotoURL:       &photoURL,
		GenderML:       &genderML,
		RoleID:         role.ID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal menyimpan data pengguna", 500)
	}

	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req userrequest.LoginUserRequest) (string, error) {
	if strings.TrimSpace(req.NIK) == "" {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "NIK wajib diisi", 400)
	}

	if !utils.IsValidNIK(req.NIK) {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "NIK harus terdiri dari 16 digit angka", 400)
	}

	user, err := s.userRepo.FindByNIK(ctx, req.NIK)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrNotFound, "NIK tidak valid", 400)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password salah", 400)
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal membuat token autentikasi", 500)
	}

	return token, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID int, req userrequest.UpdateUserRequest) error {
	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Address != "" {
		user.Address = &req.Address
	}

	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update user", 500)
	}

	return nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, userID int) error {
	fmt.Printf("User %d logged out\n", userID)
	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
