package roleservice

import (
	"context"
	rolerequest "marketly-app/internal/dto/request/role_request"
	"marketly-app/internal/models"
	rolerepository "marketly-app/internal/repositories/role_repository"
	errorresponse "marketly-app/pkg/constant/error_response"
	"strings"
)

type RoleServiceImpl struct {
	roleRepo rolerepository.IRoleRepository
}

func NewRoleServiceImpl(roleRepo rolerepository.IRoleRepository) IRoleService {
	return &RoleServiceImpl{roleRepo: roleRepo}
}

func (s *RoleServiceImpl) CreateRole(ctx context.Context, req rolerequest.CreateRoleRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Role name is required", 400)
	}

	role := &models.Role{
		Name: req.Name,
	}
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleServiceImpl) GetAllRole(ctx context.Context, page, limit int, search string) ([]*models.Role, int, error) {
	offset := (page - 1) * limit
	return s.roleRepo.FindAll(ctx, limit, offset, search)
}

func (s *RoleServiceImpl) GetByIdRole(ctx context.Context, roleId int) (*models.Role, error) {
	return s.roleRepo.FindById(ctx, roleId)
}

func (s *RoleServiceImpl) UpdateRole(ctx context.Context, roleId int, req rolerequest.UpdateRoleRequest) error {
	role, err := s.roleRepo.FindById(ctx, roleId)
	if err != nil {
		return err
	}

	role.Name = req.Name
	return s.roleRepo.Update(ctx, roleId, role)
}

func (s *RoleServiceImpl) DeleteRole(ctx context.Context, roleId int) error {
	return s.roleRepo.Delete(ctx, roleId)
}
