package smtp

import (
	"bufio"
	"io"
	"net/mail"
	"strings"

	"github.com/emersion/go-smtp"
	"github.com/jarqvi/courier/internal/log"
)

type Session struct {
	AllowedDomains []string
	Hostname       string
	From           string
	To             string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.From = from

	addr, err := mail.ParseAddress(from)
	if err != nil {
		log.Logger.Error("invalid sender address: ", from, err)

		return &smtp.SMTPError{
			Code:    550,
			Message: "Invalid sender",
			EnhancedCode: smtp.EnhancedCode{
				5, 1, 1,
			},
		}
	}

	if err := sanitizeEmailAddress(addr.Address); err != nil {
		log.Logger.Error("invalid sender address: ", from, err)

		return &smtp.SMTPError{
			Code:    550,
			Message: "Invalid sender",
			EnhancedCode: smtp.EnhancedCode{
				5, 1, 1,
			},
		}
	}

	domain := strings.Split(addr.Address, "@")[1]

	if err := checkSPF(domain, s.Hostname); err != nil {
		log.Logger.Error("spf check failed: ", err)

		return &smtp.SMTPError{
			Code:    550,
			Message: "SPF check failed",
			EnhancedCode: smtp.EnhancedCode{
				5, 7, 1,
			},
		}
	}

	log.Logger.Info("mail from: ", from)

	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	addr, err := mail.ParseAddress(to)
	if err != nil {
		log.Logger.Error("invalid recipient: ", to, err)

		return &smtp.SMTPError{
			Code:    550,
			Message: "Invalid recipient",
			EnhancedCode: smtp.EnhancedCode{
				5, 1, 1,
			},
		}
	}

	domain := strings.Split(addr.Address, "@")[1]

	if !isAllowedDomain(domain, s.AllowedDomains) {
		return &smtp.SMTPError{
			Code:    550,
			Message: "Relay access denied",
			EnhancedCode: smtp.EnhancedCode{
				5, 7, 1,
			},
		}
	}

	s.To = to
	log.Logger.Info("rcpt to: ", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	reader := bufio.NewReader(r)

	var headers []string
	var body string
	isHeaders := true

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if isHeaders && strings.TrimSpace(line) == "" {
			isHeaders = false
			continue
		}

		if isHeaders {
			headers = append(headers, line)
		} else {
			body += line
		}
	}

	log.Logger.Info("Headers:")
	for _, header := range headers {
		log.Logger.Info("%s", strings.TrimSpace(header))
	}

	log.Logger.Info("Body:")
	log.Logger.Info("%s", body)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
