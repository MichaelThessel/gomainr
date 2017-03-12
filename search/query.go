package search

import "strings"

// BuildQuery builds domain names from given parts
func (s *Search) BuildQuery(first []string, second []string, tlds []string, tldSubstitutions bool) []string {
	var baseDomains []string

	// Build all possible combinations of first and second level
	for i := 0; i < len(first); i++ {
		if len(second) == 0 {
			baseDomains = append(baseDomains, first[i])
		} else {
			for j := 0; j < len(second); j++ {
				baseDomains = append(baseDomains, first[i]+second[j])
			}
		}
	}

	// Append TLDs
	var domains []string
	for _, domain := range baseDomains {
		for i := 0; i < len(tlds); i++ {
			domains = append(domains, domain+"."+tlds[i])
		}
	}

	// Add TLD substiturion
	if tldSubstitutions {
		tldSub(&baseDomains, &domains)
	}

	return domains
}

// tldSub substites the end of a base domain with a TLD if the end matches a
// existing TLD.
// i.e: superyachts super.yachts
func tldSub(baseDomains *[]string, domains *[]string) {
	// Add domains with TLD substitutions
	for _, domain := range *baseDomains {
		for _, tld := range validTlds {
			if strings.HasSuffix(strings.ToLower(domain), tld) {
				*domains = append(*domains, domain[:len(domain)-len(tld)]+"."+tld)
			}
		}
	}
}
