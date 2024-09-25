package smtp

import (
	"fmt"
	"net"
	"strings"

	"github.com/jarqvi/courier/internal/dns"
	"github.com/jarqvi/courier/internal/log"
)

func isAllowedDomain(domain string, allowedDomains []string) bool {
	for _, d := range allowedDomains {
		if d == domain {
			return true
		}
	}

	return false
}

func sanitizeEmailAddress(email string) error {
	if strings.Contains(email, "<") || strings.Contains(email, ">") {
		return fmt.Errorf("invalid characters in email address")
	}

	return nil
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

	log.Logger.Debug("domain: ", domain)
	log.Logger.Debug("hostname: ", hostname)
	log.Logger.Debug("ips: ", ips)
	log.Logger.Debug("txt records: ", txtRecords)

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
