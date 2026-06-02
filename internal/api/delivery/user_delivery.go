package delivery

import (
	"net/http"

	
	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/dtos"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/gin-gonic/gin"
	
)


type UserDelivary struct {
	UserUseCase *usecase.UserUseCase
	Cfg bootstrap.Config
}

func NewUserDelivary(userUseCase *usecase.UserUseCase, cfg bootstrap.Config) *UserDelivary {
	return &UserDelivary{
		UserUseCase: userUseCase,
		Cfg: cfg,
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
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"token": token})
} 
