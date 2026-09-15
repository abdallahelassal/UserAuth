package validator

import (
	"context"
	"strings"
	"sync"

	"github.com/abdallahelassal/UserAuth/domain"
)


type DisposableValidator struct{
	blocked map[string]struct{}
	mu sync.RWMutex
}

var defaultBlockedDomains = []string{
	"mailinator.com", "tempmail.com", "10minutemail.com",
	"guerrillamail.com", "throwawaymail.com", "yopmail.com",
	"fakeinbox.com", "sharklasers.com", "getairmail.com",
	"burnermail.io", "temp-mail.org", "mailnesia.com",
}

func NewDisposableValidator()*DisposableValidator{
	v := &DisposableValidator{blocked: make(map[string]struct{})}
	for _, b := range defaultBlockedDomains {
		v.blocked[strings.ToLower(b)] = struct{}{}
	}
	return v
}

func (d *DisposableValidator) AddDomain(domainName string){
	d.mu.Lock()
	defer d.mu.Unlock()
	d.blocked[strings.ToLower(domainName)] = struct{}{}
}

func (d *DisposableValidator) Check(ctx context.Context, emailStr string)(bool,error){
	if emailStr == ""{
		return false, domain.ErrInvalidFormat
	}
	parts := strings.Split(emailStr, "@")
	if len(parts) != 2{
		return false, domain.ErrInvalidFormat
	}
	domainPart := parts[1]
	emailDomain := strings.ToLower(domainPart)

	d.mu.RLock()
	_, blocked := d.blocked[emailDomain]
	d.mu.RUnlock()
	if blocked{
		return false, domain.ErrDisposableEmail
	}

	return true, nil
}

