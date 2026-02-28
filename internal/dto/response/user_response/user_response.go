	package response

import (
	"marketly-app/internal/models"
	"marketly-app/pkg/utils"
)

type UserResponse struct {
	ID          int    `json:"id"`
	Role        string `json:"role"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	PhotoURL    string `json:"photo_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Role:        user.Role.Name,
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: derefStr(user.PhoneNumber),
		Address:     derefStr(user.Address),
		PhotoURL:    derefStr(user.PhotoURL),
		CreatedAt:   utils.FormatDate(user.CreatedAt),
		UpdatedAt:   utils.FormatDate(user.UpdatedAt),
	}
}
