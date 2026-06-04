package dtos

import (
	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
)

type RoleCreateRequest struct{
	Name 			string		`json:"name" validate:"required,min=3,max=100"`
	PermissionIDs 	[]uuid.UUID `json:"permission_ids" validate:"required,dive,uuid"`
}

type RoleUpdateRequest struct{
	ID				uuid.UUID	`json:"-"`
	Name			string		`json:"name" validate:"required,min=3,max=100"`
	PermissionIDs 	[]uuid.UUID	`json:"permission_ids" validate:"required,dive,uuid"`
}

type RoleWithUserPermission struct{
	ID 			uuid.UUID		`json:"id"`
	Name 		string			`json:"name"`
	Users 		[]domain.User			`json:"users"`
	Permissions []domain.Permission 	`json:"permission"`
}
