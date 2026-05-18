package repository

import (
	"context"

	"github.com/abdallahelassal/UserAuth/domain"
	"gorm.io/gorm"
)



type PersonalAccessTokenRepository struct{
	db *gorm.DB
}

func NewPersonalAccessTokenRepository(db *gorm.DB)*PersonalAccessTokenRepository{
	return &PersonalAccessTokenRepository{
		db: db,
	}
}

func (r *PersonalAccessTokenRepository) Create(ctx context.Context, model *domain.PersonalAccessToken) error {
	dbmodel := FromDomainPersonalAccessToken(model)
	if err := r.db.WithContext(ctx).Create(dbmodel).Error; err != nil {
		return err
	}
	return nil
}
