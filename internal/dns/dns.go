package dns

import (
	"fmt"
	"os"

	"github.com/jarqvi/courier/internal/log"
	"github.com/miekg/dns"
)

type DNS struct {
	config *dns.ClientConfig
}

var Client *DNS

func Init() error {
	DNS_CONFIG_PATH := os.Getenv("DNS_CONFIG_PATH")
	if DNS_CONFIG_PATH == "" {
		DNS_CONFIG_PATH = "/etc/resolv.conf"
	}

	config, err := dns.ClientConfigFromFile(DNS_CONFIG_PATH)
	if err != nil {
		return fmt.Errorf("failed to get dns config: %w", err)
	}

	Client = &DNS{
		config: config,
	}

	log.Logger.Info("dns client initialized")

	return nil
}
