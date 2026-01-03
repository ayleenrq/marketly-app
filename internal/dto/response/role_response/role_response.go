package response

import (
	"marketly-app/internal/models"
	"marketly-app/pkg/utils"
)

type RoleResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToRoleResponse(role models.Role) RoleResponse {
	return RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: utils.FormatDate(role.CreatedAt),
		UpdatedAt: utils.FormatDate(role.UpdatedAt),
	}
}
