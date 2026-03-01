package categoryrequest

type CreateCategoryRequest struct {
	Name string `json:"name" form:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" form:"name"`
}
