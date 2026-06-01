package delivery

import (
	"net/http"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/internal/bootstrap"
	"github.com/abdallahelassal/UserAuth/internal/dtos"
	"github.com/abdallahelassal/UserAuth/internal/usecase"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	encryptedPassword , err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	req.Password = string(encryptedPassword)

	user := &domain.User{
		Username: req.UserName,
		Email: req.Email,
		Password: string(encryptedPassword),
		IsActive: true,
	}
	err = d.UserUseCase.Signup(g, user)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}