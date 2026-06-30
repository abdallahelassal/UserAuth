package repository

import (
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type PersonalAccessToken struct{
	ID			uuid.UUID 	`gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash	string 		`gorm:"not null;uniqueIndex"`
	UserID		uuid.UUID	`gorm:"index;type:uuid;not null;"`

	User 		*User 		`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	TokenName	string		`gorm:"not null"`
	LastUsedAt	*time.Time 	`gorm:"type:timestamptz"`
	ExpiresAt	*time.Time	`gorm:"type:timestamptz"`
	CreatedAt  	time.Time	`gorm:"autoCreateTime;index"`
}

func (p *PersonalAccessToken) BeforeCreate(tx *gorm.DB) error{
	if p.ID == uuid.Nil{
		p.ID = uuid.New()
	}
	return nil
}


func (p *PersonalAccessToken) ToDomainPersonalToken()*domain.PersonalAccessToken{
	return &domain.PersonalAccessToken{
		ID: p.ID,
		TokenHash: p.TokenHash,
		UserID: p.UserID,
		TokenName: p.TokenName,
		LastUsedAt: p.LastUsedAt,
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
		LastUsedAt: p.LastUsedAt,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}
