package user_route

import (
	userhandler "marketly-app/internal/handlers/user_handler"
	userrepository "marketly-app/internal/repositories/user_repository"
	userservice "marketly-app/internal/services/user_service"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRoutes(e *echo.Group, db *gorm.DB, cld *cloudinary.Cloudinary) {
	userRepo := userrepository.NewUserRepositoryImpl(db)
	userService := userservice.NewUserServiceImpl(userRepo, cld)
	userHandler := userhandler.NewUserHandler(userService)

	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)

	auth := e.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	auth.GET("/me", userHandler.GetProfileUser)
	auth.PUT("/me", userHandler.UpdateProfileUser)
 	auth.POST("/logout", userHandler.LogoutUser)
}
