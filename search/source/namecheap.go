package source

import (
	gonc "github.com/billputer/go-namecheap"
)

// NameCheapConfig holds the configuration for the namecheap.com source
type NameCheapConfig struct {
	APIUser  string
	APIToken string
	UserName string
	Enabled  bool
}

// NameCheap handles namecheap.com API requests
type NameCheap struct {
	config *NameCheapConfig
}

// NewNameCheap returns a new NameCheap instance
func NewNameCheap(config *NameCheapConfig) Source {
	nc := new(NameCheap)

	nc.config = config

	return nc
}

// IsAvailable checks if a domain is available
func (nc *NameCheap) IsAvailable(domain string) (bool, error) {
	client := gonc.NewClient(nc.config.APIUser, nc.config.APIToken, nc.config.UserName)

	result, err := client.DomainsCheck(domain)

	if err == nil {
		return result[0].Available, err
	}

	return false, err
}
