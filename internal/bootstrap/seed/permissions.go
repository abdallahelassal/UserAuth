package seed

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/abdallahelassal/UserAuth/internal/repository"
)

func SeedPermissions(db *gorm.DB) map[string]repository.Permission {

	permissions := []string{
		"roles:view",
		"roles:create",
		"roles:delete",
		"permissions:manage",
	}

	result := map[string]repository.Permission{}

	for _, name := range permissions {
		var perm repository.Permission

		err := db.Where("name = ?", name).First(&perm).Error
		if err == gorm.ErrRecordNotFound {

			perm = repository.Permission{
				ID:   uuid.New(),
				Name: name,
			}

			db.Create(&perm)
			log.Println("Inserted permission:", name)
		}

		result[name] = perm
	}

	return result
}