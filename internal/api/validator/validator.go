package validator

import (
	"context"


)

type EmailValidator interface{
	Check(ctx context.Context, email string) (bool , error)
	Name() string
}