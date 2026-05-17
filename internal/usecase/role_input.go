package usecase

import "github.com/google/uuid"



type RoleCreateInput struct{
	Name 			string		`json:"name"`
	PermissionIDs 	[]uuid.UUID `json:"permission_ids"`
}

type RoleUpdateInput struct{
	ID				uuid.UUID	`json:"-"`
	Name			string		`json:"name" `
	PermissionIDs 	[]uuid.UUID	`json:"permission_ids"`
	
}


type RoleOutput struct {
    ID   uuid.UUID
    Name string
}