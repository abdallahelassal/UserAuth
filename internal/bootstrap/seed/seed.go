
package seed

import (
	

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {

	permissions := SeedPermissions(db)
	roles := SeedRoles(db)
	users := SeedUsers(db)

	SeedRolePermissions(db, roles, permissions)
	SeedUserRoles(db, users, roles)
}