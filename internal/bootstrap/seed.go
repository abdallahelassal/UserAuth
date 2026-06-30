package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository"
	"github.com/abdallahelassal/UserAuth/pkg/bcrypt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) {
	perms := []domain.Permission{
		{Name: "roles:view"},
		{Name: "roles:create"},
		{Name: "roles:delete"},
		{Name: "permissions:manage"},
	}

	for _, p := range perms {
		db.Where("name = ?", p.Name).FirstOrCreate(&p)
	}
}

func SeedRoles(db *gorm.DB) {
	role := domain.Role{Name: "admin"}
	db.Where("name = ?", role.Name).FirstOrCreate(&role)
}

func seedRoles(repo repository.RoleRepository) {
	roles := []string{"admin", "user"}
	ctx := context.Background()
	
	for _, r := range roles {
		_, err := repo.FindByName(ctx, r)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			repo.Create(ctx, &domain.Role{
				Base: domain.Base{ID: uuid.New()},
				Name: r,
			})
		}
	}
}

func SeedRolePermissions(db *gorm.DB) {
	var admin domain.Role
	var perm domain.Permission

	db.Where("name = ?", "admin").First(&admin)
	db.Where("name = ?", "permissions:manage").First(&perm)

	if admin.ID == uuid.Nil || perm.ID == uuid.Nil {
		return
	}

	db.Model(&admin).Association("Permissions").Append(&perm)
}

func SeedUsers(db *gorm.DB) {
	var user domain.User

	err := db.Where("email = ?", "admin@test.com").First(&user).Error

	hashed, _ := bcrypt.HashPassword("123456")

	if err != nil {
		// create
		user = domain.User{
			UserName: "admin",
			Email:    "admin@test.com",
			Password: string(hashed),
		}

		if err := db.Create(&user).Error; err != nil {
			panic(err)
		}

		fmt.Println("✅ admin user created")

	} else {
		// 🔥 update password لو موجود
		db.Model(&user).Update("password", string(hashed))

		fmt.Println("♻️ admin password updated")
	}
}
func SeedUserRoles(db *gorm.DB) {
	var user domain.User
	var role domain.Role

	db.Where("email = ?", "admin@test.com").First(&user)
	db.Where("name = ?", "admin").First(&role)

	if user.ID == uuid.Nil || role.ID == uuid.Nil {
		return
	}

	db.Model(&user).Association("Roles").Append(&role)
}