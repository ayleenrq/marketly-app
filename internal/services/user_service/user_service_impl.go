package userservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	if strings.TrimSpace(req.Username) == "" {
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
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nomor Handphone wajib diisi", 400)
	}
	if strings.TrimSpace(req.Address) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Alamat wajib diisi", 400)
	}

	if !utils.IsValidEmail(req.Email) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
	}

	existsUsername, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan username", 500)
	}
	if existsUsername != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username sudah digunakan", 409)
	}

	existsEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan email", 500)
	}
	if existsEmail != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Email sudah digunakan", 409)
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

	user := &models.User{
		Username:    req.Username,
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashedPass,
		PhoneNumber: &req.PhoneNumber,
		Address:     &req.Address,
		RoleID:      role.ID,
	}

	if req.PhotoFile != nil {
		photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/users")
		if err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
		}
		user.PhotoURL = &photoURL
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal menyimpan data pengguna", 500)
	}

	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req userrequest.LoginUserRequest) (string, error) {
	if strings.TrimSpace(req.Password) == "" {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Kata sandi wajib diisi", 400)
	}

	email := strings.TrimSpace(req.Email)
	username := strings.TrimSpace(req.Username)

	if email == "" && username == "" {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Email atau Username wajib diisi", 400)
	}

	var user *models.User
	var err error

	if email != "" {
		if !utils.IsValidEmail(email) {
			return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
		}

		user, err = s.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return "", errorresponse.NewCustomError(errorresponse.ErrNotFound, "Email tidak valid", 400)
		}
	} else {
		user, err = s.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return "", errorresponse.NewCustomError(errorresponse.ErrNotFound, "Username tidak valid", 400)
		}
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

func (u *UserServiceImpl) UpdateProfile(ctx context.Context, userID int, req userrequest.UpdateUserRequest) error {
	user, err := u.userRepo.FindById(ctx, userID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name is required", 400)
	}

	user.Name = req.Name

	if strings.TrimSpace(req.Username) != "" && req.Username != user.Username {
		existsUsername, err := u.userRepo.FindByUsername(ctx, req.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan username", 500)
		}
		if existsUsername != nil && existsUsername.ID != user.ID {
			return errorresponse.NewCustomError(errorresponse.ErrExists, "Username sudah digunakan", 409)
		}
		user.Username = req.Username
	}

	if strings.TrimSpace(req.PhoneNumber) != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	if strings.TrimSpace(req.Address) != "" {
		user.Address = &req.Address
	}

	if req.PhotoFile != nil {
		photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/users")
		if err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
		}
		user.PhotoURL = &photoURL
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update user profile", 500)
	}

	return nil
}

func (u *UserServiceImpl) UpdatePhoto(ctx context.Context, userID int, req userrequest.UpdatePhotoUserRequest) error {
	if req.PhotoFile == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Photo file is required", 400)
	}

	user, err := u.userRepo.FindById(ctx, userID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}

	photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/users")
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
	}

	user.PhotoURL = &photoURL

	if err := u.userRepo.Update(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update photo", 500)
	}

	return nil
}

func (u *UserServiceImpl) ChangePassword(ctx context.Context, userID int, req userrequest.ChangePasswordUserRequest) error {
	if strings.TrimSpace(req.OldPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password lama wajib diisi", 400)
	}
	if strings.TrimSpace(req.NewPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password baru wajib diisi", 400)
	}

	user, err := u.userRepo.FindById(ctx, userID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password lama salah", 400)
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Gagal meng-hash password", 400)
	}

	user.Password = hashed

	if err := u.userRepo.Update(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to change password", 500)
	}

	return nil
}

func (u *UserServiceImpl) ChangeEmail(ctx context.Context, userID int, req userrequest.ChangeEmailUserRequest) error {
	if strings.TrimSpace(req.NewEmail) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Email baru wajib diisi", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password wajib diisi", 400)
	}
	if !utils.IsValidEmail(req.NewEmail) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
	}

	user, err := u.userRepo.FindById(ctx, userID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "User not found", 404)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password salah", 400)
	}

	existsEmail, err := u.userRepo.FindByEmail(ctx, req.NewEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan email", 500)
	}
	if existsEmail != nil && existsEmail.ID != user.ID {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Email sudah digunakan", 409)
	}

	user.Email = req.NewEmail

	if err := u.userRepo.Update(ctx, user); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to change email", 500)
	}

	return nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, userID int) error {
	fmt.Printf("User %d logged out\n", userID)
	return nil
}
