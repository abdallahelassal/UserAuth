package container

import (

	"time"


	"github.com/abdallahelassal/UserAuth/internal/api/delivery"
	"github.com/abdallahelassal/UserAuth/internal/api/middelware"
	"github.com/abdallahelassal/UserAuth/internal/api/route"
	"github.com/abdallahelassal/UserAuth/internal/api/validator"
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

	//validators
	disposale := validator.NewDisposableValidator()
	dns := validator.NewDNSValidator()
	regex := validator.NewRegexValidator()
	smtp := validator.NewSMTPValidator("no-reply@yourdomain.com")
	
	//Repository
	userRepo 			:= repository.NewUserRepository(db)
	roleRepo 			:= repository.NewRoleRepository(db)
	permissionRepo 		:= repository.NewPermissionRepository(db)

	//usecase
	emailUsecase 		:= usecase.NewEmailUsecase(
		[]validator.EmailValidator{disposale,dns,regex,smtp},
		5*time.Second,
	)
	userUsecase 		:= usecase.NewUserUseCase(userRepo, roleRepo, permissionRepo, emailUsecase, 5*time.Second,cfg.JWTConfig.AccessTokenSecret, time.Duration(cfg.JWTConfig.AccessExpiration))
	roleUsecase			:= usecase.NewRoleUseCase(roleRepo, 5 * time.Second)
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