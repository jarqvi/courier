package smtp

import (
	"fmt"
	"os"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/jarqvi/courier/internal/db"
	"github.com/jarqvi/courier/internal/log"
)

var Server *smtp.Server
var ServerError error


func Init() {
	allowedDomains, err := db.Instance.GetAllDomains()
	if err != nil {
		ServerError = fmt.Errorf("failed to get allowed domains: %w", err)
		return
	}

	var allowedDomainNames []string
	for _, domain := range allowedDomains {
		allowedDomainNames = append(allowedDomainNames, domain.Name)
	}

	s := smtp.NewServer(&Backend{AllowedDomains: allowedDomainNames})

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost"
	}

	s.Addr = ":25"
	s.Domain = domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	go func() {
		if err := s.ListenAndServe(); err != nil {
			ServerError = fmt.Errorf("failed to start smtp server: %w", err)
		}
	}()

	log.Logger.Info("smtp server started on ", s.Addr)

	Server = s
}

func Shutdown() {
	log.Logger.Info("shutting down smtp server")
	Server.Close()
}
