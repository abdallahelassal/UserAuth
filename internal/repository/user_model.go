package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type User struct {
	ID		uuid.UUID	`gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName	string		`gorm:"not null"`
	Email		string		`gorm:"unique;not null"`
	Password	string		`gorm:"not null"`
	IsActive	bool		`gorm:"default:true"`
	CreatedAt 	time.Time 	`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time 	`gorm:"autoUpdateTime"`

	Roles 			[]*Role						`gorm:"many2many:user_has_roles;constraint:OnDelete:CASCADE;"`
	Permissions 	[]*Permission				`gorm:"many2many:user_has_permissions;constraint:OnDelete:CASCADE;"`
	Tokens 			[]*PersonalAccessToken		`gorm:"foreignKey:UserID"`
}

func (User) TableName()string{
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User)	ToDomain()*domain.User{
	return &domain.User{
		Base: domain.Base{
			ID: 		u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
	}
}

func FromDomain(u *domain.User) *User{
	return &User{
		ID: u.ID,
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}