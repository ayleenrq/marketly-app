package adminrepository

import (
	"context"
	"marketly-app/internal/models"

	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	DB *gorm.DB
}

func NewAdminRepositoryImpl(db *gorm.DB) IAdminRepository {
	return &AdminRepositoryImpl{DB: db}
}

func (r *AdminRepositoryImpl) Create(ctx context.Context, admin *models.User) error {
	return r.DB.WithContext(ctx).Create(admin).Error
}

func (r *AdminRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var admin models.User
	if err := r.DB.WithContext(ctx).Preload("Role").Joins("JOIN roles ON roles.id = users.role_id").Where("users.email = ?", email).Where("LOWER(roles.name) = LOWER(?)", "ADMIN").First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a *AdminRepositoryImpl) FindRoleAdmin(ctx context.Context) (*models.Role, error) {
	var role models.Role
	if err := a.DB.WithContext(ctx).Where("LOWER(name) = LOWER(?)", "admin").First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *AdminRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*models.User, int, error) {
	var (
		admins []*models.User
		count int64
	)

	query := r.DB.WithContext(ctx).Model(&models.User{}).Preload("Role").Joins("JOIN roles ON roles.id = users.role_id").Where("LOWER(roles.name) = LOWER(?)", "ADMIN")

	if search != "" {
		query = query.Where(
			"users.name ILIKE ? OR users.email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, int(count), nil
}

func (r *AdminRepositoryImpl) FindByAdminID(ctx context.Context, adminID int) (*models.User, error) {
	var admin models.User
	if err := r.DB.WithContext(ctx).Preload("Role").Joins("JOIN roles ON roles.id = users.role_id").Where("users.id = ?", adminID).Where("LOWER(roles.name) = LOWER(?)", "ADMIN").First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepositoryImpl) Update(ctx context.Context, admin *models.User) error {
	return r.DB.WithContext(ctx).Save(admin).Error
}
