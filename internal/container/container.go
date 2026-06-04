package container

import (
	"time"

	"github.com/abdallahelassal/UserAuth/internal/api/delivery"
	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/repository"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"go.uber.org/zap"
	
	"gorm.io/gorm"
)


type Container struct{
	UserDelivary *delivery.UserDelivary
	Cfg bootstrap.Config
	Logger *zap.Logger
}

func NewContainer(db *gorm.DB, logger *zap.Logger, cfg bootstrap.Config) *Container {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, 5*time.Second,cfg.JWTConfig.AccessTokenSecret, time.Duration(cfg.JWTConfig.AccessExpiration)*time.Hour)
	userDelivery := delivery.NewUserDelivary(userUsecase)
	return &Container{
		UserDelivary: userDelivery,
		Cfg: cfg,
		Logger: logger,
	}	
}