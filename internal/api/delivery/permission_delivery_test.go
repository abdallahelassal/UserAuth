package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"
	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/abdallahelassal/UserAuth/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)


type eqPermissionMatcher struct{
	expected	interface{}
}

func NeweqPermissionMatcher(expected interface{}) gomock.Matcher{
	return &eqPermissionMatcher{
		expected: expected,
	}
}

func (p *eqPermissionMatcher) Matches(x interface{}) bool{

	switch expected := p.expected.(type) {
	case usecase.PermissionInput:
		got , ok := x.(usecase.PermissionInput)
		if !ok {
			return false
		}
		return  got.Name == expected.Name
	case usecase.UpdatePermissionInput:
		got , ok := x.(usecase.UpdatePermissionInput)
		if !ok {
			return false
		}
		return got.Name == expected.Name 	
	}
	return false
}

func (p *eqPermissionMatcher) String()string{
	return "matches permission input"
}

func TestCreatePerm_handler(t *testing.T){
	gin.SetMode(gin.TestMode)
	
	t.Run("positive_create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()

		router.POST("/create", handler.Create)
		expectedPerm := usecase.PermissionInput{
			Name: faker.Name(),
		}
		body := `{"name":"` + expectedPerm.Name + `"}`

	

		usecasePerm.EXPECT().Create(gomock.Any(),gomock.Any()).DoAndReturn(func(ctx context.Context, input *usecase.PermissionInput) error {
		require.Equal(t, expectedPerm.Name, input.Name)
		return nil
	})
		req , _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(body))
		router.ServeHTTP(w, req)

		require.Equal(t,http.StatusOK, w.Code)
	})
	t.Run("negative_create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()

		router.POST("/create", handler.Create)

		body := `{"name": "ahmed"}`
		req , _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(body))
		usecasePerm.EXPECT().Create(gomock.Any(),gomock.Any()).Return(domain.ErrInternalServer)
		router.ServeHTTP(w,req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
func TestFindAll_handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("positive_FindAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()

		router.GET("/permission", handler.FindAll)
		req , _ := http.NewRequest(http.MethodGet, "/permission", nil)
		
		expected := []usecase.PermissionOutput{
			{ID: uuid.New(), Name: faker.Name()},
			{ID: uuid.New(), Name: faker.Name()},
		}
		usecasePerm.EXPECT().FindAllPermissions(gomock.Any()).Return(expected,nil)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusOK, w.Code)
	})
	t.Run("negative_FindAll" , func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()

		router.GET("/permission", handler.FindAll)
		req , _ := http.NewRequest(http.MethodGet, "/permission", nil)
		usecasePerm.EXPECT().FindAllPermissions(gomock.Any()).Return(nil , domain.ErrInternalServer)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusInternalServerError, w.Code)
	})
}

func TestFindPermissionByUserID_handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("positive_FindPermissionByUserID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()
		userID := uuid.New()

		router.GET("/permissions/:id", handler.FindPermissionByUserID)
		req , _ := http.NewRequest(http.MethodGet, "/permissions/"+ userID.String(),nil )
		expected:= []usecase.GetUserPermissions{
			{Name: faker.Name(), UserID: userID},
			{Name: faker.Name(), UserID: userID},
		}
		usecasePerm.EXPECT().GetPermissionsByUserID(gomock.Any(), userID).Return(expected,nil)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), expected[0].Name)
	})
	t.Run("negative_FindPermissionByUserID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()
		userID := uuid.New()

		router.GET("/permissions/:id", handler.FindPermissionByUserID)
		req, _ := http.NewRequest(http.MethodGet, "/permissions/"+ userID.String(), nil)
		usecasePerm.EXPECT().GetPermissionsByUserID(gomock.Any(), userID).Return(nil,domain.ErrInternalServer)
		router.ServeHTTP(w,req)
		require.Equal(t,http.StatusInternalServerError, w.Code)
	})
}

func TestFindPermissionByRoleID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("positive_FindPermissionByRoleID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()
		

		router.GET("/permissions", handler.FindPermissionByRoleID)
		roleID1 := uuid.New()
		roleID2 := uuid.New()
		ids := []uuid.UUID{roleID1, roleID2}
		expectedRoles := []usecase.GetPermissionsByRoleIDs{
			{Name: faker.Name(), RoleID: roleID1}, 
			{Name: faker.Name(), RoleID: roleID2}, 
		}
	
		usecasePerm.EXPECT().GetPermissionByRoleIDs(gomock.Any(),ids).Return(expectedRoles,nil)
		req, _ := http.NewRequest(http.MethodGet, "/permissions?id="+roleID1.String()+","+roleID2.String(), nil)
		router.ServeHTTP(w,req)

		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("negative_FindPermissionByRoleIDs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecasePerm := mocks.NewMockPermissionUsecase(ctrl)
		handler := NewPermissionDelivery(usecasePerm)
		router := gin.New()
		w := httptest.NewRecorder()
		
		router.GET("/permissions", handler.FindPermissionByRoleID)

		roleID1 := uuid.New()
		roleID2 := uuid.New()
		ids := []uuid.UUID{roleID1, roleID2}
		
		usecasePerm.EXPECT().GetPermissionByRoleIDs(gomock.Any(), ids).Return(nil, domain.ErrInternalServer)
		req , _ := http.NewRequest(http.MethodGet, "/permissions?id="+roleID1.String()+","+roleID2.String(),nil)
		router.ServeHTTP(w,req)

		require.Equal(t,http.StatusInternalServerError, w.Code)
	})
}
