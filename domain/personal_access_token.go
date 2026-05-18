package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)



type PersonalAccessToken struct{
	ID			uuid.UUID 	`json:"-"`
	TokenHash	string 		`json:"-"`
	UserID		uuid.UUID	`json:"user_id"`
	TokenName	string		`json:"token_name"`
	LastUseAt	*time.Time 	`json:"last_use_at"`
	ExpiresAt	*time.Time	`json:"expires_at"`
	CreatedAt  	time.Time	`json:"created_at"`
}

type PersonalAccessTokenRepository interface{
	Create(ctx context.Context, hash *PersonalAccessToken)error
	FindByToken(ctx context.Context,hash string)(*PersonalAccessToken,error)
	Delete(ctx context.Context,hash string)error
	DeleteByUserID(ctx context.Context, userID uuid.UUID)error
	UpdateLastUseAt(ctx context.Context, token string)error
}

