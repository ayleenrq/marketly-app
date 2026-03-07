package productresponse

import "marketly-app/internal/models"

type ProductResponse struct {
	ID         int     `json:"id"`
	SellerID   int     `json:"seller_id"`
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name"`
	Price      int     `json:"price"`
	Stock      int     `json:"stock"`
	ImageURL   *string `json:"image_url"`
}

func ToProductResponse(product models.Product) ProductResponse {
	return ProductResponse{
		ID:         product.ID,
		SellerID:   product.SellerID,
		CategoryID: product.CategoryID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		ImageURL:   product.ImageURL,
	}
}