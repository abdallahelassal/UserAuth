package validator

import (
	"context"

	"github.com/abdallahelassal/UserAuth/domain"
)

type EmailValidator interface{
	Check(ctx context.Context, email domain.User) (bool , error)
	Name() string
}