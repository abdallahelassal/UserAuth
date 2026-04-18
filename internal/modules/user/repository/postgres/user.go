package postgres

import (
	"context"
	"time"

	"github.com/abdallahelassal/UserAuth.git/domain"
	"github.com/abdallahelassal/UserAuth.git/internal/modules/user/repository"
	"gorm.io/gorm"
)


type  postgresUserRepository struct{
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB)*postgresUserRepository{
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User)error{
	return r.db.WithContext(ctx).Create(user).Error 
}

func (r *postgresUserRepository) fetch(ctx context.Context,query *gorm.DB ,limit int) ([]domain.User,error){
	var users []domain.User
	err := r.db.WithContext(ctx).
			Limit(limit + 1).
			Find(&users).Error
	return users , err			
}

func (r *postgresUserRepository) Fetch(ctx context.Context, cursor string, limit int)([]domain.User,string,error){
	var users 		[]domain.User
	var cursorTime 	time.Time
	var err 		error

	// decode cursor 
	if cursor != ""{
		cursorTime , err = repository.DecodeCursor(cursor)
		if err != nil {
			return nil , "", domain.ErrBadParamInput
		}
	}

	// build query
	query := r.db.Order("created_at DESC")

	if cursor != ""{
		query = r.db.Where("created_at < ?", cursorTime)
	}

	// call fetch 

	users , err = r.fetch(ctx, query, limit)
	if err != nil {
		return nil , "",err
	}

	// next cursor
	
	var nextCursor = ""
	
	if len(users) > limit {
		last := users[limit-1]
		nextCursor = repository.EncodeCursor(last.CreatedAT)
		users = users[:limit]
	}
	return users , nextCursor , nil
}





