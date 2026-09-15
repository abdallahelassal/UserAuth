package validator

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/abdallahelassal/UserAuth/domain"
)

type DNSVlaidator struct{
	resolver *net.Resolver
}

func NewDNSValidator()*DNSVlaidator{
	return &DNSVlaidator{resolver: &net.Resolver{}}
}

func (d *DNSVlaidator) Name()string{
	return "dns_mx"
}

func (d *DNSVlaidator) Check(ctx context.Context, emailSTR string)(bool,error){
	email := strings.Trim(emailSTR, "@")
	if email == ""{
		return false, domain.ErrInvalidFormat
	}

	mxRecords, err := d.resolver.LookupMX(ctx, email)
	if err != nil{
		return false, classifyDNSError(err)
	}
	if len(mxRecords) == 0 {
		return false, domain.ErrDomainNotFound
	}
	return true, nil
}


func classifyDNSError(err error)error{
	var dnsErr net.DNSError
	if errors.As(err, &dnsErr) && dnsErr.IsNotFound{
		return domain.ErrDomainNotFound.With("", err)
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled){
		return domain.ErrTimeout.With("", err)
	}
	return domain.ErrDomainNotFound.With("", err)
}