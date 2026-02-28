package response

import (
	"marketly-app/internal/models"
	"marketly-app/pkg/utils"
)

type SellerResponse struct {
	ID               int    `json:"id"`
	Role             string `json:"role"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Address          string `json:"address"`
	PhotoURL         string `json:"photo_url"`
	StoreName        string `json:"store_name,omitempty"`
	StoreDescription string `json:"store_description,omitempty"`
	IsVerified       bool   `json:"is_verified,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func ToSellerResponse(seller models.User) SellerResponse {
	return SellerResponse{
		ID:               seller.ID,
		Role:             seller.Role.Name,
		Username:         seller.Username,
		Name:             seller.Name,
		Email:            seller.Email,
		PhoneNumber:      derefStr(seller.PhoneNumber),
		Address:          derefStr(seller.Address),
		PhotoURL:         derefStr(seller.PhotoURL),
		StoreName:        derefStr(seller.StoreName),
		StoreDescription: derefStr(seller.StoreDescription),
		IsVerified:       derefBool(seller.IsVerified),
		CreatedAt:        utils.FormatDate(seller.CreatedAt),
		UpdatedAt:        utils.FormatDate(seller.UpdatedAt),
	}
}
