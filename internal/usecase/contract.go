package usecase

import (
	"context"

	"github.com/google/uuid"
)


type UserUsecase interface{
	Signup(ctx context.Context,req CreateUserInput)error
	FindByID(ctx context.Context,userID uuid.UUID)(*FindByIDOutput,error)
	GetByEmail(ctx context.Context, email string)(*UserOutput,error)
	GetByName(ctx context.Context,name string)(*UserOutput,error)
	Login(ctx context.Context,req LoginUserInput)(string,error)
	AssignRole(ctx context.Context,userID uuid.UUID,roleID uuid.UUID)error
	GetFullProfile(ctx context.Context,userID uuid.UUID)(*FullProfile,error)
}