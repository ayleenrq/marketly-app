package categoryservice

import (
	"context"
	categoryrequest "marketly-app/internal/dto/request/category_request"
	"marketly-app/internal/models"
	categoryrepository "marketly-app/internal/repositories/category_repository"
	errorresponse "marketly-app/pkg/constant/error_response"
	"strings"
)

type CategoryServiceImpl struct {
	categoryRepo categoryrepository.ICategoryRepository
}

func NewCategoryServiceImpl(categoryRepo categoryrepository.ICategoryRepository) ICategoryService {
	return &CategoryServiceImpl{categoryRepo: categoryRepo}
}

func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, req categoryrequest.CreateCategoryRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Category name is required", 400)
	}

	category := &models.Category{
		Name: req.Name,
	}
	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryServiceImpl) GetAllCategory(ctx context.Context, page, limit int, search string) ([]*models.Category, int, error) {
	offset := (page - 1) * limit
	return s.categoryRepo.FindAll(ctx, limit, offset, search)
}

func (s *CategoryServiceImpl) GetByIdCategory(ctx context.Context, categoryId int) (*models.Category, error) {
	return s.categoryRepo.FindById(ctx, categoryId)
}

func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, categoryId int, req categoryrequest.UpdateCategoryRequest) error {
	category, err := s.categoryRepo.FindById(ctx, categoryId)
	if err != nil {
		return err
	}

	category.Name = req.Name
	return s.categoryRepo.Update(ctx, categoryId, category)
}

func (s *CategoryServiceImpl) DeleteCategory(ctx context.Context, categoryId int) error {
	return s.categoryRepo.Delete(ctx, categoryId)
}
