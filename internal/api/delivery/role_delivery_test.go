package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/abdallahelassal/UserAuth/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type eqRoleMatcher struct {
	expected interface{}
}

func NeweqRoleMatcher(expected interface{}) gomock.Matcher {
	return &eqRoleMatcher{
		expected: expected,
	}
}

func (m *eqRoleMatcher) Matches(x interface{}) bool {

	switch expected := m.expected.(type) {
	case usecase.RoleCreateInput:
		got, ok := x.(usecase.RoleCreateInput)
		if !ok {
			return false
		}
		return got.Name == expected.Name && reflect.DeepEqual(got.PermissionIDs, expected.PermissionIDs)
	case usecase.RoleUpdateInput:
		got, ok := x.(usecase.RoleUpdateInput)
		if !ok {
			return false
		}
		return got.ID == expected.ID && got.Name == expected.Name && reflect.DeepEqual(got.PermissionIDs, expected.PermissionIDs)
	case usecase.RoleDeleteInput:
		got, ok := x.(usecase.RoleDeleteInput)
		if !ok {
			return false
		}
		return got.ID == expected.ID
	}
	return false

}

func (m *eqRoleMatcher) String() string {
	return "matchs role (create/update/delete input)"
}

func TestFindAll_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("positive_FindAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecaseRole := mocks.NewMockRoleUsecase(ctrl)
		handler := NewRoleDelivery(usecaseRole)
		w := httptest.NewRecorder()
		router := gin.New()
		router.GET("/roles", handler.FindAll)
		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		expectedRoles := []*usecase.RoleOutput{
			{ID: uuid.New(), Name: faker.Name()},
			{ID: uuid.New(), Name: faker.Name()},
		}
		usecaseRole.EXPECT().FindAll(gomock.Any()).Return(expectedRoles, nil)

		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("negative_FindAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecaseRole := mocks.NewMockRoleUsecase(ctrl)
		handler := NewRoleDelivery(usecaseRole)
		w := httptest.NewRecorder()
		router := gin.New()
		router.GET("/roles", handler.FindAll)
		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		usecaseRole.EXPECT().FindAll(gomock.Any()).Return(nil, domain.ErrInternalServer)

		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestFindByID_Handler(t *testing.T){
	gin.SetMode(gin.TestMode)

	t.Run("posetive_FindByID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecaseRole := mocks.NewMockRoleUsecase(ctrl)
		handler := NewRoleDelivery(usecaseRole)
		w := httptest.NewRecorder()
		roleID := uuid.New()
		router := gin.New()
		router.GET("/role/:id", handler.FindByID)
		req, _ := http.NewRequest(http.MethodGet,"/role/"+roleID.String(),nil)

		expectedRole := &usecase.RoleOutput{
			ID: roleID,
			Name: faker.Name(),
		}
		usecaseRole.EXPECT().FindByID(gomock.Any(), roleID).Return(expectedRole, nil)
		router.ServeHTTP(w,req)
		require.Equal(t, http.StatusOK,w.Code)
	})
	t.Run("negative_FindByID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseRole := mocks.NewMockRoleUsecase(ctrl)
		handler := NewRoleDelivery(usecaseRole)
		w := httptest.NewRecorder()
		roleID := uuid.New()

		router := gin.New()
		router.GET("/role/:id", handler.FindByID)
		req , _ := http.NewRequest(http.MethodGet,"/role/"+roleID.String(),nil)
		usecaseRole.EXPECT().FindByID(gomock.Any(),roleID).Return(nil,domain.ErrInternalServer)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusInternalServerError,w.Code)
	})
}

func TestCreate_handler(t *testing.T){
	gin.SetMode(gin.TestMode)
	t.Run("positive_CreateRole", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecaseRole := mocks.NewMockRoleUsecase(ctrl)
		handler := NewRoleDelivery(usecaseRole)
		router := gin.New()
		w := httptest.NewRecorder()
		permissionID := []uuid.UUID{uuid.New(), uuid.New()}
		expectedRole := usecase.RoleCreateInput{
			Name: faker.Name(),
			PermissionIDs: permissionID,
		}
		router.POST("/role/create", handler.Create)
		
		body := `{"name": "` + expectedRole.Name + `", "permission_ids": ["` + permissionID[0].String() + `", "` + permissionID[1].String() + `"]}`
		
		req , _ := http.NewRequest(http.MethodPost,"/role/create", bytes.NewBufferString(body))
		usecaseRole.EXPECT().Create(gomock.Any(),NeweqRoleMatcher(expectedRole)).Return(nil)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusOK,w.Code)

	})
}
