package roleservice

import (
	"context"
	rolerequest "marketly-app/internal/dto/request/role_request"
	"marketly-app/internal/models"
)

type IRoleService interface {
	CreateRole(ctx context.Context, req rolerequest.CreateRoleRequest) error
	GetAllRole(ctx context.Context, page, limit int, search string) ([]*models.Role, int, error)
	GetByIdRole(ctx context.Context, roleId int) (*models.Role, error)
	UpdateRole(ctx context.Context, roleId int, req rolerequest.UpdateRoleRequest) error
	DeleteRole(ctx context.Context, roleId int) error
}
