package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Base struct {
	UUID		uuid.UUID		`json:"uuid"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type User struct {
	Base
	UserName     string 	`json:"user_name"`
	Email    	string 		`json:"email"`
	Password 	string 		`json:"password"`
	IsActive 	bool			`json:"is_active"`
}

func (b *Base) BeforeCreate(tx *gorm.DB)(err error){
	b.UUID = uuid.Must(uuid.NewV6())
	return
}
