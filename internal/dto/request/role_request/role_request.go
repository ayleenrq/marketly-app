package rolerequest

type CreateRoleRequest struct {
	Name string `json:"name" form:"name"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" form:"name"`
}
