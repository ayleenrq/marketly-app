package productresponse

import (
	"marketly-app/internal/models"
	"marketly-app/pkg/utils"
)

type SellerInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Price       int          `json:"price"`
	Stock       int          `json:"stock"`
	ImageURL    *string      `json:"image_url"`
	IsActive    bool         `json:"is_active"`
	Seller      SellerInfo   `json:"seller"`
	Category    CategoryInfo `json:"category"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

func ToProductResponse(product models.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		IsActive:    product.IsActive,
		Seller: SellerInfo{
			ID:   product.SellerID,
			Name: product.User.Name,
		},
		Category: CategoryInfo{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		},
		CreatedAt: utils.FormatDate(product.CreatedAt),
		UpdatedAt: utils.FormatDate(product.UpdatedAt),
	}
}
