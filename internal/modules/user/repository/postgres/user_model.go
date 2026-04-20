package postgres

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type UserModel struct {
	UUID		uuid.UUID	`gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName	string		`gorm:"not null"`
	Email		string		`gorm:"unique;not null"`
	Password	string		`gorm:"not null"`
	IsActive	bool		`gorm:"default:true"`
	CreatedAt 	time.Time 	`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time 	`gorm:"autoUpdateTime"`
}

func (UserModel) TableName()string{
	return "users"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	id , err := uuid.NewV6()
	if err != nil {
		return err
	}

	u.UUID = id
	return nil
}

func (u *UserModel)	FromDomain()*domain.User{
	return &domain.User{
		Base: domain.Base{
			UUID: 		u.UUID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
	}
}

func  ToDomain(u *domain.User) *UserModel{
	return &UserModel{
		UUID: u.UUID,
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
		IsActive: u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}