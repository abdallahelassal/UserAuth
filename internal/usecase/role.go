
package usecase

import (
	"context"

	
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