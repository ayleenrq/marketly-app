package sellerrepository

import (
	"context"
	"marketly-app/internal/models"
)

type ISellerRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindRoleSeller(ctx context.Context) (*models.Role, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.User, int64, error)
	FindById(ctx context.Context, sellerID int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}
