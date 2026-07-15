package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)



type eqPermissionMatcher struct{
	expectedName string
}

func NeweqPermissionMatcher(name string)gomock.Matcher{
	return &eqPermissionMatcher{
		expectedName: name,
	}
}

func (p *eqPermissionMatcher) Matches(x interface{}) bool{
	permission , ok := x.(*domain.Permission)
	if !ok {
		return false
	}
	return permission.Name == p.expectedName
}

func (p *eqPermissionMatcher) String()string{
	return "match permission failed"
}

func TestCreatPermission_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	permiRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewPermissionUsecase(permiRepo, time.Second)
	ctx := context.Background()
	
	t.Run("Positive_CreatePermission", func(t *testing.T) {
	
		expectedPermission := &PermissionInput{
			Name: faker.Name(),
		}
		permiRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&domain.Permission{})).Return(nil)
		err := usecase.Create(ctx,expectedPermission)
		require.NoError(t,err)
		require.NotEmpty(t, expectedPermission.Name)
	})

	t.Run("negative_CreatePermission", func(t *testing.T) {

		expectedPermission := PermissionInput{
			Name: faker.Name(),
		}
		permiRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("failed create premission"))
		err := usecase.Create(ctx,&expectedPermission)
		require.Error(t,err)
	})
}

func TestFindAllPermissions_usecase(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permiRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewPermissionUsecase(permiRepo, time.Second)
	ctx := context.Background()

	t.Run("positive_FindAllPermission", func(t *testing.T) {
		expectedPermissions := []domain.Permission{
			{Base: domain.Base{ID: uuid.New()},Name: faker.Name()},
			{Base: domain.Base{ID: uuid.New()},Name: faker.Name()},
			{Base: domain.Base{ID: uuid.New()},Name: faker.Name()},
		}
		permiRepo.EXPECT().FindAllPermissions(gomock.Any()).Return(expectedPermissions,nil)
		permissions, err := usecase.FindAllPermissions(ctx)
		require.NoError(t,err)
		require.Len(t,permissions, len(expectedPermissions))
	})
	t.Run("negative_FindAllPermission", func(t *testing.T) {
		permiRepo.EXPECT().FindAllPermissions(gomock.Any()).Return(nil,errors.New("failed fetch permissions"))
		permissions , err := usecase.FindAllPermissions(ctx)
		require.Error(t, err)
		require.Empty(t,permissions)
	})
}

func TestGetPermissionsByUserID_usecase(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permiRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewPermissionUsecase(permiRepo, time.Second)
	ctx := context.Background()

	t.Run("positive_GetPermissionsByUserID", func(t *testing.T) {
		userID := uuid.New()
		expectedPermissions := []domain.Permission{
			{Base: domain.Base{ID: uuid.New()}, Name: faker.Name()},
			{Base: domain.Base{ID: uuid.New()}, Name:faker.Name()},
		}

		permiRepo.EXPECT().GetPermissionsByUserID(gomock.Any(),userID).Return(expectedPermissions, nil)
		permissions , err := usecase.GetPermissionsByUserID(ctx,userID)
		require.NoError(t,err)
		require.Len(t,permissions, len(expectedPermissions))
	})
	t.Run("nigative_GetPermissionsByUserID", func(t *testing.T) {
		userID := uuid.New()

		permiRepo.EXPECT().GetPermissionsByUserID(gomock.Any(), userID).Return(nil, errors.New("permissions not found"))
		permissions , err := usecase.GetPermissionsByUserID(ctx,userID)
		require.Error(t,err)
		require.Nil(t, permissions)
	})
}

func TestGetPermissionsByRoleIDs_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permiRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewPermissionUsecase(permiRepo,time.Second)
	ctx := context.Background()

	t.Run("positive_GetPermissionByRoleIDs", func(t *testing.T) {
		roleIDs := []uuid.UUID{uuid.New()}

		expectedPermissions := []domain.Permission{
			{Base: domain.Base{ID: uuid.New()}, Name: faker.Name()},
			{Base: domain.Base{ID: uuid.New()}, Name:faker.Name()},
		}
		permiRepo.EXPECT().GetPermissionByRoleIDs(gomock.Any(), roleIDs).Times(1).Return(expectedPermissions , nil)
		permissions , err := usecase.GetPermissionByRoleIDs(ctx,roleIDs)
		require.NoError(t,err)
		require.NotNil(t,permissions)
	})
	t.Run("negative_GetPermissionByRoleIDs", func(t *testing.T) {
		rolesIDs := []uuid.UUID{uuid.New()}

		permiRepo.EXPECT().GetPermissionByRoleIDs(gomock.Any(),rolesIDs).Times(1).Return(nil,errors.New("permission not found"))
		permissions , err := usecase.GetPermissionByRoleIDs(ctx,rolesIDs)
		require.Error(t,err)
		require.Nil(t,permissions)
	})
}