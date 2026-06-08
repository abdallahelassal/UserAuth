package delivery

import (
	"net/http"

	"github.com/abdallahelassal/UserAuth/internal/dtos"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



type RoleDelivery struct{
	RoleUsecase usecase.RoleUseCase
}

func NewRoleDelivery(roleUsecase usecase.RoleUseCase)*RoleDelivery{
	return &RoleDelivery{
		RoleUsecase: roleUsecase,
	}
}

func (r *RoleDelivery) FindAll(g *gin.Context){
	ctx := g.Request.Context()

	
	 roles , err := r.RoleUsecase.FindAll(ctx)
	if err != nil {
		g.JSON(http.StatusBadRequest,gin.H{"error":"roles not found"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"roles":roles,})
}

func (r *RoleDelivery) FindByID(g *gin.Context){
	ctx := g.Request.Context()
	
	paramID := g.Param("id")
	
	roleID , err := uuid.Parse(paramID)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}
	
	role , err := r.RoleUsecase.FindByID(ctx,roleID)
	if err != nil {
		g.JSON(http.StatusBadRequest,gin.H{"error":"role not found"})
		return
	}
	g.JSON(http.StatusOK,gin.H{"role":role})
}

func (r *RoleDelivery) Create(g *gin.Context){
	var req dtos.RoleCreateRequest
	ctx := g.Request.Context()

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest,gin.H{"error":"request not required"})
		return
	}
	
	input := usecase.RoleCreateInput{
		Name: req.Name,
	}

	if err := r.RoleUsecase.Create(ctx,input) ;err != nil {
		g.JSON(http.StatusBadRequest,gin.H{"error":"role not created"})
		return
	}
	g.JSON(http.StatusOK,gin.H{"message":"role created scuccessfuly"})
}

func (r *RoleDelivery) Update(g *gin.Context){
	ctx := g.Request.Context()
	param := g.Param("id")
	roleID , err  := uuid.Parse(param)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error":" request not required"})
		
		return
	}
	id := usecase.RoleUpdateInput{
		ID: roleID,
	}
	if err := r.RoleUsecase.Update(ctx,id) ; err != nil {
			g.JSON(http.StatusInternalServerError,gin.H{"error":"role not update"})
		return
	}

	g.JSON(http.StatusOK,gin.H{"message":"role updated scuccessfuly"})
}

func (r *RoleDelivery) Delete(g *gin.Context){
	ctx := g.Request.Context()
	param := g.Param("id")
	roleID , err  := uuid.Parse(param)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error":"invalid role id"})
		return
	}
	id := usecase.RoleDeleteInput{
		ID: roleID,
	}

	if err := r.RoleUsecase.Delete(ctx,id.ID) ; err != nil {
			g.JSON(http.StatusInternalServerError,gin.H{"error":"role not delete"})
		return
	}

	g.JSON(http.StatusOK,gin.H{"message":"role deleted scuccessfuly"})
}