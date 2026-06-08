package route

import (
	"github.com/abdallahelassal/UserAuth/internal/api/delivery"
	"github.com/abdallahelassal/UserAuth/internal/api/middelware"
	

	"github.com/gin-gonic/gin"
)


type Handler struct{
	r 			*gin.Engine
	userHandler *delivery.UserDelivary
	roleHandler *delivery.RoleDelivery
	permissionHandler *delivery.PermissionDelivery
	permMiddleware 	*middelware.PermissionMiddelWare
	authMiddleware gin.HandlerFunc
}
func NewHandler(
	r *gin.Engine,
	user *delivery.UserDelivary,
	role *delivery.RoleDelivery,
	perm *delivery.PermissionDelivery,
	permMiddleware *middelware.PermissionMiddelWare,
	auth gin.HandlerFunc,
) *Handler {
	return &Handler{
		r:              r,
		userHandler:    user,
		roleHandler:    role,
		permissionHandler: perm,
		permMiddleware: permMiddleware,
		authMiddleware: auth,
	}
}

func (h *Handler) SetupRoutes(){

	
	api := h.r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/signup", h.userHandler.Signup)
		auth.POST("/login", h.userHandler.Login)
	}
	
	user := api.Group("/user")
	user.Use(h.authMiddleware)
	{
		user.GET("/:id",h.userHandler.Profile)
		user.PUT("/:id/roles",h.userHandler.AssignRoles)
	}
	roles := api.Group("/roles")
	roles.Use(h.authMiddleware)
	{
		roles.GET("/",h.permMiddleware.Required("roles:view"),h.roleHandler.FindAll)
		roles.GET("/:id", h.permMiddleware.Required("roles:view"),h.roleHandler.FindByID)
		roles.POST("/create", h.permMiddleware.Required("roles:manage"), h.roleHandler.Create)
		roles.PUT("/:id",h.permMiddleware.Required("roles:manage"),h.roleHandler.Update)
		roles.DELETE("/:id",h.permMiddleware.Required("roles:manage"),h.roleHandler.Delete)
	}
	permissions := api.Group("/permissions")
	permissions.Use(h.authMiddleware)
	{
		permissions.GET("/", h.permMiddleware.Required("permissions:view"),h.permissionHandler.FindAll)
		permissions.GET("/user/:id",h.permMiddleware.Required("permissions:view"), h.permissionHandler.FindPermissionByUserID)
	}



} 
