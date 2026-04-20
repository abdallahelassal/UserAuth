package postgres

import (
	"context"

	"github.com/abdallahelassal/UserAuth/domain"
)

type UserRepository interface {
	Fetch(ctx context.Context,cursor string,limit int)(*[]domain.User,error)
	CreateUser(ctx context.Context,user *domain.User) error
	GetByEmail(ctx context.Context,email string)(*domain.User,error)
	Update(ctx context.Context,uuid string) error
	Delete(ctx context.Context,uuid string)error
}