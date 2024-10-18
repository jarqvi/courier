package smtp

import (
	"fmt"
	"net"
	"strings"

	"github.com/jarqvi/courier/internal/dns"
)

func isAllowedDomain(domain string, allowedDomains []string) bool {
	for _, d := range allowedDomains {
		if d == domain {
			return true
		}
	}

	return false
}

func checkSPF(domain string, hostname string) error {
	ips, err := dns.Client.ResolveARecord(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve A record: %w", err)
	}

	txtRecords, err := dns.Client.ResolveTXT(domain)
	if err != nil {
		return fmt.Errorf("failed to resolve SPF record: %w", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			parts := strings.Split(record, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "ip4:") {
					ipRange := strings.TrimPrefix(part, "ip4:")
					_, ipNet, err := net.ParseCIDR(ipRange)
					if err != nil {
						ip := net.ParseIP(ipRange)
						if ip == nil {
							return fmt.Errorf("invalid ip4 format: %s", ipRange)
						}
						ipNet = &net.IPNet{
							IP:   ip,
							Mask: net.CIDRMask(32, 32),
						}
					}

					for _, ip := range ips {
						if ipNet.Contains(net.ParseIP(ip)) {
							return nil
						}
					}
				}

				if strings.HasPrefix(part, "include:") {
					includeDomain := strings.TrimPrefix(part, "include:")
					if err := checkSPF(includeDomain, hostname); err == nil {
						return nil
					}
				}

				if strings.HasPrefix(part, "redirect=") {
					redirectDomain := strings.TrimPrefix(part, "redirect=")
					if err := checkSPF(redirectDomain, hostname); err == nil {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("no valid SPF record found")
}

func checkDKIM(headers []string) error {
	var dkimHeaders []string

	for _, header := range headers {
		if strings.HasPrefix(header, "DKIM-Signature") {
			dkimHeaders = strings.Split(strings.TrimSpace(strings.Split(header, ":")[1]), ";")
		}
	}

	if len(dkimHeaders) == 0 {
		return fmt.Errorf("no DKIM header found")
	}

	var selector string

	for _, header := range dkimHeaders {
		parts := strings.Split(header, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "s=") {
				selector = strings.TrimSpace(strings.Split(part, "=")[1])
				break
			}
		}

		if selector != "" {
			break
		}
	}

	if selector == "" {
		return fmt.Errorf("no selector found")
	}

	var domain string
	
	for _, header := range dkimHeaders {
		parts := strings.Split(header, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "d=") {
				domain = strings.TrimSpace(strings.Split(part, "=")[1])
				break
			}
		}

		if domain != "" {
			break
		}
	}

	if domain == "" {
		return fmt.Errorf("no selector found")
	}

	fqdn := fmt.Sprintf("%s._domainkey.%s", selector, domain)
	txtRecords, err := dns.Client.ResolveTXT(fqdn)
	if err != nil {
		return fmt.Errorf("failed to resolve DKIM record: %w", err)
	}

	fmt.Println(txtRecords)
	return nil
}
