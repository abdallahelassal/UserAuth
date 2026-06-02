package domain

import "errors"




var (
	
	ErrInternalServer = errors.New("internal server Error")

	ErrNotFound = errors.New("your request items not found")

	ErrConfilct = errors.New("your item is already exit")

	ErrBadParamInput = errors.New("givin param is not vaild")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidToken = errors.New("invalid token")

	ErrTokenExpired = errors.New("token has expired")
	
	ErrRoleNameRequired = errors.New("role name is required")
)