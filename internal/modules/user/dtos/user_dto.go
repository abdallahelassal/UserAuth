package dtos

type CreateUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`	
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

type FetchUserResponse struct{
	Users		[]UserResponse		`json:"users"`
	NextCursor	string				`json:"next_cursor,omitempty"`
}