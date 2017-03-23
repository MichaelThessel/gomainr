package source

import (
	"github.com/domainr/dnsr"
)

// DNSConfig holds the configuration for the DNS source
type DNSConfig struct {
	Enabled bool
}

// DNS handles checking availablility of domain names via DNS
type DNS struct {
	config   *DNSConfig
	resolver *dnsr.Resolver
}

// NewDNS returns a new DNS instance
func NewDNS(config *DNSConfig) Source {
	return &DNS{
		config:   config,
		resolver: dnsr.New(0),
	}
}

// IsAvailable checks if a domain is available
func (dns *DNS) IsAvailable(domain string) (_ bool, err error) {
	_, err = dns.resolver.ResolveErr(domain, "TXT")
	if err == dnsr.NXDOMAIN {
		return true, nil
	}
	return
}
