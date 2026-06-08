package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"

	"github.com/google/uuid"
)


type PermissionUsecase interface{
	FindAllPermissions(ctx context.Context)([]GetPermissions,error)
	GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]GetUserPermissions,error)
	GetPermissionByRoleIDs(ctx context.Context,roleIDs []uuid.UUID)([]GetPermissionsByRoleIDs,error)
}

type permissionUsecase struct{
	PremissionRepo domain.PermissionRepository
	ContextTimeOut time.Duration
}

func NewPermissionUsecase(permissionRepo domain.PermissionRepository, time time.Duration)*permissionUsecase{
	return &permissionUsecase{
		PremissionRepo: permissionRepo,
		ContextTimeOut: time,
	}
}

func (p *permissionUsecase) FindAllPermissions(ctx context.Context)([]GetPermissions,error){
	ctx , cancel := context.WithTimeout(ctx,p.ContextTimeOut)
	defer cancel()

	permissions , err := p.PremissionRepo.FindAllPermissions(ctx)
	if err != nil {
		return []GetPermissions{},err
	}

	output := make([]GetPermissions,len(permissions))

	for i , v := range permissions{
		output[i] = GetPermissions{
			Name: v.Name,
		}
	}
	return output,nil
}

func (p *permissionUsecase) GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]GetUserPermissions,error){
	ctx , cancel := context.WithTimeout(ctx , p.ContextTimeOut)
	defer cancel()
	
	if userID == uuid.Nil{
		return nil,errors.New("user id not required")
	}

	permissions , err := p.PremissionRepo.GetPermissionsByUserID(ctx,userID)
	if err != nil {
		return nil,err
	}
	output := make([]GetUserPermissions,len(permissions))

	for i , v := range permissions{
		output[i] = GetUserPermissions{
			Name: v.Name,
			UserID: userID,
		}
	}
	return output ,nil
}

func (p *permissionUsecase) GetPermissionByRoleIDs(ctx context.Context,rolesIDs []uuid.UUID)([]GetPermissionsByRoleIDs,error){
	ctx , cancel := context.WithTimeout(ctx , p.ContextTimeOut)
	defer cancel()
	if len(rolesIDs) == 0 {
		return nil , errors.New("roles id is not found")
	}
	
	permissions , err := p.PremissionRepo.GetPermissionByRoleIDs(ctx,rolesIDs)
	if err != nil {
		return nil , errors.New("roles not found")
	}
	
	output := make([]GetPermissionsByRoleIDs,len(permissions))
	for i , v := range permissions{
		output[i] = GetPermissionsByRoleIDs{
			Name: v.Name,
			RoleID: uuid.Nil,
		}
	}
	return output,nil 
}