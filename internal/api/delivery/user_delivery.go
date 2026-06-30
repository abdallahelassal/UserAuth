package delivery

import (
	"log"
	"net/http"

	"github.com/abdallahelassal/UserAuth/internal/dtos"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type UserDelivary struct {
	UserUseCase *usecase.UserUseCase
	
}

func NewUserDelivary(userUseCase *usecase.UserUseCase) *UserDelivary {
	return &UserDelivary{
		UserUseCase: userUseCase,
		
	}
}

func (d *UserDelivary) Signup(g *gin.Context){
	var req dtos.CreateUserRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	_, err := d.UserUseCase.GetByEmail(g, req.Email);
	if  err == nil {
		g.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}


	ctx := g.Request.Context()
	input := usecase.CreateUserInput{
		UserName: req.UserName,
		Email: req.Email,
		Password: req.Password,
	}

	if err := d.UserUseCase.Signup(ctx, input); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (d *UserDelivary) Login(g *gin.Context){
	var req dtos.LoginUserRequest
	if err := g.ShouldBindJSON(&req);err != nil {
		 Json(g,Error(http.StatusBadGateway,"invalid request body"))
		 return
	}

	ctx := g.Request.Context()
	input := usecase.LoginUserInput{
		Email: req.Email,
		Password: req.Password,
	}
	token , err :=  d.UserUseCase.Login(ctx,input)
	if err != nil {
		log.Printf("login err %+v \n", err)
		
		g.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"token": token})
} 

func (d *UserDelivary) Profile(g *gin.Context){
	ctx  := g.Request.Context()
	idParam := g.Param("id")
	
	userID , err := uuid.Parse(idParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	


	user , err := d.UserUseCase.FindByID(ctx,userID)
		if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"user": user})

}

func (d *UserDelivary) AssignRoles(g *gin.Context){
	ctx := g.Request.Context()
	idParam := g.Param("id")
	userID , err := uuid.Parse(idParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id"`
	}
	if err := g.ShouldBindJSON(&req);err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error":"invailed role_id"})
		return
	}
	roleID := req.RoleID

	
	if err := d.UserUseCase.AssignRole(ctx,userID,roleID); err!= nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error":"not found roles "})
		return
	}
	g.JSON(http.StatusOK, gin.H{"userRoles": "assigned successfully"})

}

func (u *UserDelivary) Me(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	result, err := u.UserUseCase.GetFullProfile(ctx, userID)
	if err != nil {
		
		c.JSON(500, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(200, result)
}