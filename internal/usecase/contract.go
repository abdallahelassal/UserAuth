package usecase

import (
	"context"

	"github.com/google/uuid"
)


type UserUsecase interface{
	Signup(ctx context.Context,req CreateUserInput)error
	FindByID(ctx context.Context,userID uuid.UUID)(*FindByIDOutput,error)
	GetByEmail(ctx context.Context, email string)(*UserOutput,error)
	GetByName(ctx context.Context,name string)(*UserOutput,error)
	Login(ctx context.Context,req LoginUserInput)(string,error)
	AssignRole(ctx context.Context,userID uuid.UUID,roleID uuid.UUID)error
	GetFullProfile(ctx context.Context,userID uuid.UUID)(*FullProfile,error)
}
type RoleUsecase interface {
	Create(ctx context.Context, req RoleCreateInput)error
	Update(ctx context.Context, req RoleUpdateInput)error
	FindByID(ctx context.Context,roleID uuid.UUID)(*RoleOutput,error)
	FindAll(ctx context.Context)([]*RoleOutput,error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetRolesByUserID(ctx context.Context,userID uuid.UUID)([]*RoleOutput,error)
}
type PermissionUsecase interface{
	FindAllPermissions(ctx context.Context)([]PermissionOutput,error)
	GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]GetUserPermissions,error)
	GetPermissionByRoleIDs(ctx context.Context,roleIDs []uuid.UUID)([]GetPermissionsByRoleIDs,error)
	Create(ctx context.Context,perm *PermissionInput)error
}