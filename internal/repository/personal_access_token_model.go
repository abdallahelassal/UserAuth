package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type PersonalAccessToken struct{
	ID			uuid.UUID 	`gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash	string 		`gorm:"not null"`
	UserID		uuid.UUID	`gorm:"index;type:uuid;not null"`
	TokenName	string		`gorm:"not null"`
	LastUseAt	*time.Time 	`gorm:"type:timestamptz"`
	ExpiresAt	*time.Time	`gorm:"type:timestamptz"`
	CreatedAt  	time.Time	`gorm:"autoCreateTime;index"`
}

func (p *PersonalAccessToken) BeforeCreate(tx *gorm.DB) error{
	if p.ID == uuid.Nil{
		p.ID = uuid.New()
	}
	return nil
}


func (p *PersonalAccessToken) ToDomain()*domain.PersonalAccessToken{
	return &domain.PersonalAccessToken{
		ID: p.ID,
		TokenHash: p.TokenHash,
		UserID: p.UserID,
		TokenName: p.TokenName,
		LastUseAt: p.LastUseAt,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}

func FromDomainPersonalAccessToken(p *domain.PersonalAccessToken)*PersonalAccessToken{
	return &PersonalAccessToken{
		ID: p.ID,
		TokenHash: p.TokenHash,
		UserID: p.UserID,
		TokenName: p.TokenName,
		LastUseAt: p.LastUseAt,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}
