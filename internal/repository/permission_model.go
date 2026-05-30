package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Permission struct{
	ID 			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Name 		string		`gorm:"not null;"`
	CreatedAt 	time.Time 	`gorm:"autoCreateTime"`
	UpdatedAt	time.Time 	`gorm:"autoUpdateTime"`

	Roles 		[]*Role 	`gorm:"many2many:role_has_permissions;constraint:OnDelete:CASCADE;"`
	Users		[]*User		`gorm:"many2many:user_has_permissions;constraint:OnDelete:CASCADE;"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB)error{
	if p.ID == uuid.Nil{
		p.ID = uuid.New()

	}
	return nil
}

func (p *Permission) ToDomainPermission()*domain.Permission{
	return &domain.Permission{
		Base: domain.Base{
			ID: p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		Name: p.Name,

	}
}

func FromDomainPermission(permission *domain.Permission )*Permission{
	return&Permission{
		ID: permission.ID,
		Name: permission.Name,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

