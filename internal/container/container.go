package container

import (
	"time"

	"github.com/abdallahelassal/UserAuth/internal/api/delivary"
	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/modules/user/repository"
	"github.com/abdallahelassal/UserAuth/internal/modules/user/usecase"
	"go.uber.org/zap"
	"gorm.io/gorm"
)


type Container struct{
	UserDelivary *delivary.UserDelivary
	Cfg bootstrap.Config
	Logger *zap.Logger
}

func NewContainer(db *gorm.DB, logger *zap.Logger, cfg bootstrap.Config) *Container {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo, 5*time.Second)
	userDelivary := delivary.NewUserDelivary(userUsecase,cfg)
	return &Container{
		UserDelivary: userDelivary,
		Cfg: cfg,
		Logger: logger,
	}	
}