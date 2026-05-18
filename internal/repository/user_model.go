package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type UserModel struct {
	ID		uuid.UUID	`gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName	string		`gorm:"not null"`
	Email		string		`gorm:"unique;not null"`
	Password	string		`gorm:"not null"`
	IsActive	bool		`gorm:"default:true"`
	CreatedAt 	time.Time 	`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time 	`gorm:"autoUpdateTime"`

	Tokens 		[]PersonalAccessToken	`gorm:"foreignKey:userID"`
}

func (UserModel) TableName()string{
	return "users"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *UserModel)	ToDomain()*domain.User{
	return &domain.User{
		Base: domain.Base{
			UUID: 		u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
	}
}

func FromDomain(u *domain.User) *UserModel{
	return &UserModel{
		ID: u.UUID,
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}