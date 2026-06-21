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
	permissionRepo domain.PermissionRepository
	contextTimeout time.Duration
	jwtSecret     string
	jwtExpiry     time.Duration
}

func NewUserUseCase(userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	permissionRepo domain.PermissionRepository,
	timeout time.Duration,
	jwtSecret string,
	jwtExpiary time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		permissionRepo: permissionRepo,
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
	if user == nil || user.ID == uuid.Nil {
	return errors.New("invalid user")
}

if role == nil || role.ID == uuid.Nil {
	return errors.New("invalid role")
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

func (uc *UserUseCase) GetFullProfile(ctx context.Context, userID uuid.UUID) (*FullProfile, error) {

	// 1. user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. roles
	roles, err := uc.roleRepo.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. permissions from roles
	roleIDs := make([]uuid.UUID, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	rolePermissions, err := uc.permissionRepo.GetPermissionByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	// 4. direct permissions
	directPermissions, _ := uc.permissionRepo.GetPermissionsByUserID(ctx, userID)

	// 5. merge (simple version)
	allPermissions := append(rolePermissions, directPermissions...)

	return &FullProfile{
		User:        ToUserOutput(*user),
		Roles:       ToRoleOutputs(roles),
		Permissions: ToPermissionOutputs(allPermissions),
	}, nil
}
func ToUserOutput(u domain.User) UserOutput {
    return UserOutput{
        ID:       u.ID.String(),
        Email:    u.Email,
        UserName: u.UserName,
    }
}

func ToRoleOutputs(roles []domain.Role) []RoleOutput {
	out := make([]RoleOutput, len(roles))

	for i, r := range roles {
		out[i] = ToRoleOutput(r)
	}

	return out
}
func ToRoleOutput(r domain.Role) RoleOutput {
	return RoleOutput{
		ID:   r.Base.ID,
		Name: r.Name,
	}
}
func ToPermissionOutput(p domain.Permission) PermissionOutput {
	return PermissionOutput{
		ID:   p.Base.ID.String(),
		Name: p.Name,
	}
}
func ToPermissionOutputs(perms []domain.Permission) []PermissionOutput {
	out := make([]PermissionOutput, len(perms))

	for i, p := range perms {
		out[i] = ToPermissionOutput(p)
	}

	return out
}