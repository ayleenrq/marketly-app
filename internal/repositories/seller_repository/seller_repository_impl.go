package sellerrepository

import (
	"context"
	"marketly-app/internal/models"

	"gorm.io/gorm"
)

type SellerRepositoryImpl struct {
	DB *gorm.DB
}

func NewSellerRepositoryImpl(db *gorm.DB) ISellerRepository {
	return &SellerRepositoryImpl{DB: db}
}

func (r *SellerRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *SellerRepositoryImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Role").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.username = ?", username).
		Where("LOWER(roles.name) = LOWER(?)", "SELLER").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SellerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Role").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.email = ?", email).
		Where("LOWER(roles.name) = LOWER(?)", "SELLER").
		First(&user).Error; err != nil {
			return nil, err
	}
	return &user, nil
}

func (r *SellerRepositoryImpl) FindRoleSeller(ctx context.Context) (*models.Role, error) {
	var role models.Role
	if err := r.DB.WithContext(ctx).Where("LOWER(name) = LOWER(?)", "seller").First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *SellerRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*models.User, int64, error) {
	var (
		users []*models.User
		count int64
	)

	query := r.DB.WithContext(ctx).Model(&models.User{}).Preload("Role").Joins("JOIN roles ON roles.id = users.role_id").Where("LOWER(roles.name) = LOWER(?)", "SELLER")

	if search != "" {
		query = query.Where(
			"users.name ILIKE ? OR users.email ILIKE ? OR users.username ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *SellerRepositoryImpl) FindById(ctx context.Context, sellerID int) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Role").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.id = ?", sellerID).
		Where("LOWER(roles.name) = LOWER(?)", "SELLER").
		First(&user).Error; err != nil {
			return nil, err
	}
	return &user, nil
}

func (r *SellerRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}
