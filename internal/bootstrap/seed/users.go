package seed

import (
	"log"

	"github.com/abdallahelassal/UserAuth/internal/repository"
	"github.com/abdallahelassal/UserAuth/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
func SeedUsers(db *gorm.DB) map[string]repository.User {

	result := map[string]repository.User{}

	var user repository.User
	email := "admin@test.com"

	err := db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {

		hashedPassword, _ := bcrypt.HashPassword("123456")

		user = repository.User{
			ID:       uuid.New(),
			UserName: "admin",
			Email:    email,
			Password: hashedPassword,
			IsActive: true,
		}

		db.Create(&user)
		log.Println("Inserted admin user")
	}

	result["admin"] = user

	return result
}