package dns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

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
