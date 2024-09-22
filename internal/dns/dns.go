package dns

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

type DNS struct {
	config *dns.ClientConfig
}

func Init() (*DNS, error) {
	DNS_CONFIG_PATH := os.Getenv("DNS_CONFIG_PATH")
	if DNS_CONFIG_PATH == "" {
		DNS_CONFIG_PATH = "/etc/resolv.conf"
	}

	config, err := dns.ClientConfigFromFile(DNS_CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to get dns config: %w", err)
	}

	return &DNS{
		config: config,
	}, nil
}
