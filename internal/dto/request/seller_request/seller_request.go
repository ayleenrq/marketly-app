package sellerrequest

import "mime/multipart"

type RegisterSellerRequest struct {
	Username         string                `json:"username" form:"username"`
	Name             string                `json:"name" form:"name"`
	Email            string                `json:"email" form:"email"`
	Password         string                `json:"password" form:"password"`
	PhoneNumber      string                `json:"phone_number" form:"phone_number"`
	Address          string                `json:"address" form:"address"`
	PhotoFile        *multipart.FileHeader `json:"photo_file" form:"photo_file"`
	StoreName        string                `json:"store_name" form:"store_name"`
	StoreDescription string                `json:"store_description" form:"store_description"`
}

type LoginSellerRequest struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UpdateSellerRequest struct {
	Username         string                `json:"username" form:"username"`
	Name             string                `json:"name" form:"name"`
	PhoneNumber      string                `json:"phone_number" form:"phone_number"`
	Address          string                `json:"address" form:"address"`
	StoreName        string                `json:"store_name" form:"store_name"`
	StoreDescription string                `json:"store_description" form:"store_description"`
	PhotoFile        *multipart.FileHeader `json:"photo_file" form:"photo_file"`
}

type ChangePasswordSellerRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

type ChangeEmailSellerRequest struct {
	NewEmail string `json:"new_email" form:"new_email"`
	Password string `json:"password" form:"password"`
}

type UpdatePhotoSellerRequest struct {
	PhotoFile *multipart.FileHeader `json:"photo_file" form:"photo_file"`
}
