package repository

import (
	"context"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct{
	db 	*gorm.DB
}

func NewPermissionRepository(db *gorm.DB)*permissionRepository{
	return &permissionRepository{
		db: db,
	}
}


func (r *permissionRepository) FindAllPermissions(ctx context.Context)([]domain.Permission,error){
	var dbModels []Permission
	if err := r.db.WithContext(ctx).Find(&dbModels).Error;err != nil {
		return nil,err
	}
	permissions := make([]domain.Permission,len(dbModels))

	for i , v := range dbModels{
		permissions[i] = *v.ToDomainPermission()
	}
	return permissions,nil
}
func (r *permissionRepository) GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]domain.Permission,error){
	var user User
	if err := r.db.WithContext(ctx).Preload("permissions").First(&user,"id = ?", userID).Error; err != nil {
		return nil , err
	}
	permissions := make([]domain.Permission,len(user.Permissions))
	for i , v := range user.Permissions {
		permissions[i] = *v.ToDomainPermission()
	}
	return  permissions , nil 
}

func (r  *permissionRepository) GetPermissionByRoleIDs(ctx context.Context,roleIDs []uuid.UUID)([]domain.Permission,error){
	var role Role

	if err := r.db.WithContext(ctx).Preload("permissions").First(&role,"id IN ?", roleIDs).Error; err != nil {
		return nil , err
	}

	permissions := make([]domain.Permission,len(role.Permissions))

	for i , v := range role.Permissions{
		permissions[i] = *v.ToDomainPermission()
	}
	return permissions , nil
}
