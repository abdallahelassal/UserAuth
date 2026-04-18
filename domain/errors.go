package domain

import "errors"




var (
	
	ErrIternalServer = errors.New("internal server Error")

	ErrNotFound = errors.New("your request items not found")

	ErrConfilct = errors.New("your item is already exit")

	ErrBadParamInput = errors.New("givin param is not vaild")
)