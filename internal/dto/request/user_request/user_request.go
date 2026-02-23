package userrequest

import "mime/multipart"

type RegisterUserRequest struct {
	Username    string                `json:"username" form:"username"`
	Name        string                `json:"name" form:"name"`
	Email       string                `json:"email" form:"email"`
	Password    string                `json:"password" form:"password"`
	PhoneNumber string                `json:"phone_number" form:"phone_number"`
	Address     string                `json:"address" form:"address"`
	Role        string                `json:"role" form:"role"`
	PhotoFile   *multipart.FileHeader `json:"photo_file" form:"photo_file"`
}

type LoginUserRequest struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UpdateUserRequest struct {
	Username    string                `json:"username" form:"username"`
	Name        string                `json:"name" form:"name"`
	PhoneNumber string                `json:"phone_number" form:"phone_number"`
	Address     string                `json:"address" form:"address"`
	PhotoFile   *multipart.FileHeader `json:"photo_file" form:"photo_file"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" form:"new_email"`
	Password string `json:"password" form:"password"`
}

type UpdatePhotoRequest struct {
	PhotoFile *multipart.FileHeader `json:"photo_file" form:"photo_file"`
}
