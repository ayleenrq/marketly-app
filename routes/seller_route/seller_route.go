package sellerroute

import (
	sellerhandler "marketly-app/internal/handlers/seller_handler"
	sellerrepository "marketly-app/internal/repositories/seller_repository"
	sellerservice "marketly-app/internal/services/seller_service"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SellerRoutes(e *echo.Group, db *gorm.DB, cld *cloudinary.Cloudinary) {
	sellerRepo := sellerrepository.NewSellerRepositoryImpl(db)
	sellerService := sellerservice.NewSellerServiceImpl(sellerRepo, cld)
	sellerHandler := sellerhandler.NewSellerHandler(sellerService)

	e.POST("/register", sellerHandler.RegisterSeller)
	e.POST("/login", sellerHandler.LoginSeller)

	auth := e.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	
	auth.GET("/me", sellerHandler.GetProfileSeller)
	auth.PUT("/me", sellerHandler.UpdateProfileSeller)
	auth.PUT("/me/photo", sellerHandler.UpdatePhotoProfile)
	auth.PUT("/me/password", sellerHandler.ChangePassword)
	auth.PUT("/me/email", sellerHandler.ChangeEmail)
	auth.POST("/logout", sellerHandler.LogoutSeller)
}
