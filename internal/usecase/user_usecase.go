package usecase

import (
	"context"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/pkg/bcrypt"
)

type UserUseCase struct {
	userRepo domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(userRepo domain.UserRepository, timeout time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		contextTimeout: timeout,
	}
}


func (u *UserUseCase) Signup(ctx context.Context, user *domain.User) error {
	ctx , cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	
	hashPassword , err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword
	user.IsActive = true
	return u.userRepo.Create(ctx, user)	
}

func (u *UserUseCase) GetByEmail(ctx context.Context, email string)(*domain.User,error){
	ctx , cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.GetByEmail(ctx, email)
}

func (u *UserUseCase) GetByName(ctx context.Context, name string)(*domain.User,error){
	ctx , cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.GetByName(ctx, name)
}
