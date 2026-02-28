package sellerservice

import (
	"context"
	sellerrequest "marketly-app/internal/dto/request/seller_request"
	"marketly-app/internal/models"
)

type ISellerService interface {
	Register(ctx context.Context, req sellerrequest.RegisterSellerRequest) error
	Login(ctx context.Context, req sellerrequest.LoginSellerRequest) (string, error)
	GetProfile(ctx context.Context, sellerID int) (*models.User, error)
	GetAllSeller(ctx context.Context, page, limit int, search string) ([]*models.User, int64, error)
	GetByIdSeller(ctx context.Context, sellerID int) (*models.User, error)
	UpdateProfile(ctx context.Context, sellerID int, req sellerrequest.UpdateSellerRequest) error
	UpdatePhoto(ctx context.Context, sellerID int, req sellerrequest.UpdatePhotoSellerRequest) error
	ChangeEmail(ctx context.Context, sellerID int, req sellerrequest.ChangeEmailSellerRequest) error
	ChangePassword(ctx context.Context, sellerID int, req sellerrequest.ChangePasswordSellerRequest) error
	Logout(ctx context.Context, sellerID int) error
}