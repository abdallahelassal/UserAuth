package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/abdallahelassal/UserAuth/domain/ports"
)

type EmailUsecase struct{
	Validators []ports.EmailValidator
	Timeout time.Duration
}

func NewEmailUsecase(validators []ports.EmailValidator, timeout time.Duration)*EmailUsecase{
	return &EmailUsecase{
		Validators: validators,
		Timeout: timeout,
	}
}

func (e *EmailUsecase) Validate(ctx context.Context, emailStr string) domain.ValidationResult {
	email := domain.User{Email: emailStr}

	ctx , cancel:= context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	results := make(chan domain.CheckResult, len(e.Validators) )
	var wg sync.WaitGroup

	for _, v :=range e.Validators {
		wg.Add(1)
		go func(validator ports.EmailValidator) {
			defer wg.Done()

			ok , err := validator.Check(ctx,email)
			results <- domain.CheckResult{Check: validator.Name(), Ok: ok, Err: err}
		}(v)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return  e.aggregate(ctx , results, len(e.Validators))
}

func (e *EmailUsecase) aggregate(ctx context.Context,
	result <-chan domain.CheckResult,
	total int)domain.ValidationResult{
		var failures []domain.ValidationError
		collected := 0

		for collected < total{
			select {
			case r , ok := <- result:
				if !ok {
					return domain.NewValidationResult(failures)
				}
				collected++ 
				if !r.Ok{
					failures = append(failures, *failureOf(r))
				}
			case <- ctx.Done():
				failures = append(failures, *domain.ErrTimeout.With("pipline", ctx.Err()))
				return domain.NewValidationResult(failures)	
			}
		}
		return domain.NewValidationResult(failures)
	}

	func failureOf(r domain.CheckResult) *domain.ValidationError{
		ve := domain.AsValidationError(r.Err)
		if ve.Check == ""{
			ve = ve.With(r.Check, ve.Cause())
		}
		return  ve
	}

	func (e *EmailUsecase) ValidateSequntial(ctx context.Context, emailStr string)domain.ValidationResult{
		email := domain.User{Email: emailStr}

		var failures []domain.ValidationError
		for _, v := range e.Validators {
			ok , err := v.Check(ctx, email)
			if !ok {
				failures = append(failures, *failureOf(domain.CheckResult{Check: v.Name(),Ok: ok,Err: err}))
			}
		}
		return domain.NewValidationResult(failures)
	}
