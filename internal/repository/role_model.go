package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Role struct{
	ID 				uuid.UUID 		`gorm:"type:uuid;primaryKey"`
	Name 			string			`gorm:"uniqueIndex;size:20;not null"`
	CreatedAt 		time.Time		`gorm:"autoCreateTime"`
	UpdatedAt 		time.Time		`gorm:"autoUpdateTime"`

	Users 			[]*User			`gorm:"many2many:user_has_roles;constraint:OnDelete:CASCADE;"`
	Permissions 	[]*Permission	`gorm:"many2many:role_has_permissions;constraint:OnDelete:CASCADE;"`
}

func (r *Role) BeforeCreate(tx *gorm.DB)error{
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *Role) ToDomainRole()*domain.Role{
	return &domain.Role{
		Base: domain.Base{
			ID: r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Name: r.Name,
	}
}

func  FromDomainRole(role *domain.Role)*Role{
	return &Role{
		ID: role.ID,
		Name: role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
