package domain

import (
	"context"

	"github.com/google/uuid"
)


type Permission struct{
	Base
	Name string `json:"name"`
}

type PermissionRepository interface{
	FindAllPermissions(ctx context.Context)([]Permission,error)
	GetPermissionsByUserID(ctx context.Context,userID uuid.UUID)([]Permission,error)
	GerPermissionsByRolesID(ctx context.Context,roleIS uuid.UUID)([]Permission,error)
}