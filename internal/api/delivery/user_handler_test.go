package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/abdallahelassal/UserAuth/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqUserMatcher struct{
	expectName	string
	expectEmail string
	plainPass	string
}
func NeweqUserMatcher(name,email,password string)gomock.Matcher{
	return &eqUserMatcher{
		expectName: name,
		expectEmail: email,
		plainPass: password,
	}
}

func (m *eqUserMatcher) Matches(x interface{})bool{
	c , ok := x.(usecase.CreateUserInput)
	if !ok {
		return false
	}
	return c.Email == m.expectEmail && c.UserName == m.expectName && c.Password == m.plainPass
}

func (m *eqUserMatcher) String()string{
	return "matcher user failed"
}

func TestSignup_handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserDelivary(usecase)

	t.Run("positive_signup",func(t *testing.T) {
		reqBody := `{
			"user_name": "abdallah",
			"email" : "test@test.com",
			"password" : "25556660"
		}`
		w := httptest.NewRecorder()
		c , _ := gin.CreateTestContext(w)

		req , _ := http.NewRequest(http.MethodPost,"/signup", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type","application/json")
		c.Request =req
		gomock.InOrder(
		usecase.EXPECT().GetByEmail(gomock.Any(),"test@test.com").Return(nil,errors.New("user not found")),
		usecase.EXPECT().Signup(gomock.Any(),NeweqUserMatcher("abdallah", "test@test.com", "25556660")).Return(nil),
		)
		handler.Signup(c)
		
		require.Equal(t, http.StatusCreated, w.Code)
	})
	
}
func TestLogin_handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecaseUser := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserDelivary(usecaseUser)
	
	t.Run("positive_signin", func(t *testing.T) {
		reqBody := `{
			"email" : "test@test.com",
			"password" : "123456"
		}`
		w := httptest.NewRecorder()
		c , _ := gin.CreateTestContext(w)

		req , _ := http.NewRequest(http.MethodPost,"/signin",bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type","application/json")
		c.Request = req
	
		usecaseUser.EXPECT().Login(gomock.Any(), gomock.AssignableToTypeOf(usecase.LoginUserInput{})).Return("token", nil)
		handler.Login(c)
		require.Equal(t,http.StatusOK,w.Code)
	})
	t.Run("negative_signin", func(t *testing.T) {
		reqBody := `{
			"email" : "test@test.com",
			"password" : "wrongpassword"
		}`
		w := httptest.NewRecorder()
		c , _ := gin.CreateTestContext(w)

		req , _ := http.NewRequest(http.MethodPost,"/signin",bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type","application/json")
		c.Request = req
	
		usecaseUser.EXPECT().Login(gomock.Any(), gomock.AssignableToTypeOf(usecase.LoginUserInput{})).Return("", domain.ErrInvalidCredentials)
		handler.Login(c)
		require.Equal(t,http.StatusUnauthorized,w.Code)
	})
}
func TestProfile_handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecaseUser := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserDelivary(usecaseUser)

	t.Run("positive_profile", func(t *testing.T) {
		userID := uuid.New()

		w := httptest.NewRecorder()
		c , _ := gin.CreateTestContext(w)
		req , _ := http.NewRequest(http.MethodGet,"/:id",bytes.NewBufferString(userID.String()))
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}

		expectedUser := &usecase.FindByIDOutput{
			UserName: "testuser",
			Email: "test@test.com",
			IsActive: true,
			Roles: []usecase.RoleOutput{
				{
					Name: "admin",
			
				},
			},	
		}
		usecaseUser.EXPECT().FindByID(gomock.Any(), userID).Return(expectedUser, nil)
		handler.Profile(c)
		require.Equal(t,http.StatusOK,w.Code)		

	})
	t.Run("negative_profile", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		userID := uuid.New()

		w := httptest.NewRecorder()
		c , _ := gin.CreateTestContext(w)
		req , _ := http.NewRequest(http.MethodGet,"/:id",bytes.NewBufferString(userID.String()))
		c.Request =req
		c.Params = gin.Params{gin.Param{Key: "id",Value: userID.String()}}
		
		usecaseUser.EXPECT().FindByID(gomock.Any(),userID).Return(nil,domain.ErrInternalServer)
		handler.Profile(c)
		require.Equal(t,http.StatusInternalServerError,w.Code)
	})
}

func TestAssignRolesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecaseUser := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserDelivary(usecaseUser)

	setupRequest := func(userID string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(
			http.MethodPost,
			"/users/"+userID+"/assign-roles",
			body,
		)

		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: userID}}

		return c, w
	}

	t.Run("positive_assign_roles", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		type assignRoleRequest struct {
			RoleID string `json:"role_id"`
		}
		
		reqBody := assignRoleRequest{
			RoleID: roleID.String(),
		}

		jsonBody , _ := json.Marshal(reqBody)


		c, w := setupRequest(userID.String(), bytes.NewReader(jsonBody))

		usecaseUser.EXPECT().
			AssignRole(gomock.Any(), userID, roleID).
			Return(nil)

		handler.AssignRoles(c)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("negative_internal_error", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		type assignRoleRequest struct{
			RolesID string `json:"role_id"`
		}
		reqBody := assignRoleRequest{
			RolesID: roleID.String(),
		}
		jsonBody , _ := json.Marshal(reqBody)


		c, w := setupRequest(userID.String(), bytes.NewReader(jsonBody))

		usecaseUser.EXPECT().
			AssignRole(gomock.Any(), userID, roleID).
			Return(domain.ErrInternalServer)

		handler.AssignRoles(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

}


func TestMeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecaseUser := mocks.NewMockUserUsecase(ctrl)
	handler := NewUserDelivary(usecaseUser)

	t.Run("positive_me", func(t *testing.T) {
	userID := uuid.New()
	w := httptest.NewRecorder()

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("id", userID.String())
	})

	router.GET("/user/me", handler.Me)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/user/me",
		nil,
	)

	expectedProfile := &usecase.FullProfile{
		User: usecase.UserOutput{
			Email: faker.Email(),
			UserName: faker.Name(),
		},
	}

	usecaseUser.EXPECT().
		GetFullProfile(gomock.Any(), userID).
		Return(expectedProfile, nil)

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
})
}