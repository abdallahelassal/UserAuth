package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Base struct {
	UUID		string		`gorm:"type:uuid;primary_key"`
	CreatedAT	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type User struct {
	
	Name     string 	`json:"name" gorm:"type:varchar(255);not null"`
	Email    string 	`json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string 	`json:"password" gorm:"type:varchar(255);not null"`
}

func (b *Base) BeforeCreate(tx *gorm.DB)(err error){
	b.UUID = uuid.New().String()
	return
}
