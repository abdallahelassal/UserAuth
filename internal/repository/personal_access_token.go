package repository

import (
	"context"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
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

func (r *PersonalAccessTokenRepository) FindByToken(ctx context.Context,hash string)(*domain.PersonalAccessToken,error){
	var dbModel PersonalAccessToken

	if err := r.db.WithContext(ctx).Where("token_hash = ?",hash).Take(&dbModel).Error; err != nil {
		return nil,err
	}
	return dbModel.ToDomainPersonalToken(),nil
}

 func (r *PersonalAccessTokenRepository) UpdateLastUseAt(ctx context.Context,hash string)error{
	result := r.db.WithContext(ctx).Model(&PersonalAccessToken{}).Where("token_hash = ?", hash).Update("last_use_at = ?",time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return  gorm.ErrRecordNotFound
	}
	return nil
}
func (r *PersonalAccessTokenRepository) Delete(ctx context.Context,hash string)error{
	result := r.db.WithContext(ctx).Where("token_hash = ?", hash).Delete(&PersonalAccessToken{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return  gorm.ErrRecordNotFound
	}
	return nil

}

func (r *PersonalAccessTokenRepository) DeleteByUserID(ctx context.Context,userID uuid.UUID)error{
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PersonalAccessToken{})

		if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return  gorm.ErrRecordNotFound
	}
	return result.Error

}
 