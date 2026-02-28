package adminrepository

import (
	"context"
	"marketly-app/internal/models"
)

type IAdminRepository interface {
	Create(ctx context.Context, admin *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindRoleAdmin(ctx context.Context) (*models.Role, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.User, int, error)
	FindByAdminID(ctx context.Context, adminID int) (*models.User, error)
	Update(ctx context.Context, admin *models.User) error
}
