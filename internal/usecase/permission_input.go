package usecase

import "github.com/google/uuid"



type PermissionInput struct{
	Name	string
}

type UpdatePermissionInput struct {
	Name string
}

type GetUserPermissions struct{
	Name	string
	UserID 		uuid.UUID
}

type GetPermissionsByRoleIDs struct{
	Name string
	RoleID uuid.UUID
}

type PermissionOutput struct {
	ID   string
	Name string
}