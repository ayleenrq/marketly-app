package categoryrepository

import (
	"context"
	"marketly-app/internal/models"
)

type ICategoryRepository interface {
	Create(ctx context.Context, data *models.Category) error
	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Category, int, error)
	FindById(ctx context.Context, categoryId int) (*models.Category, error)
	Update(ctx context.Context, categoryId int, data *models.Category) error
	Delete(ctx context.Context, categoryId int) error
}
