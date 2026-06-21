package repository

import (
	"context"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"

	"gorm.io/gorm"
)


type  UserRepository struct{
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB)*UserRepository{
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User)error{
	model := FromDomain(user)
	
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
	
}

func (r *UserRepository) GetByEmail(ctx context.Context,email string)(*domain.User,error){
	var model User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	user := model.ToDomain()
	return user, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string)(*domain.User, error){
	var model User
	err := r.db.WithContext(ctx).Where("user_name = ?", name).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound{
			return nil, domain.ErrUserNotFound
		}
		return nil,err
	}
	user := model.ToDomain()
	return user,nil
}

func (r *UserRepository) FindByID(ctx context.Context,userID uuid.UUID)(*domain.User, error){
	var model User
	if err := r.db.WithContext(ctx).Where("ID = ?" , userID).Preload("Roles").First(&model).Error; err!= nil {
		return nil , err 
	}
	user := model.ToDomain()
	return  user, nil
}

func (r *UserRepository) fetchUser(ctx context.Context,query *gorm.DB ,limit int) ([]domain.User,error){
	var models []User
	err := query.WithContext(ctx).
			Limit(limit + 1).
			Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]domain.User,len(models))
	for i, model := range models {
		users[i] = *model.ToDomain()
	}
	return users, nil
}


func (r *UserRepository) Fetch(ctx context.Context, cursor string, limit int) ([]domain.User, string, error) {

	var models []User
	var cursorTime time.Time
	var err error

	// decode cursor
	if cursor != "" {
		cursorTime, err = DecodeCursor(cursor)
		if err != nil {
			return nil, "", domain.ErrBadParamInput
		}
	}

	// build query
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if cursor != "" {
		query = query.Where("created_at < ?", cursorTime)
	}

	// fetch + 1
	err = query.Limit(limit + 1).Find(&models).Error
	if err != nil {
		return nil, "", err
	}

	// next cursor logic
	var nextCursor string

	hasNext := len(models) > limit
	if hasNext {
		last := models[limit-1]
		nextCursor = EncodeCursor(last.CreatedAt)
		models = models[:limit]
	}

	// map to domain
	users := make([]domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, *m.ToDomain())
	}

	return users, nextCursor, nil
}

func (r *UserRepository) AssignRole(ctx context.Context,id uuid.UUID,roleID uuid.UUID)error{
	var user User
	if err := r.db.WithContext(ctx).First(&user,"id = ?", id).Error; err != nil {
		return err 
	}
	var role Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
	return err
	}
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Append(&role)
}


