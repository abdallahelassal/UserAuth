package domain

import (
	"context"

	
)

type UserRepository interface {
	Fetch(ctx context.Context,cursor string,limit int)(*[]User,error)
	CreateUser(ctx context.Context,user *User) error
	GetByEmail(ctx context.Context,email string)(*User,error)
	Update(ctx context.Context,uuid string) error
	Delete(ctx context.Context,uuid string)error
}