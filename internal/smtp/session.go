package smtp

import (
	"io"
	"strings"

	"github.com/emersion/go-smtp"
	"github.com/jarqvi/courier/internal/log"
)

type Session struct {
	AllowedDomains []string
	From           string
	To             string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.From = from

	log.Logger.Info("mail from: ", from)

	return nil
}

func (s *Session) Rcpt(to string, opts  *smtp.RcptOptions) error {
	domain := strings.Split(to, "@")[1]

	for _, allowedDomain := range s.AllowedDomains {
		if domain == allowedDomain {
			s.To = to
			log.Logger.Info("rcpt to:", to)
			return nil
		}
	}

	return smtp.ErrServerClosed
}

func (s *Session) Data(r io.Reader) error {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Logger.Info("data: %s", string(buf[:n]))
	}

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
