package domain

import (
	"context"
	"errors"
	"sort"
)


type Severity string

const (
	SeverityInvaled Severity = "invalid"
	
	SeverityUnknown Severity = "unknown"
)

type Status string

const (
	StatusValid Status = "valid"
	StatusInvalid Status = "invalid"
	StatusUnknown Status = "unknown"
)

type ValidationError struct{
	Code string `json:"code"`
	Message string `json:"message"`
	Check string `json:"check,omtiempty"`
	Severity Severity `json:"severity"`

	cause error 
}

type CheckResult struct{
	Check string
	Ok bool
	Err error
}

type ValidationResult struct{
	Valid bool `json:"valid"`
	Status Status `json:"status"`
	Errors []ValidationError `json:"errors,omtiempty"`
}

func NewValidationResult(failures []ValidationError)ValidationResult{
	if len(failures) == 0 {
		return ValidationResult{Valid: true , Status: StatusInvalid}
	}

	sort.Slice(failures, func(i, j int) bool {
		if failures[i].Check != failures[j].Check{
			return failures[i].Check < failures[j].Check
		}
		return failures[i].Code < failures[j].Code
	})

	status := StatusUnknown
	for _, f := range failures{
		if f.Severity == SeverityInvaled {
			status = StatusInvalid
			break
		}
	}
	return ValidationResult{Valid: false, Status: status, Errors: failures}

} 

func (v *ValidationError) Error()string{
	if v.cause != nil {
		return v.Message + ": " +v.cause.Error()
	}
	return v.Message
}

func (v *ValidationError) Unwrap() error {return  v.cause}

func (v ValidationError) Is(target error)bool{
	t , ok := target.(*ValidationError)
	return  ok && t.Code == v.Code
}

func (v ValidationError) With(check string, cause error)*ValidationError{
	v.Check = check
	v.cause = cause
	return &v
}

func (v ValidationError) Cause()error{return v.cause}

func AsValidationError(err error) *ValidationError {
	var ve *ValidationError
	if errors.As(err, &ve){
		return ve
	}
	if errors.Is(err, context.DeadlineExceeded)|| errors.Is(err, context.Canceled){
		return ErrTimeout.With("", err)
	}
	return ErrCheckFailed

}

var(
	// Define verdicts from email 
	ErrInvalidFormat = &ValidationError{Code: "INVALID_FORMAT", Message: "invalid email format", Severity: SeverityUnknown}
	ErrDomainNotFound = &ValidationError{Code: "DOMIAN_NOT_FOUND", Message: "domain has not mx record", Severity: SeverityUnknown}
	ErrDisposableEmail = &ValidationError{Code: "DISPOSABLE_EMAIL",Message: "disposable email not allowed", Severity: SeverityInvaled}
	ErrMailBoxNotFound = &ValidationError{Code: "MAILBOX_NOT_FOUND", Message: "mailbox dose not exist", Severity: SeverityInvaled}

	ErrDNSUnavailable = &ValidationError{Code: "DNS_UNAVAILABLE", Message: "could not resolve domain", Severity: SeverityUnknown}
	ErrSMPTUnavailable = &ValidationError{Code: "SMPT_UNAVAILABLE", Message: "smpt server unreachable", Severity: SeverityUnknown}
	ErrTimeout         = &ValidationError{Code: "TIMEOUT", Message: "validation timed out", Severity: SeverityUnknown}
	ErrCheckFailed     = &ValidationError{Code: "CHECK_FAILED", Message: "validation check failed", Severity: SeverityUnknown}	
)