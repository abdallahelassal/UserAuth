package seed



import (
	"log"

	"gorm.io/gorm"
	"github.com/abdallahelassal/UserAuth/internal/repository"
)

func SeedUserRoles(db *gorm.DB, users map[string]repository.User, roles map[string]repository.Role) {

	adminUser := users["admin"]
	adminRole := roles["admin"]

	var count int64

	db.Table("user_has_roles").
		Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).
		Count(&count)

	if count == 0 {
		err := db.Exec(`
			INSERT INTO user_has_roles (user_id, role_id)
			VALUES (?, ?)`,
			adminUser.ID, adminRole.ID,
		).Error

		if err != nil {
			log.Println("❌ failed to link user with role:", err)
			return
		}

		log.Println("✅ linked admin user with admin role")
	}
}