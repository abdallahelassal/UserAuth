package seed

import (

	"gorm.io/gorm"
	"github.com/abdallahelassal/UserAuth/internal/repository"
)
func SeedRolePermissions(
	db *gorm.DB,
	roles map[string]repository.Role,
	permissions map[string]repository.Permission,
) {

	adminRole := roles["admin"]

	adminPermissions := []string{
		"roles:view",
		"roles:create",
		"roles:delete",
		"permissions:manage",
	}

	for _, permName := range adminPermissions {

		var exists int64

		db.Table("role_has_permission").
			Where("role_id = ? AND permission_id = ?", adminRole.ID, permissions[permName].ID).
			Count(&exists)

		if exists == 0 {
			db.Exec(`
				INSERT INTO role_has_permission (role_id, permission_id)
				VALUES (?, ?)`,
				adminRole.ID, permissions[permName].ID,
			)
		}
	}
}