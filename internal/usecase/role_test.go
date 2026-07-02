package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	

	"github.com/golang/mock/gomock"
)





type eqRoleMatcher struct{
	expectedName string
}

func NeweqRoleMatcher(name string) gomock.Matcher {
	return &eqRoleMatcher{
		expectedName: name,
	}
}

func (r *eqRoleMatcher) Matches(x interface{}) bool{
	role , ok := x.(*domain.Role)
	
	if !ok {
		return false
	}
	return role.Name == r.expectedName
}

func (r *eqRoleMatcher) String() string {
	return "matches role fields"
}
func TestCreateRole_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mocks.NewMockRoleRepository(ctrl)
	usecase := NewRoleUseCase(roleRepo, time.Second)

	ctx := context.Background()

	t.Run("positive_CreateRole", func(t *testing.T) {
		permissionIDs := []uuid.UUID{uuid.New()}
		roleID := uuid.New()

		expectedRole := RoleCreateInput{
			Name: faker.Name(),
			PermissionIDs: permissionIDs,
		}
		
		roleRepo.EXPECT().Create(gomock.Any(),NeweqRoleMatcher(expectedRole.Name)).
		DoAndReturn(func (ctx context.Context, role *domain.Role) error {
			role.ID = roleID
			return nil
		})
		roleRepo.EXPECT().AssignPermission(gomock.Any(),roleID,permissionIDs).Return(nil)
		err := usecase.Create(ctx,expectedRole)
		require.NoError(t,err)
	})
	t.Run("negative_CreateRole", func(t *testing.T) {
		permissioIDs 	:= []uuid.UUID{uuid.New()}
		

		expextedRole := RoleCreateInput{
			Name: faker.Name(),
			PermissionIDs: permissioIDs,
		}
		roleRepo.EXPECT().Create(gomock.Any(),NeweqRoleMatcher(expextedRole.Name)).Return(errors.New("failed create role"))
		err := usecase.Create(ctx,expextedRole)
		require.Error(t,err)
	})
}

func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mocks.NewMockRoleRepository(ctrl)
	usecase := NewRoleUseCase(roleRepo, time.Second)
	ctx := context.Background()

	t.Run("positive_FindRoleByID", func(t *testing.T) {
		roleID := uuid.New()
		expectedRole := &domain.Role{
			Base: domain.Base{ID: roleID},
			Name: faker.Name(),
		}
		roleRepo.EXPECT().FindByID(gomock.Any(), roleID).Return(expectedRole, nil)
		role , err := usecase.FindByID(ctx,roleID)
		require.NoError(t,err)
		require.NotNil(t,role)
	})

}