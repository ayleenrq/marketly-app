package categoryroute

import (
	categoryhandler "marketly-app/internal/handlers/category_handler"
	categoryrepository "marketly-app/internal/repositories/category_repository"
	categoryservice "marketly-app/internal/services/category_service"
	middleware "marketly-app/internal/middlewares"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CategoryRoutes(e *echo.Group, db *gorm.DB) {
	categoryRepo := categoryrepository.NewCategoryRepositoryImpl(db)
	categoryService := categoryservice.NewCategoryServiceImpl(categoryRepo)
	categoryHandler := categoryhandler.NewCategoryHandler(categoryService)

	e.POST("/create", categoryHandler.CreateCategory, middleware.AllowRoles("admin"))
    e.GET("/all", categoryHandler.GetAllCategory)
    e.GET("/:categoryId", categoryHandler.GetByIdCategory)
    e.PUT("/:categoryId/edit", categoryHandler.UpdateCategory, middleware.AllowRoles("admin"))
    e.DELETE("/:categoryId/delete", categoryHandler.DeleteCategory, middleware.AllowRoles("admin"))
}
