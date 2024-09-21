package dns

import (
	"fmt"
	"os"
	"strings"

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

func (d *DNS) ResolveARecord(domain string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, d.config.Servers[0]+":"+d.config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve a record: %w", err)
	}

	var ips []string
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			if strings.HasSuffix(a.A.String(), ".") {
				a.A = a.A[:len(a.A)-1]
			}

			ips = append(ips, a.A.String())
		}
	}

	return ips, nil
}

func (d *DNS) ResolveMXRecord(domain string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, d.config.Servers[0]+":"+d.config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve mx record: %w", err)
	}

	var mxs []string
	for _, ans := range r.Answer {
		if mx, ok := ans.(*dns.MX); ok {
			if strings.TrimSuffix(mx.Mx, ".") != mx.Mx {
				mx.Mx = strings.TrimSuffix(mx.Mx, ".")
			}

			mxs = append(mxs, mx.Mx)
		}
	}

	return mxs, nil
}
