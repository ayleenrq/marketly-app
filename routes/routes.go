package routes

import (
	admin_route "marketly-app/routes/admin_route"
	category_route "marketly-app/routes/category_route"
	role_route "marketly-app/routes/role_route"
	seller_route "marketly-app/routes/seller_route"
	user_route "marketly-app/routes/user_route"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, cld *cloudinary.Cloudinary) {
	v1 := e.Group("/api/v1")

	role_route.RoleRoutes(v1.Group("/role"), db)
	admin_route.AdminRoutes(v1.Group("/admin"), db)
	user_route.UserRoutes(v1.Group("/user"), db, cld)
	seller_route.SellerRoutes(v1.Group("/seller"), db, cld)
	category_route.CategoryRoutes(v1.Group("/category"), db)
}
