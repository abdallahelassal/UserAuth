package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Base struct {
	UUID		string		`json:"uuid"`
	CreatedAT	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type User struct {
	Base
	Name     string 	`json:"name"`
	Email    string 	`json:"email"`
	Password string 	`json:"password"`
}

func (b *Base) BeforeCreate(tx *gorm.DB)(err error){
	b.UUID = uuid.New().String()
	return
}
