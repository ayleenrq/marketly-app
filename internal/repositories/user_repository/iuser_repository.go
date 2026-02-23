package userrepository

import (
	"context"
	"marketly-app/internal/models"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindRoleUser(ctx context.Context) (*models.Role, error)
	FindById(ctx context.Context, userID int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}
