package validator

import (
	"context"
	"errors"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	
)


type SMTPValidator struct{
	fromAdder 	string 
	timeout 	time.Duration
	resolver 	*net.Resolver
}

func NewSMTPValidator(fromAdder string) *SMTPValidator{
	return &SMTPValidator{
		fromAdder: fromAdder,
		timeout: 5* time.Second,
		resolver: &net.Resolver{},
	}
}

func (s *SMTPValidator) Name()string{
	return  "smtp"
}

func (s *SMTPValidator) Check(ctx context.Context, emailStr string) (bool,error){
	
	if emailStr == ""{
		return false, domain.ErrInvalidFormat
	}
	parts := strings.Split(emailStr, "@")
	if len(parts) != 2 {
		return false, domain.ErrInvalidFormat
	}

	domainPart := parts[1]

	mxRecords , err := s.resolver.LookupMX(ctx, domainPart)
	if err != nil {
		return false, classifyDNSError(err)
	}
	if len(mxRecords) == 0 {
		return false , domain.ErrDomainNotFound
	}

	
	dialer := &net.Dialer{Timeout: s.timeout}
	for _ , mx := range mxRecords {
		host := strings.TrimPrefix(mx.Host, ".")

		conn , err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, "25"))
		if err != nil {
			continue
		}

		client , err := smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			continue
		}
		defer client.Quit()
		if err := client.Hello("validate.localhost"); err != nil {
			return false, domain.ErrSMPTUnavailable.With("", err)
		}
		if err := client.Mail(s.fromAdder); err != nil {
			return false, domain.ErrSMPTUnavailable.With("", err)
		}
		if err := client.Rcpt(emailStr); err != nil {
			return false, classifyRcptError(err)
		}
		return true, nil
	}
	return true , nil

} 

func classifyRcptError(err error)error{
	var protoerr textproto.Error
	if errors.As(err, &protoerr) && protoerr.Code >= 500 && protoerr.Code <= 600 {
		return domain.ErrDomainNotFound
	} 
	return domain.ErrSMPTUnavailable.With("", err)
}






