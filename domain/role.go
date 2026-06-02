package domain

import (
	"context"

	
	"github.com/google/uuid"
)




type Role struct {
	Base
	Name 		string		`json:"name"`
	
}


type RoleRepository interface{
	Create(ctx context.Context,r *Role)error
	AssignPermission(ctx context.Context,roleID uuid.UUID,permID []uuid.UUID)error
	RemoveAllPermission(ctx context.Context,roleID uuid.UUID)error
	FindByID(ctx context.Context,ID uuid.UUID)(*Role,error)
	FindAll(ctx context.Context)([]Role,error)
	GetRolesByUserID(ctx context.Context,userID uuid.UUID)([]Role,error)
	Update(ctx context.Context,r *Role)error
	Delete(ctx context.Context,id uuid.UUID)error
	FindByName(ctx context.Context,name string)(*Role,error)
}

