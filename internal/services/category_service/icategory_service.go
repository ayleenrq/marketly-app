package categoryservice

import (
	"context"
	categoryrequest "marketly-app/internal/dto/request/category_request"
	"marketly-app/internal/models"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, req categoryrequest.CreateCategoryRequest) error
	GetAllCategory(ctx context.Context, page, limit int, search string) ([]*models.Category, int, error)
	GetByIdCategory(ctx context.Context, categoryId int) (*models.Category, error)
	UpdateCategory(ctx context.Context, categoryId int, req categoryrequest.UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, categoryId int) error
}
