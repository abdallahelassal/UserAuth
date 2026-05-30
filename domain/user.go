package domain

import (
	"time"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Base struct {
	ID		uuid.UUID	`json:"uuid"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type User struct {
	Base
	Username     string 	`json:"username"`
	Email    	string 		`json:"email"`
	Password 	string 		`json:"password"`
	IsActive 	bool			`json:"is_active"`
}

func (b *Base) BeforeCreate(tx *gorm.DB)(err error){
	id , err := uuid.NewV6()
	if err != nil {
		return err
	}
	b.ID = id

	return
}

type UserRepository interface {
	//Fetch(ctx context.Context,cursor string,limit int)(*[]User,error)
	Create(ctx context.Context,user *User) error
	GetByEmail(ctx context.Context,email string)(*User,error)
	GetByName(ctx context.Context,name string)(*User,error)
	//Update(ctx context.Context,uuid string) error
	//Delete(ctx context.Context,uuid string)error
}