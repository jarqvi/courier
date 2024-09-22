package dns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

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
