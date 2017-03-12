package search

import (
	"github.com/MichaelThessel/gomainr/cache"
	"github.com/MichaelThessel/gomainr/search/source"
)

// Search struct
type Search struct {
	cache  *cache.Cache
	source source.Source
}

var cacheTTL int64 = 86400

// New returns a new Search struct
func New(source source.Source, cache *cache.Cache) *Search {
	s := new(Search)

	s.source = source
	s.cache = cache

	return s
}

// IsAvailable checks the availability of a domain
func (s *Search) IsAvailable(domain string) (bool, error) {
	var available bool

	// Try to load results from cache
	cached, err := s.cache.Get(domain)
	if err == nil && len(cached) > 0 {
		if cached[0] == 't' {
			return true, nil
		}
		return false, nil
	}

	// Fetch from API and save to cache
	available, err = s.source.IsAvailable(domain)
	if err != nil {
		return false, err
	}

	// Cache response
	cached = make([]byte, 1, 1)
	if available {
		cached[0] = 't'
	} else {
		cached[0] = 'f'
	}
	s.cache.Save(domain, cached, cacheTTL)

	return available, nil
}
