package productrequest

import "mime/multipart"

type CreateProductRequest struct {
	CategoryID  int                   `json:"category_id" form:"category_id"`
	Name        string                `json:"name" form:"name"`
	Description string                `json:"description" form:"description"`
	Price       int                   `json:"price" form:"price"`
	Stock       int                   `json:"stock" form:"stock"`
	ImageFile   *multipart.FileHeader `json:"image_file" form:"image_file"`
}

type UpdateProductRequest struct {
	CategoryID  int    `json:"category_id" form:"category_id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Price       int    `json:"price" form:"price"`
	Stock       int    `json:"stock" form:"stock"`
	IsActive    *bool  `json:"is_active" form:"is_active"`
}

type UpdateProductImageRequest struct {
	ImageFile *multipart.FileHeader `json:"image_file" form:"image_file"`
}
