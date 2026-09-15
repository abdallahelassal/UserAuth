package validator

import (
	"context"
	"regexp"
	"strings"

	"github.com/abdallahelassal/UserAuth/domain"
)


type RegexValidator struct{}

func NewRegexValidator()*RegexValidator{
	return &RegexValidator{}
}
func (r *RegexValidator) Name()string{
	return "regex"
}

func (r *RegexValidator) Check(ctx context.Context, emailStr string)(bool,error){
	var regex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if emailStr == "" {
		return false, domain.ErrInvalidFormat
	}
	parts := strings.Split(emailStr, "@")
	if len(parts) != 2 {
		return  false, domain.ErrInvalidFormat
	}
	domainParts := parts[1]

	if !regex.MatchString(domainParts) {
		return false, domain.ErrInvalidFormat
	}


	return true, nil 
}