package adminservice

import (
	"context"
	adminrequest "marketly-app/internal/dto/request/admin_request"
	"marketly-app/internal/models"
)

type IAdminService interface {
	Register(ctx context.Context, req adminrequest.RegisterAdminRequest) error
	Login(ctx context.Context, req adminrequest.LoginAdminRequest) (string, error)
	GetProfile(ctx context.Context, adminID int) (*models.User, error)
	GetAllAdmin(ctx context.Context, page, limit int, search string) ([]*models.User, int, error)
	GetByIdAdmin(ctx context.Context, adminID int) (*models.User, error)
	UpdateProfile(ctx context.Context, adminID int, req adminrequest.UpdateProfileRequest) error
	Logout(ctx context.Context, adminID int) error
}
