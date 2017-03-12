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
	config *DNSConfig
}

// NewDNS returns a new DNS instance
func NewDNS(config *DNSConfig) Source {
	dns := new(DNS)

	dns.config = config

	return dns
}

// IsAvailable checks if a domain is available
func (dns *DNS) IsAvailable(domain string) (bool, error) {
	r := dnsr.New(100000)
	_, err := r.ResolveErr(domain, "TXT")

	return err != nil && err == dnsr.NXDOMAIN, nil
}
