package userservice

import (
	"context"
	userrequest "marketly-app/internal/dto/request/user_request"
	"marketly-app/internal/models"
)

type IUserService interface {
	Register(ctx context.Context, req userrequest.RegisterUserRequest) error
	Login(ctx context.Context, req userrequest.LoginUserRequest) (string, error)
	GetProfile(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, req userrequest.UpdateUserRequest) error
	UpdatePhoto(ctx context.Context, userID int, req userrequest.UpdatePhotoUserRequest) error
	ChangeEmail(ctx context.Context, userID int, req userrequest.ChangeEmailUserRequest) error
	ChangePassword(ctx context.Context, userID int, req userrequest.ChangePasswordUserRequest) error
	Logout(ctx context.Context, userID int) error
}
