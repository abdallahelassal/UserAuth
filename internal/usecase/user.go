package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"

	"github.com/abdallahelassal/UserAuth/pkg/bcrypt"
	"github.com/abdallahelassal/UserAuth/pkg/jwt"
)

type UserUseCase struct {
	userRepo domain.UserRepository
	roleRepo domain.RoleRepository
	contextTimeout time.Duration
	jwtSecret     string
	jwtExpiry     time.Duration
}

func NewUserUseCase(userRepo domain.UserRepository,roleRepo domain.RoleRepository, timeout time.Duration, jwtSecret string,jwtExpiary time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		contextTimeout: timeout,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiary,
	}
}


func (u *UserUseCase) Signup(ctx context.Context, req CreateUserInput) error {
	ctx , cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()


	if req.Email == "" || req.UserName == "" || req.Password == "" {
		return errors.New("all fields are required")
	}

	hashedPassword ,err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}


	user := &domain.User{
		UserName: req.UserName,
		Email: req.Email,
		Password: hashedPassword,
		IsActive: true,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return err
	}

	role, err := u.roleRepo.FindByName(ctx, "user")	
	if err != nil {
	
		role = &domain.Role{Name: "user"}
		if err := u.roleRepo.Create(ctx, role); err != nil {
			return err
		}
	}
		if err := u.userRepo.AssignRole(ctx, user.ID, role.ID); err != nil {
		return err
	}
	
	return  nil 
}

func (u *UserUseCase) FindByID(ctx context.Context,userID uuid.UUID)(FindByIDOutput,error){
	ctx , cancel := context.WithTimeout(ctx,u.contextTimeout)
	defer cancel()

	if userID == uuid.Nil{
		return FindByIDOutput{}, errors.New("invalid user ID")
	}
	user , err := u.userRepo.FindByID(ctx,userID)
	if err != nil {
		return FindByIDOutput{},nil 
	}

	output := FindByIDOutput{
		Email: user.Email,
		UserName: user.UserName,
		IsActive: user.IsActive,
		
	}
	return output,nil
	
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


func (u *UserUseCase) Login(ctx context.Context , req LoginUserInput) (string,error){
	ctx , cancel := context.WithTimeout(ctx , u.contextTimeout)
	defer cancel()

	if req.Email == "" || req.Password == ""{
		return "",errors.New("all fields are required")
	}
	
	user , err  := u.userRepo.GetByEmail(ctx,req.Email)
	if err != nil {
		return "",err
	}




	err = bcrypt.ComparePassword(req.Password, user.Password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token , err := jwt.CreateAccessToken(user, u.jwtSecret, int(u.jwtExpiry))
	if err != nil {
	return "", err
}
	return token, nil
}
func (u *UserUseCase) AssignRole(ctx context.Context,userID uuid.UUID,roleID uuid.UUID)error{
	ctx , cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if userID == uuid.Nil || roleID == uuid.Nil {
		return errors.New("invalid user ID or role ID")
	}
	return u.userRepo.AssignRole(ctx, userID, roleID)
}