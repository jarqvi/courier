package dns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func (d *DNS) ResolveTXT(domain string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, d.config.Servers[0]+":"+d.config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve txt record: %w", err)
	}

	var txts []string
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.TXT); ok {
			txts = append(txts, strings.Join(a.Txt, ""))
		}
	}

	return txts, nil
}
