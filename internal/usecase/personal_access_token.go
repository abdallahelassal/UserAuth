package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	
	"github.com/google/uuid"
)

type PersonalAccessTokenUsecase interface{
	Create(ctx context.Context, req CreatePersonalAccessTokenInput)error
	FindByToken(ctx context.Context,tokenHash string)(GetByTokenOutput,error)
	Delete(ctx context.Context,tokenHash string)error
	DeleteByUserID(ctx context.Context,userID uuid.UUID)error
	UpdateLastUsed(ctx context.Context,tokenID uuid.UUID)error
}

type personalAccessTokenUsecase struct{
	PersonalAccessTokenRepo domain.PersonalAccessTokenRepository
	ContextTimeOut time.Duration
}

func NewPersonalAccessTokenUsecase(personalAccessToken domain.PersonalAccessTokenRepository, timeOut time.Duration )*personalAccessTokenUsecase{
	return &personalAccessTokenUsecase{
		PersonalAccessTokenRepo: personalAccessToken,
		ContextTimeOut: timeOut,
	}
}

func (p *personalAccessTokenUsecase) Create(ctx context.Context,req CreatePersonalAccessTokenInput)error{
	ctx , cancel := context.WithTimeout(ctx,p.ContextTimeOut)
	defer cancel()

	if req.UserID == uuid.Nil && req.TokenName == "" {
		return errors.New("a user not found")
	}

	pToken := &domain.PersonalAccessToken{
		UserID: req.UserID,
		ExpiresAt: req.ExpiresAt,
		TokenName: req.TokenName,
	}

	err := p.PersonalAccessTokenRepo.Create(ctx, pToken)
	if err != nil {
		return err
	}
	return nil
}

func (p *personalAccessTokenUsecase) FindByToken(ctx context.Context, tokenHash string)(GetByTokenOutput,error){
	ctx , cancel := context.WithTimeout(ctx , p.ContextTimeOut)
	defer cancel()

	if tokenHash == ""{
		return GetByTokenOutput{}, errors.New("validate token")
	}

	token , err := p.PersonalAccessTokenRepo.FindByToken(ctx,tokenHash)
	if err != nil {
		return GetByTokenOutput{} , err
	}

	return GetByTokenOutput{
		ID: token.ID,
		UserID: token.UserID,
		TokenName: token.TokenName,
		LastUseAt: token.LastUseAt,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	} , nil

}
func (p *personalAccessTokenUsecase) UpdateLastUsed(ctx context.Context,tokenHash string )error{
	ctx , cancel := context.WithTimeout(ctx,p.ContextTimeOut)
	defer cancel()


	if tokenHash == "" {
		return errors.New("validate token")
	}
	return p.PersonalAccessTokenRepo.UpdateLastUseAt(ctx,tokenHash)

}

func (p *personalAccessTokenUsecase) Delete(ctx context.Context,tokenHash string)error{
	ctx , cancel := context.WithTimeout(ctx,p.ContextTimeOut)
	defer cancel()
	if err := p.PersonalAccessTokenRepo.Delete(ctx,tokenHash); err != nil {
		return err
	}
	return nil 
}

func (p *personalAccessTokenUsecase) DeleteByUserID(ctx context.Context,userID uuid.UUID)error{
	ctx , cancel := context.WithTimeout(ctx,p.ContextTimeOut)
	defer cancel()
	if err := p.PersonalAccessTokenRepo.DeleteByUserID(ctx,userID); err != nil {
		return err
	}
	return nil
}
