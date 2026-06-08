package middelware

import (
	"net/http"

	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionMiddelWare struct {
	PermissionUsecase usecase.PermissionUsecase
	RoleUsecase       usecase.RoleUseCase
}

func NewPermissionMiddelware(pu usecase.PermissionUsecase, ru usecase.RoleUseCase) *PermissionMiddelWare {
	return &PermissionMiddelWare{
		PermissionUsecase: pu,
		RoleUsecase:       ru,
	}
}
func (p *PermissionMiddelWare) Required(requireParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not required"})
			c.Abort()
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
			c.Abort()
			return
		}

		paramID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			c.Abort()
			return
		}
		directPermission, err := p.PermissionUsecase.GetPermissionsByUserID(ctx, paramID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user"})
			c.Abort()
			return
		}
		for _, v := range directPermission {
			if v.Name == requireParam {
				c.Next()
				return
			}
		}

		roles , err := p.RoleUsecase.GetRolesByUserID(ctx,paramID)
		if err != nil || len(roles) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"error": "user not authenticated"})
			c.Abort()
			return 
		}
		var roleIDs []uuid.UUID
		
			for _, r := range roles {
				
					roleIDs = append(roleIDs, r.ID)
				
			}

		permissionFromRole, err := p.PermissionUsecase.GetPermissionByRoleIDs(ctx,roleIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get permissions from roles"})
			c.Abort()
			return
		}

		for _, v := range permissionFromRole {
			if v.Name == requireParam {
				c.Next()
				return
			}
		}
		

		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		c.Abort()
	}
}
