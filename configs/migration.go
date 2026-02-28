package configs

import (
	"fmt"
	"marketly-app/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	fmt.Println("Running migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)

	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	fmt.Println("Migrations completed!")
}
