package delivery

import (
	"net/http"

	
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type PermissionDelivery struct {
	PermissionUsecase usecase.PermissionUsecase
}

func NewPermissionDelivery(permissionUsecase usecase.PermissionUsecase)*PermissionDelivery{
	return &PermissionDelivery{
		PermissionUsecase: permissionUsecase,
	}
}

func (p *PermissionDelivery) FindAll(g *gin.Context){
	ctx := g.Request.Context()

	permissions , err := p.PermissionUsecase.FindAllPermissions(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError,gin.H{"error":"permissions not found"})
		return  
	}
	g.JSON(http.StatusOK,gin.H{"permissions":permissions})

}

func (p *PermissionDelivery) FindPermissionByUserID(g *gin.Context){
	ctx := g.Request.Context()
	

	param := g.Param("id")
	paramsID, err := uuid.Parse(param)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error":"invalid user id"})
		return
	}

	permission , err := p.PermissionUsecase.GetPermissionsByUserID(ctx,paramsID)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error":"permissions not found"})
		return 		
	}
	g.JSON(http.StatusOK, gin.H{"permissions":permission})
}
func (p *PermissionDelivery) Create(g *gin.Context){
	ctx := g.Request.Context()

	var req struct{
		Name string `json:"name" binding:"required"`
	}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest,gin.H{"error":"invalid request"})
		return 
	}

	perm := &usecase.PermissionInput{
		Name: req.Name,
	}

	if err := p.PermissionUsecase.Create(ctx,perm); err != nil {
		g.JSON(http.StatusInternalServerError,gin.H{"error":"failed to create permission"})
		return 
	}
	g.JSON(http.StatusOK,gin.H{"message":"permission created successfully"})
}