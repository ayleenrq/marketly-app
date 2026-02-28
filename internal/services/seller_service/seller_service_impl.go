package sellerservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sellerrequest "marketly-app/internal/dto/request/seller_request"
	"marketly-app/internal/models"
	sellerrepo "marketly-app/internal/repositories/seller_repository"
	errorresponse "marketly-app/pkg/constant/error_response"
	"marketly-app/pkg/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"gorm.io/gorm"
)

type SellerServiceImpl struct {
	sellerRepo sellerrepo.ISellerRepository
	cloudinary *cloudinary.Cloudinary
}

func NewSellerServiceImpl(sellerRepo sellerrepo.ISellerRepository, cld *cloudinary.Cloudinary) ISellerService {
	return &SellerServiceImpl{sellerRepo: sellerRepo, cloudinary: cld}
}

func (s *SellerServiceImpl) Register(ctx context.Context, req sellerrequest.RegisterSellerRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Username wajib diisi", 400)
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
	if strings.TrimSpace(req.StoreName) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nama toko wajib diisi", 400)
	}
	if strings.TrimSpace(req.StoreDescription) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Deskripsi toko wajib diisi", 400)
	}

	if !utils.IsValidEmail(req.Email) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
	}

	existsUsername, err := s.sellerRepo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan username", 500)
	}
	if existsUsername != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username sudah digunakan", 409)
	}

	existsEmail, err := s.sellerRepo.FindByEmail(ctx, req.Email)
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

	role, err := s.sellerRepo.FindRoleSeller(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Role 'seller' not found.", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to get role seller", 500)
	}

	seller := &models.User{
		Username:         req.Username,
		Name:             req.Name,
		Email:            req.Email,
		Password:         hashedPass,
		PhoneNumber:      &req.PhoneNumber,
		Address:          &req.Address,
		StoreName:        &req.StoreName,
		StoreDescription: &req.StoreDescription,
		RoleID:           role.ID,
	}

	if req.PhotoFile != nil {
		photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/sellers")
		if err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
		}
		seller.PhotoURL = &photoURL
	}

	if err := s.sellerRepo.Create(ctx, seller); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal menyimpan data penjual", 500)
	}

	return nil
}

func (s *SellerServiceImpl) Login(ctx context.Context, req sellerrequest.LoginSellerRequest) (string, error) {
	if strings.TrimSpace(req.Password) == "" {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Kata sandi wajib diisi", 400)
	}

	email := strings.TrimSpace(req.Email)
	username := strings.TrimSpace(req.Username)

	if email == "" && username == "" {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Email atau Username wajib diisi", 400)
	}

	var seller *models.User
	var err error

	if email != "" {
		if !utils.IsValidEmail(email) {
			return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
		}

		seller, err = s.sellerRepo.FindByEmail(ctx, email)
		if err != nil {
			return "", errorresponse.NewCustomError(errorresponse.ErrNotFound, "Email tidak valid", 400)
		}
	} else {
		seller, err = s.sellerRepo.FindByUsername(ctx, username)
		if err != nil {
			return "", errorresponse.NewCustomError(errorresponse.ErrNotFound, "Username tidak valid", 400)
		}
	}

	if strings.ToLower(seller.Role.Name) != "seller" {
		return "", errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "Unauthorized access", 401)
	}

	if !utils.CheckPasswordHash(req.Password, seller.Password) {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password salah", 400)
	}

	token, err := utils.GenerateToken(seller.ID, seller.RoleID)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal membuat token autentikasi", 500)
	}

	return token, nil
}

func (s *SellerServiceImpl) GetProfile(ctx context.Context, sellerID int) (*models.User, error) {
	seller, err := s.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}
	return seller, nil
}

func (s *SellerServiceImpl) GetAllSeller(ctx context.Context, page, limit int, search string) ([]*models.User, int64, error) {
	offset := (page - 1) * limit
	return s.sellerRepo.FindAll(ctx, limit, offset, search)
}

func (s *SellerServiceImpl) GetByIdSeller(ctx context.Context, sellerID int) (*models.User, error) {
	seller, err := s.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}
	return seller, nil
}

func (u *SellerServiceImpl) UpdateProfile(ctx context.Context, sellerID int, req sellerrequest.UpdateSellerRequest) error {
	seller, err := u.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name is required", 400)
	}

	seller.Name = req.Name

	if strings.TrimSpace(req.Username) != "" && req.Username != seller.Username {
		existsUsername, err := u.sellerRepo.FindByUsername(ctx, req.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan username", 500)
		}
		if existsUsername != nil && existsUsername.ID != seller.ID {
			return errorresponse.NewCustomError(errorresponse.ErrExists, "Username sudah digunakan", 409)
		}
		seller.Username = req.Username
	}

	if strings.TrimSpace(req.PhoneNumber) != "" {
		seller.PhoneNumber = &req.PhoneNumber
	}

	if strings.TrimSpace(req.Address) != "" {
		seller.Address = &req.Address
	}

	if req.PhotoFile != nil {
		photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/sellers")
		if err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
		}
		seller.PhotoURL = &photoURL
	}

	if err := u.sellerRepo.Update(ctx, seller); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update seller profile", 500)
	}

	return nil
}

func (u *SellerServiceImpl) UpdatePhoto(ctx context.Context, sellerID int, req sellerrequest.UpdatePhotoSellerRequest) error {
	if req.PhotoFile == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Photo file is required", 400)
	}

	seller, err := u.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}

	photoURL, err := utils.UploadToCloudinary(req.PhotoFile, "marketly/sellers")
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mengunggah foto", 500)
	}

	seller.PhotoURL = &photoURL

	if err := u.sellerRepo.Update(ctx, seller); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update photo", 500)
	}

	return nil
}

func (u *SellerServiceImpl) ChangePassword(ctx context.Context, sellerID int, req sellerrequest.ChangePasswordSellerRequest) error {
	if strings.TrimSpace(req.OldPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password lama wajib diisi", 400)
	}
	if strings.TrimSpace(req.NewPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password baru wajib diisi", 400)
	}

	seller, err := u.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}

	if !utils.CheckPasswordHash(req.OldPassword, seller.Password) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password lama salah", 400)
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Gagal meng-hash password", 400)
	}

	seller.Password = hashed

	if err := u.sellerRepo.Update(ctx, seller); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to change password", 500)
	}

	return nil
}

func (u *SellerServiceImpl) ChangeEmail(ctx context.Context, sellerID int, req sellerrequest.ChangeEmailSellerRequest) error {
	if strings.TrimSpace(req.NewEmail) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Email baru wajib diisi", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password wajib diisi", 400)
	}
	if !utils.IsValidEmail(req.NewEmail) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Format email tidak valid", 400)
	}

	seller, err := u.sellerRepo.FindById(ctx, sellerID)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "Seller not found", 404)
	}

	if !utils.CheckPasswordHash(req.Password, seller.Password) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password salah", 400)
	}

	existsEmail, err := u.sellerRepo.FindByEmail(ctx, req.NewEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Gagal mendapatkan email", 500)
	}
	if existsEmail != nil && existsEmail.ID != seller.ID {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Email sudah digunakan", 409)
	}

	seller.Email = req.NewEmail

	if err := u.sellerRepo.Update(ctx, seller); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to change email", 500)
	}

	return nil
}

func (u *SellerServiceImpl) Logout(ctx context.Context, sellerID int) error {
	fmt.Printf("Seller %d logged out\n", sellerID)
	return nil
}
