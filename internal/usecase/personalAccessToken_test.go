package usecase

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/repository/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)



type eqPersonalAccessTokenMatcher struct{
	expectedTokenName string
}

func NewEqPersonalAccessTokenMatcher(tokenName string)gomock.Matcher{
	return &eqPersonalAccessTokenMatcher{
		expectedTokenName: tokenName,
	}
}

func (p *eqPersonalAccessTokenMatcher) Matches(x interface{}) bool{
	personalAccessToken , ok := x.(*domain.PersonalAccessToken)
	if !ok {
		return false
	}
	return personalAccessToken.TokenName == p.expectedTokenName
}

func (p *eqPersonalAccessTokenMatcher) String()string{
	return "match personal access token failed"
}	

func TestCreatePersonalAccessToken_usecase(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	accessTokenRep := mocks.NewMockPersonalAccessTokenRepository(ctrl)
	usecase := NewPersonalAccessTokenUsecase(accessTokenRep, time.Second)
	ctx := context.Background()

	t.Run("positive_CreatePersonalAccessToken",func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		expectedAccess := CreatePersonalAccessTokenInput{
			
			UserID: uuid.New(),
			ExpiresAt: &expiresAt,
			TokenName: faker.Name(),
		}
		accessTokenRep.EXPECT().Create(gomock.Any(),gomock.AssignableToTypeOf(&domain.PersonalAccessToken{})).Times(1).Return(nil)
		err := usecase.Create(ctx,expectedAccess)
		require.NoError(t,err)
	})
	t.Run("negative_CreatePersonalAccessToken",func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		expectedAccess := CreatePersonalAccessTokenInput{
			
			UserID: uuid.New(),
			ExpiresAt: &expiresAt,
			TokenName: faker.Name(),
		}
		accessTokenRep.EXPECT().Create(gomock.Any(),gomock.AssignableToTypeOf(&domain.PersonalAccessToken{})).Times(1).Return(exec.ErrNotFound)
		err := usecase.Create(ctx,expectedAccess)
		require.Error(t,err)
	})
}

func TestFindByToken_usecase(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	accessTokenRep := mocks.NewMockPersonalAccessTokenRepository(ctrl)
	usecase := NewPersonalAccessTokenUsecase(accessTokenRep, time.Second)
	ctx := context.Background()

	t.Run("positive_FindByToken",func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		expectedAccess := GetByTokenOutput{
			ID: uuid.New(),
			UserID: uuid.New(),
			TokenName: faker.Name(),
			LastUseAt: &expiresAt,
			ExpiresAt: &expiresAt,
			CreatedAt: time.Now(),
		}
		accessTokenRep.EXPECT().FindByToken(gomock.Any(),gomock.Any()).Times(1).Return(&domain.PersonalAccessToken{
			ID: expectedAccess.ID,
			UserID: expectedAccess.UserID,
			TokenName: expectedAccess.TokenName,
			LastUsedAt: expectedAccess.LastUseAt,
			ExpiresAt: expectedAccess.ExpiresAt,
			CreatedAt: expectedAccess.CreatedAt,
		},nil)
		result,err := usecase.FindByToken(ctx,expectedAccess.TokenName)
		require.NoError(t,err)
		require.Equal(t, expectedAccess, *result)
	})
	t.Run("negative_FindByToken",func(t *testing.T) {
		expiresAt := time.Now().Add(time.Hour)
		expectedAccess := GetByTokenOutput{
			ID: uuid.New(),
			UserID: uuid.New(),
			TokenName: faker.Name(),
			LastUseAt: &expiresAt,
			ExpiresAt: &expiresAt,
			CreatedAt: time.Now(),
		}
		accessTokenRep.EXPECT().FindByToken(gomock.Any(),gomock.Any()).Times(1).Return(nil,errors.New("token not found"))
		result,err := usecase.FindByToken(ctx,expectedAccess.TokenName)
		require.Error(t,err)
		require.Nil(t,result)
	})
}	

func TestUpdateLastUsed_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessTokenRep := mocks.NewMockPersonalAccessTokenRepository(ctrl)
	usecase := NewPersonalAccessTokenUsecase(accessTokenRep , time.Second)
	ctx := context.Background()

	t.Run("positive_updateLastUsed", func(t *testing.T) {
		tokenHash := "kdkdkdk"
	
		accessTokenRep.EXPECT().
		UpdateLastUseAt(gomock.Any(), tokenHash).
		Times(1).
		Return(nil)

		err := usecase.UpdateLastUsed(ctx,tokenHash)
		require.NoError(t,err)
		

	})
	t.Run("negative_updateLastUsed", func(t *testing.T) {
		tokenHash := "kfkfkf"

		accessTokenRep.EXPECT().UpdateLastUseAt(gomock.Any(),tokenHash).Times(1).Return(errors.New("failed update time"))

		err := usecase.UpdateLastUsed(ctx, tokenHash)
		require.Error(t,err)
	})
}

func TestDelete_usecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessTokenRep := mocks.NewMockPersonalAccessTokenRepository(ctrl)
	usecase := NewPersonalAccessTokenUsecase(accessTokenRep, time.Second)
	ctx := context.Background()

	t.Run("positive_deleteAccessToken", func(t *testing.T) {
		tokenHash := "kdkdk"
		accessTokenRep.EXPECT().Delete(gomock.Any(), tokenHash).Times(1).Return(nil)

		err := usecase.Delete(ctx,tokenHash)
		require.NoError(t,err)
	})
	
	t.Run("negative_deleteAccessToken", func(t *testing.T) {
		tokenHash := "kkfkfkf"
		accessTokenRep.EXPECT().Delete(gomock.Any(), tokenHash).Times(1).Return(errors.New("failed delete token"))

		err := usecase.Delete(ctx , tokenHash)
		require.Error(t,err)
	})
}

func TestDeleteByUserIDAccessToken_usecse(t *testing.T) {
		ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accessTokenRep := mocks.NewMockPersonalAccessTokenRepository(ctrl)
	usecase := NewPersonalAccessTokenUsecase(accessTokenRep, time.Second)
	ctx := context.Background()

	t.Run("positive_deleteByUserIDAccessToken", func(t *testing.T) {
		userID := uuid.New()
		accessTokenRep.EXPECT().DeleteByUserID(gomock.Any(), userID).Times(1).Return(nil)

		err := usecase.DeleteByUserID(ctx,userID)
		require.NoError(t,err)
	})
	
	t.Run("negative_deleteByUserIDAccessToken", func(t *testing.T) {
		userID := uuid.New()
		accessTokenRep.EXPECT().DeleteByUserID(gomock.Any(), userID).Times(1).Return(errors.New("failed delete token"))

		err := usecase.DeleteByUserID(ctx , userID)
		require.Error(t,err)
	})	
}