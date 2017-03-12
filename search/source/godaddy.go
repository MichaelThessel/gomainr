package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GoDaddyConfig holds the configuration for the namecheap.com source
type GoDaddyConfig struct {
	Key     string
	Secret  string
	Enabled bool
}

type goDaddyResponse struct {
	Available  bool
	Domain     string
	Definitive bool
	Price      int
	Currency   string
	Period     int
	Message    string
	Code       string
	Name       string
}

// GoDaddy handles godaddy.com API requests
type GoDaddy struct {
	Source
	config *GoDaddyConfig
}

// NewGoDaddy returns a new GoDaddy instance
func NewGoDaddy(config *GoDaddyConfig) Source {
	gd := new(GoDaddy)

	gd.config = config

	return gd
}

// IsAvailable checks if a domain is available
func (gd *GoDaddy) IsAvailable(domain string) (bool, error) {
	client := &http.Client{}

	v := url.Values{}
	v.Set("domain", domain)
	req, err := http.NewRequest("GET", "https://api.godaddy.com/v1/domains/available?"+v.Encode(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", gd.config.Key, gd.config.Secret))
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.New("Couldn't connect to API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.New("Couldn't read API response")
	}

	var gdResponse goDaddyResponse
	if err := json.Unmarshal(body, &gdResponse); err != nil {
		return false, errors.New("Couldn't parse API response")
	}

	if gdResponse.Message != "" {
		return false, errors.New(gdResponse.Message)
	}

	return gdResponse.Available, nil
}
