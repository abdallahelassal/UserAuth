package usecase

import (
	"time"

	"github.com/google/uuid"
)


type CreatePersonalAccessTokenInput struct{
	UserID		uuid.UUID
	TokenName	string
	ExpiresAt	*time.Time
}

type UpdatePersonalAccessTokenInput struct{
	ID			uuid.UUID
	UserID		uuid.UUID
	TokenName	string
	ExpiresAt	*time.Time
}

type DeletePersonalAccessTokenInput struct{
	ID			string
	UserID		string
}	

type GetByTokenOutput struct{
	ID			uuid.UUID
	UserID		uuid.UUID
	TokenName	string
	LastUseAt	*time.Time
	ExpiresAt	*time.Time
	CreatedAt	time.Time
}
