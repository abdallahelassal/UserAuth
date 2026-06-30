package repository

import (
	"context"
	"errors"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type RoleRepository struct{
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB)*RoleRepository{
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(ctx context.Context, model *domain.Role)error{
	dbModel := FromDomainRole(model)

	if err := r.db.WithContext(ctx).Create(&dbModel).Error; err != nil {
		return err
	}
	model.ID = dbModel.ID
	return nil
}

func (r *RoleRepository) AssignPermission(ctx context.Context,roleID uuid.UUID,permID []uuid.UUID)error{
	var role Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}
	var permissions []Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", permID).Find(&permissions).Error; err != nil {
		return err
	}
	if len(permissions) == 0 {
		return errors.New("no permissions found with the provided IDs")
	}
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Replace(permissions)
}

func (r *RoleRepository) FindByID(ctx context.Context,roleID uuid.UUID)(*domain.Role,error){
	var dbModel Role

	if err := r.db.WithContext(ctx).Where("id = ?", roleID).Take(&dbModel).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
		return nil , domain.ErrNotFound
		}
		return nil ,err
	}
	

	return dbModel.ToDomainRole(),nil

}

func (r *RoleRepository) FindByName(ctx context.Context,name string)(*domain.Role,error){
	var dbModel Role

	if err := r.db.WithContext(ctx).Where("name = ?", name).Take(&dbModel).Error; errors.Is(err,gorm.ErrRecordNotFound){
		return nil , nil
	}

	return dbModel.ToDomainRole(),nil

}

func (r *RoleRepository) FindAll(ctx context.Context)([]domain.Role,error){

	var dbModels []Role

	err := r.db.WithContext(ctx).Find(&dbModels).Error


	if err != nil {
		return nil,err
	}

	role := make([]domain.Role, len(dbModels))
	for i, dbModel := range dbModels {
		role[i] = *dbModel.ToDomainRole()
	}

	return role, nil

}

func (r *RoleRepository) GetRolesByUserID(ctx context.Context,userID uuid.UUID)([]domain.Role,error){
	var user User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	roles := make([]domain.Role, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = *role.ToDomainRole()
	}

	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, model *domain.Role)error{
	dbModel := FromDomainRole(model)

	if err := r.db.WithContext(ctx).Save(&dbModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) Delete(ctx context.Context,id uuid.UUID)error{
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Role{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}	

func (r *RoleRepository) RemoveAllPermission(ctx context.Context,roleID uuid.UUID)error{
	var role Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Clear()
}