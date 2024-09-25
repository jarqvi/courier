package smtp

import (
	"github.com/emersion/go-smtp"
)

type Backend struct {
	AllowedDomains []string
}

func (bkd *Backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	return &Session{
		AllowedDomains: bkd.AllowedDomains,
		Hostname: state.Hostname(),
	}, nil
}
