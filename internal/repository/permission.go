package repository

import (
	"context"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository struct{
	db 	*gorm.DB
}

func NewPermissionRepository(db *gorm.DB)*PermissionRepository{
	return &PermissionRepository{
		db: db,
	}
}


func (r *PermissionRepository) FindAllPermission(ctx context.Context)([]domain.Permission,error){
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
func (r *PermissionRepository) GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]domain.Permission,error){
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

func (r  *PermissionRepository) GetPermissionByRoleID(ctx context.Context,roleID uuid.UUID)([]domain.Permission,error){
	var role Role

	if err := r.db.WithContext(ctx).Preload("permissions").First(&role,"id = ?", roleID).Error; err != nil {
		return nil , err
	}

	permissions := make([]domain.Permission,len(role.Permissions))

	for i , v := range role.Permissions{
		permissions[i] = *v.ToDomainPermission()
	}
	return permissions , nil
}
