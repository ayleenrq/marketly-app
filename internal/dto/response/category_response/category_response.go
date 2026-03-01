package categoryresponse

import (
	"marketly-app/internal/models"
	"marketly-app/pkg/utils"
)

type CategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToCategoryResponse(category models.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: utils.FormatDate(category.CreatedAt),
		UpdatedAt: utils.FormatDate(category.UpdatedAt),
	}
}
