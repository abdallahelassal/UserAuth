package usecase

import (
	"context"
	
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository"

	"github.com/google/uuid"
)


type RoleUseCase interface {
    Create(ctx context.Context, req RoleCreateInput) error
    Update(ctx context.Context, req RoleUpdateInput) error

    FindByID(ctx context.Context, id uuid.UUID) (RoleOutput, error)
    FindAll(ctx context.Context) ([]RoleOutput, error)

    GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]RoleOutput, error)

    Delete(ctx context.Context, id uuid.UUID) error
}

type roleUseCase struct{
    RoleRepo    repository.RoleRepository
    ContextTimeout  time.Duration 
}

func NewRoleUseCase(roleRepo repository.RoleRepository,timeOut time.Duration) RoleUseCase{
    return &roleUseCase{
        RoleRepo: roleRepo,
        ContextTimeout: timeOut,
    }
}

func (r *roleUseCase) Create(ctx context.Context,req RoleCreateInput)error{
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()
    
    if req.Name == "" {
        return domain.ErrRoleNameRequired
    }
    role := &domain.Role{
        Name: req.Name,
    }
    if err := r.RoleRepo.Create(ctx,role);err != nil {
        return err
    }
    if len(req.PermissionIDs) > 0 {
        if err := r.RoleRepo.AssignPermission(ctx, role.ID, req.PermissionIDs); err != nil {
            return err
        }
    }
    return nil
}

func (r *roleUseCase) Update(ctx context.Context,req RoleUpdateInput)error{
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()

    if req.Name == "" {
        return domain.ErrRoleNameRequired
    }
    role := &domain.Role{
        Base: domain.Base{
            ID: req.ID,
        },
        Name: req.Name,
    }
    if err := r.RoleRepo.Update(ctx,role);err != nil {
        return err
    }
    return nil
}
func (r *roleUseCase) FindByID(ctx context.Context, id uuid.UUID)(RoleOutput,error){
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()

    role, err := r.RoleRepo.FindByID(ctx, id)
    if err != nil {
        return RoleOutput{}, err
    }
    return RoleOutput{
        ID: role.ID,
        Name: role.Name,
    }, nil
}


func (r *roleUseCase) FindAll(ctx context.Context)([]RoleOutput,error){
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()

    roles, err := r.RoleRepo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    output := make([]RoleOutput,len(roles))
    for i , v := range roles{
        output[i] = RoleOutput{
            ID: v.ID,
            Name: v.Name,
        }
    }
    
    return output, nil
}

func (r *roleUseCase) GetRolesByUserID(ctx context.Context, userID uuid.UUID)([]RoleOutput,error){
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()

    roles, err := r.RoleRepo.GetRolesByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    output := make([]RoleOutput,len(roles)) 
    for i, v := range roles {
        output[i] =  RoleOutput{
            ID: v.ID,
            Name: v.Name,
        }
    }
    return output, nil
}

func (r *roleUseCase) Delete(ctx context.Context, id uuid.UUID)error{
    ctx , cancel := context.WithTimeout(ctx, r.ContextTimeout)
    defer cancel()

    return r.RoleRepo.Delete(ctx, id)
}       
