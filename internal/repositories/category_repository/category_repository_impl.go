package categoryrepository

import (
	"context"
	"marketly-app/internal/models"

	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{DB: db}
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, data *models.Category) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *CategoryRepositoryImpl) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category

	if err := r.DB.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Category, int, error) {
	var (
		categories []*models.Category
		count int64
	)

	query := r.DB.WithContext(ctx).Model(&models.Category{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, int(count), nil
}

func (r *CategoryRepositoryImpl) FindById(ctx context.Context, categoryId int) (*models.Category, error) {
	var category models.Category

	if err := r.DB.WithContext(ctx).Where("id = ?", categoryId).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepositoryImpl) Update(ctx context.Context, categoryId int, data *models.Category) error {
	var existing models.Category

	if err := r.DB.WithContext(ctx).Where("id = ?", categoryId).First(&existing).Error; err != nil {
		return err
	}

	existing.Name = data.Name
	return r.DB.WithContext(ctx).Save(&existing).Error
}

func (r *CategoryRepositoryImpl) Delete(ctx context.Context, categoryId int) error {
	return r.DB.WithContext(ctx).Delete(&models.Category{}, "id = ?", categoryId).Error
}
