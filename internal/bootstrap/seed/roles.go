package seed

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/abdallahelassal/UserAuth/internal/repository"
)
func SeedRoles(db *gorm.DB) map[string]repository.Role {

	roles := []string{"admin", "user"}

	result := map[string]repository.Role{}

	for _, name := range roles {
		var role repository.Role

		err := db.Where("name = ?", name).First(&role).Error
		if err == gorm.ErrRecordNotFound {

			role = repository.Role{
				ID:   uuid.New(),
				Name: name,
			}

			db.Create(&role)
			log.Println("Inserted role:", name)
		}

		result[name] = role
	}

	return result
}