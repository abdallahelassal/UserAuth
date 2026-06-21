package container

import (
	"time"

	"github.com/abdallahelassal/UserAuth/internal/api/delivery"
	"github.com/abdallahelassal/UserAuth/internal/api/middelware"
	"github.com/abdallahelassal/UserAuth/internal/api/route"
	"github.com/gin-gonic/gin"

	// "github.com/abdallahelassal/UserAuth/internal/api/middelware"
	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/repository"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"go.uber.org/zap"

	"gorm.io/gorm"
)


type Container struct{
	UserDelivery 		*delivery.UserDelivary
	RoleDelivery  		*delivery.RoleDelivery
	PermissionDelivary 	*delivery.PermissionDelivery
	PremMiddelware		*middelware.PermissionMiddelWare

	Router 				*gin.Engine
	Handler 			*route.Handler
	Cfg 				bootstrap.Config
	Logger 				*zap.Logger
}

func NewContainer(db *gorm.DB, logger *zap.Logger, cfg bootstrap.Config) *Container {
	r := 				gin.Default()

	//Repository
	userRepo 			:= repository.NewUserRepository(db)
	roleRepo 			:= repository.NewRoleRepository(db)
	permissionRepo 		:= repository.NewPermissionRepository(db)

	//usecase
	userUsecase 		:= usecase.NewUserUseCase(userRepo, roleRepo,permissionRepo, 5*time.Second,cfg.JWTConfig.AccessTokenSecret, time.Duration(cfg.JWTConfig.AccessExpiration)*time.Hour)
	roleUsecase			:= usecase.NewRoleUseCase(*roleRepo, 5 * time.Second)
	permissionUsecase 	:= usecase.NewPermissionUsecase(permissionRepo, 5*time.Second)
	//delivery 
	userDelivery 		:= delivery.NewUserDelivary(userUsecase)
	roleDelivery 		:= delivery.NewRoleDelivery(roleUsecase)
	permissionDelivery 	:= delivery.NewPermissionDelivery(permissionUsecase)
	//middleware
	authMW 				:= middelware.JwtAuthMiddleware(cfg.JWTConfig.AccessTokenSecret)
	premMiddelware 		:= middelware.NewPermissionMiddelware(permissionUsecase,roleUsecase)
	//handler 
	h := route.NewHandler(
		r,
		userDelivery,
		roleDelivery,
		permissionDelivery,
		premMiddelware,
		authMW,
	)
	return &Container{
		UserDelivery: userDelivery,
		RoleDelivery: roleDelivery,
		PermissionDelivary: permissionDelivery,
		PremMiddelware: premMiddelware,
		Router: r,
		Handler: h,
		Cfg: cfg,
		Logger: logger,
	}	
}