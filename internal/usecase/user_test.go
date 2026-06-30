package usecase

import (
	"context"
	

	"testing"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository/mocks"
	"github.com/abdallahelassal/UserAuth/pkg/bcrypt"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type eqUserMatcher struct {
	expectName string
	expectEmail string
	plainPass	string
}

func NeweqUserMatcher(name,email,password string)gomock.Matcher{
	return eqUserMatcher{
		expectName: name,
		expectEmail: email,
		plainPass: password,
	}
}


func (m eqUserMatcher) Matches(x interface{})bool{
	u , ok := x.(*domain.User)
	if !ok {
		return false
	}
	return u.Email == m.expectEmail &&
	 u.UserName == m.expectName
}




func (m eqUserMatcher) String() string {
	return "matches user fields"
}

func TestCreateUser_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo 		:= mocks.NewMockUserRepository(ctrl)
	roleRepo 		:= mocks.NewMockRoleRepository(ctrl)
	permissionRepo 	:= mocks.NewMockPermissionRepository(ctrl)

	usecase := NewUserUseCase(userRepo,roleRepo,permissionRepo,time.Second,"secret", time.Hour)

	passStr := "password"
	
	
	
	user := CreateUserInput{
		
		UserName: "abdallah",
		Email: "abdallah@test.com",
		Password: passStr,
	}

	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "user",
	}
	userRepo.EXPECT().Create(gomock.Any(),gomock.Any()).DoAndReturn(func (ctx context.Context,u *domain.User) error {
		if u.Password == passStr {
			t.Fatalf("password should be hashed")
		}
		u.Base.ID = uuid.New()
		return nil
	})
	roleRepo.EXPECT().FindByName(gomock.Any(), "user").Return(role,nil)
	userRepo.EXPECT().AssignRole(gomock.Any(),gomock.Any(),gomock.Any()).Return(nil)
	if err := usecase.Signup(context.Background(),user); err != nil{
		t.Errorf("expected err %v",err)
	}
}

func TestFindByID_usecase(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	userRepo 		:= mocks.NewMockUserRepository(ctrl)
	roleRepo 		:= mocks.NewMockRoleRepository(ctrl)
	permissionRepo 	:= mocks.NewMockPermissionRepository(ctrl)
	usecase := NewUserUseCase(userRepo,roleRepo,permissionRepo,time.Second,"secret", time.Hour)
	userID := uuid.New()

	t.Run("positive_findByID",func(t *testing.T) {
		expectedUser := &domain.User{
			Base: domain.Base{ID: userID},
			UserName:faker.Name() ,
			Email: faker.Email(),
			Password: faker.Password(),
		}
		userRepo.EXPECT().FindByID(gomock.Any(), userID).Return(expectedUser,nil)
		user , err :=usecase.FindByID(ctx,userID)
		
		require.NoError(t, err)
		require.NotNil(t,user)
		require.Equal(t,expectedUser.Email, user.Email)
	})
	t.Run("negative_findByID_not_found",func(t *testing.T) {
		userRepo.EXPECT().FindByID(gomock.Any(),userID).Return(nil,gorm.ErrRecordNotFound)
		detail , err := usecase.FindByID(ctx,userID)
		require.EqualError(t, err, "user not found")
		require.Equal(t,detail, FindByIDOutput{})
	})
}
func TestGetByEmail_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	roleRepo := mocks.NewMockRoleRepository(ctrl)
	permissionRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewUserUseCase(userRepo,roleRepo,permissionRepo,time.Second,"secret", time.Hour)
	userID := uuid.New()
	ctx := context.Background()

	t.Run("positive_GetByEmail",func(t *testing.T) {
		expectedUser := &domain.User{
			Base: domain.Base{ID: userID},
			UserName: faker.Name(),
			Email: faker.Email(),
			Password: faker.Password(),
		}
		userRepo.EXPECT().GetByEmail(gomock.Any(),expectedUser.Email).Return(expectedUser,nil)
		user , err := usecase.GetByEmail(ctx,expectedUser.Email)
		require.NoError(t,err)
		require.NotNil(t,user)
		require.Equal(t,expectedUser.Email,user.Email)
	})
	t.Run("negative_GetByEmail",func(t *testing.T) {
		expecetEmail := faker.Email()
		userRepo.EXPECT().GetByEmail(gomock.Any(),expecetEmail).Return(nil, gorm.ErrRecordNotFound)
		user , err := usecase.GetByEmail(ctx,expecetEmail)
		require.Error(t,err)
		require.ErrorIs(t,err,gorm.ErrRecordNotFound)
		require.Empty(t,user)
	})


}

func TestLoginUser_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	roleRepo := mocks.NewMockRoleRepository(ctrl)
	permissionRepo := mocks.NewMockPermissionRepository(ctrl)
	usecase := NewUserUseCase(userRepo,roleRepo,permissionRepo,time.Second,"secret", time.Hour)
	ctx := context.Background()
	t.Run("positive_Login", func(t *testing.T) {
		plainPass := faker.Password()
		hash , err := bcrypt.HashPassword(plainPass)
		require.NoError(t,err)

		expectuser := &domain.User{
			Base: domain.Base{ID: uuid.New()},
			UserName: faker.Name(),
			Email: faker.Email(),
			Password: hash,
		}
		userRepo.EXPECT().GetByEmail(gomock.Any(),expectuser.Email).Return(expectuser,nil)
		token , err:=  usecase.Login(ctx,LoginUserInput{
			Email: expectuser.Email,
			Password: plainPass,
		})
		require.NoError(t,err)
		require.NotEmpty(t,token)
	})

	t.Run("negative_Login",func(t *testing.T) {
		plainPass := faker.Password()
		hash , err := bcrypt.HashPassword(plainPass)
		require.NoError(t,err)

		expectuser := &domain.User{
			Base: domain.Base{ID: uuid.New()},
			UserName: faker.Name(),
			Email: faker.Email(),
			Password: hash,
		}
		userRepo.EXPECT().GetByEmail(gomock.Any(),expectuser.Email).Return(expectuser,nil)
		token , err:=  usecase.Login(ctx,LoginUserInput{
			Email: expectuser.Email,
			Password: "wrongpassword",
		})
		require.Error(t,err)
		require.Empty(t,token)
	})

}
